// Package snowflake provides a very simple Twitter Snowflake generator and parser.
package snowflake

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/docodex/gopkg/internal"
)

const (
	// epochTimestamp is the time since which the Snowflake time is defined as the timestamp.
	// The default epoch timestamp is set to 2026-01-01 00:00:00 +0000 UTC in milliseconds.
	// You may customize this to set a different epoch timestamp for your application.
	// The epoch timestamp should be before the current time.
	epochTimestamp = 1767225600000 // time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	// timeUnit is the internal time unit of the Snowflake algorithm, measured in milliseconds.
	// The default time unit is set to 1 millisecond.
	// You may customize this to set a different time unit for your application.
	timeUnit = 1 // 1ms

	// timeBits holds the number of bits to use for timestamp.
	// The timeBits is calculated by 63 - nodeBits - sequenceBits.
	// The default timeBits is set to 42: 63-10-11.
	// The timeBits should be between 36 (inclusive) and 52 (inclusive).
	timeBits = 42

	// nodeBits holds the number of bits to use for node id.
	// The default nodeBits is set to 10.
	// You may customize this to set a different length for your application.
	// The nodeBits should be between 1 (inclusive) and 26 (inclusive).
	// Remember, you have a total (63 - timeBits) bits to share between Node/Sequence.
	nodeBits = 10

	// sequenceBits holds the number of bits to use for sequence number.
	// The default sequenceBits is set to 11.
	// You may customize this to set a different length for your application.
	// The sequenceBits should be between 1 (inclusive) and 26 (inclusive).
	// Remember, you have a total (63 - timeBits) bits to share between Node/Sequence.
	sequenceBits = 11

	// maxTolerableRollback is the maximum tolerable clock rollback duration.
	// If the clock is rolled back more than this, Generate will return an error
	// instead of busy-waiting.
	maxTolerableRollback = 2000 // 2 seconds in milliseconds
)

const (
	timeShift    = nodeBits + sequenceBits
	maxTimestamp = -1 ^ (-1 << timeBits) // (1 << timeBits) - 1

	nodeShift = sequenceBits
	maxNode   = -1 ^ (-1 << nodeBits) // (1 << nodeBits) - 1
	nodeMask  = maxNode << nodeShift

	maxSequence  = -1 ^ (-1 << sequenceBits) // (1 << sequenceBits) - 1
	sequenceMask = maxSequence
)

var (
	ErrOverTimeLimit     = errors.New("over the timestamp limit")
	ErrOverNodeLimit     = errors.New("over the node id limit")
	ErrOverSequenceLimit = errors.New("over the sequence number limit")
	ErrCheckNodeFailed   = errors.New("check node id failed")
	ErrClockRollback     = errors.New("clock rolled back too far, exceeds maximum tolerable rollback")
)

// Snowflake is a distributed unique ID generator inspired by twitter snowflake.
// By default, a Snowflake ID is composed of
// - 42 bits for time in units of 1 msec
// - 10 bits for a node id
// - 11 bits for a sequence number
//
// Node represents the unique ID of a snowflake instance.
// The default node is set to the lower 8 bits of the private IP address.
// You may customize this to set a different value for your application.
//
// CheckNode validates the uniqueness of a node id.
// If checkNode returns false, the instance will not be created.
// If checkNode is nil, no validation is done.
type Snowflake struct {
	mu sync.Mutex

	epoch    time.Time
	elapsed  int64
	node     int64
	sequence int64

	checkNode func(node int64) bool
}

func New(opts ...Option) (*Snowflake, error) {
	s := &Snowflake{node: -1}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	if s.node == -1 {
		// The default node is set to the lower 8 bits of the private IP address.
		node, err := internal.Lower8BitPrivateIPv4()
		if err != nil {
			return nil, fmt.Errorf("failed to set node id: %w", err)
		}
		s.node = int64(node)
	}

	if s.node > maxNode {
		return nil, ErrOverNodeLimit
	}
	if s.checkNode != nil && !s.checkNode(s.node) {
		return nil, ErrCheckNodeFailed
	}

	// Construct s.epoch so that it carries a monotonic clock reading.
	// time.Unix(...) returns a Time with only all-clock data (no monotonic reading).
	// If we stored that directly in s.epoch, every subsequent time.Since(s.epoch) call
	// would fall back to wall-clock subtraction, making the ID generator vulnerable to
	// system clock adjustments (e.g. NTP corrections) which could produce duplicate or
	// out-of-order timestamps.
	//
	// time.Now() returns a Time that includes both wall-clock and monotonic readings.
	// By computing the Duration between now and the epoch wall-clock, then adding that
	// Duration back to now via now.Add(...), the result preserves the monotonic reading.
	// This ensures time.Since(s.epoch) always uses the monotonic clock, making elapsed
	// time measurements immune to wall-clock jumps.
	now := time.Now()
	s.epoch = now.Add(time.Unix(epochTimestamp/1_000, (epochTimestamp%1_000)*1_000_000).Sub(now))

	// NOTE: setting s.sequence = maxSequence here would cause the first Generate() call
	// to wrap the sequence to 0 and wait for the next time tick when now == s.elapsed (== 0),
	// avoiding generating an ID with all-zero timestamp and sequence bits (which could be
	// confused with an uninitialized value). In practice this never triggers because the
	// epoch (2026-01-01) is in the past, so now >> 0 on the first call, and the else branch
	// resets sequence to 0 directly. Left commented out as unnecessary.
	//
	// s.sequence = maxSequence

	return s, nil
}

// Generate creates and returns a unique snowflake ID.
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same node id
func (s *Snowflake) Generate() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Since(s.epoch).Milliseconds() / timeUnit
	if now == s.elapsed {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// wait to next time unit
			for now <= s.elapsed {
				now = time.Since(s.epoch).Milliseconds() / timeUnit
			}
		}
	} else {
		// Clock rollback detection: when now < s.elapsed, the clock has gone backward.
		// With the monotonic clock trick in New(), this should never happen under normal
		// conditions. It is kept as defense-in-depth for exceptional scenarios such as
		// VM live migration or OS-level monotonic clock anomalies.
		//
		// Since both are signed int64, if now > s.elapsed (normal forward progression),
		// s.elapsed-now is negative which is always <= maxTolerableRollback.
		// Only when now < s.elapsed does s.elapsed-now become a positive value,
		// and we check if the rollback exceeds the maximum tolerable threshold.
		if s.elapsed-now > maxTolerableRollback {
			return 0, ErrClockRollback
		}
		// wait to next time unit
		for now < s.elapsed {
			time.Sleep(time.Millisecond) // avoid busy-wait CPU spinning
			now = time.Since(s.epoch).Milliseconds() / timeUnit
		}
		s.sequence = 0
	}
	if now > maxTimestamp {
		return 0, ErrOverTimeLimit
	}

	s.elapsed = now

	return (now << timeShift) | (s.node << nodeShift) | (s.sequence), nil
}

// Timestamp returns an int64 unix timestamp in milliseconds of the snowflake ID time.
func Timestamp(id int64) int64 {
	return (id>>timeShift)*timeUnit + epochTimestamp
}

// Node returns an int64 of the snowflake ID node id.
func Node(id int64) int64 {
	return (id & nodeMask) >> nodeShift
}

// Sequence returns an int64 of the snowflake ID sequence number.
func Sequence(id int64) int64 {
	return id & sequenceMask
}

// Compose creates a snowflake ID from its components.
// The time parameter should be the time when the ID was generated.
// The node parameter should be between 0 and 2^s.nodeBits-1 (inclusive).
// The sequence parameter should be between 0 and 2^s.sequenceBits-1 (inclusive).
func Compose(t time.Time, node, sequence int64) (int64, error) {
	elapsed := (t.UnixMilli() - epochTimestamp) / timeUnit
	if elapsed < 0 || elapsed > maxTimestamp {
		return 0, ErrOverTimeLimit
	}
	if node < 0 || node > maxNode {
		return 0, ErrOverNodeLimit
	}
	if sequence < 0 || sequence > maxSequence {
		return 0, ErrOverSequenceLimit
	}
	return (elapsed << timeShift) | (node << nodeShift) | (sequence), nil
}

// Decompose returns a set of snowflake ID parts.
func Decompose(id int64) map[string]int64 {
	timestamp := Timestamp(id)
	node := Node(id)
	sequence := Sequence(id)
	return map[string]int64{
		"id":        id,
		"timestamp": timestamp,
		"node":      node,
		"sequence":  sequence,
	}
}
