// Package snowflake provides a very simple Twitter snowflake generator and parser.
package snowflake

import (
	"errors"
	"sync"
	"time"

	"github.com/docodex/gopkg/internal"
)

const (
	// Epoch is the time since which the snowflake time is defined as the timestamp.
	// The default epoch is set to 2025-01-01 00:00:00 +0000 UTC in milliseconds.
	// You may customize this to set a different epoch for your application.
	// The epoch should be before the current time.
	epoch = 1735689600000 // time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	// TimeBits holds the number of bits to use for timestamp.
	// The timeBits is calculated by 63 - nodeBits - sequenceBits.
	// The default timeBits is set to 42: 63-10-11.
	// The timeBits should be between 36 (inclusive) and 52 (inclusive).
	timeBits = 42

	// NodeBits holds the number of bits to use for Node.
	// NodeBits holds the number of bits to use for node id.
	// The default nodeBits is set to 10.
	// You may customize this to set a different length for your application.
	// The nodeBits should be between 1 (inclusive) and 26 (inclusive).
	// Remember, you have a total (63 - timeBits) bits to share between Node/Sequence.
	nodeBits = 10

	// SequenceBits holds the number of bits to use for sequence number.
	// The default sequenceBits is set to 11.
	// You may customize this to set a different length for your application.
	// The sequenceBits should be between 1 (inclusive) and 26 (inclusive).
	// Remember, you have a total (63 - timeBits) bits to share between Node/Sequence.
	sequenceBits = 11
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
	ErrOverNodeLimit     = errors.New("over the ndoe id limit")
	ErrOverSequenceLimit = errors.New("over the sequence number limit")
	ErrCheckNodeFailed   = errors.New("check node id failed")
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

	elapsed  int64
	node     int64
	sequence int64

	checkNode func(node int64) bool
}

func Default() *Snowflake {
	s := &Snowflake{node: -1}
	s.prepare()
	return s
}

func New(opts ...Option) (*Snowflake, error) {
	s := Default()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	if s.node > maxNode {
		return nil, ErrOverNodeLimit
	}
	if s.checkNode != nil && !s.checkNode(s.node) {
		return nil, ErrCheckNodeFailed
	}
	s.prepare()
	return s, nil
}

func (s *Snowflake) prepare() {
	if s.node == -1 {
		// The default node is set to the lower 8 bits of the private IP address.
		node := int64(internal.Lower8BitPrivateIPv4())
		if node > maxNode {
			// If over node limit, 0 would be used.
			s.node = 0
		} else {
			s.node = node
		}
	}
	// s.sequence = s.maxSequence
}

// Generate creates and returns a unique snowflake ID.
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same node id
func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli() - epoch
	if now == s.elapsed {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// wait to next time unit: for-loop or sleep
			for now <= s.elapsed {
				now = time.Now().UnixMilli() - epoch
			}
		}
	} else {
		s.sequence = 0
	}
	s.elapsed = now

	return (s.elapsed << timeShift) | (s.node << nodeShift) | (s.sequence)
}

// Timestamp returns an int64 unix timestamp in milliseconds of the snowflake ID time.
func Timestamp(id int64) int64 {
	return (id >> timeShift) + epoch
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
	elapsed := t.UnixMilli() - epoch
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
