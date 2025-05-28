// Package snowflake provides a very simple Twitter snowflake generator and parser.
package snowflake

import (
	"errors"
	"sync"
	"time"

	"github.com/docodex/gopkg/internal"
)

const (
	defaultEpoch = 1735689600000 // time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	defaultTimeUnit = 10 // 10 msec, in unit of msec

	defaultTimeBits     = 39
	defaultNodeBits     = 11
	defaultSequenceBits = 13
)

var (
	ErrInvalidEpochTime    = errors.New("invalid epoch time")
	ErrInvalidTimeUnit     = errors.New("invalid time unit")
	ErrInvalidTimeBits     = errors.New("bit length for timestamp should be between 36 and 52")
	ErrInvalidNodeBits     = errors.New("bit length for node id should be between 1 and 26")
	ErrInvalidSequenceBits = errors.New("bit length for sequence number should be between 1 and 26")
	ErrOverTimeLimit       = errors.New("over the timestamp limit")
	ErrOverNodeLimit       = errors.New("over the ndoe id limit")
	ErrOverSequenceLimit   = errors.New("over the sequence number limit")
	ErrCheckNodeFailed     = errors.New("check node id failed")
	ErrNoPrivateAddress    = errors.New("no private ip address")
)

// Snowflake is a distributed unique ID generator inspired by twitter snowflake.
// By default, a Snowflake ID is composed of
// - 39 bits for time in units of 10 msec
// - 11 bits for a node id
// - 13 bits for a sequence number
//
// Epoch is the time since which the snowflake time is defined as the timestamp.
// The default epoch is set to 2025-01-01 00:00:00 +0000 UTC in milliseconds.
// You may customize this to set a different epoch for your application.
// The epoch should be before the current time.
// Otherwise, an error would be returned.
//
// TimeUnit is the time unit of snowflake, in unit of msec.
// The default timeUnit is 10 msec.
// You may customize this to set a different epoch for your application.
// TimeUnit should be between 1 msec (inclusive) and 1 sec (inclusive).
//
// NodeBits holds the number of bits to use for node id.
// The default nodeBits is set to 11.
// You may customize this to set a different length for your application.
// The nodeBits should be between 1 (inclusive) and 26 (inclusive).
// Otherwise, an error would be returned.
//
// SequenceBits holds the number of bits to use for sequence number.
// The default sequenceBits is set to 13.
// You may customize this to set a different length for your application.
// The sequenceBits should be between 1 (inclusive) and 26 (inclusive).
// Otherwise, an error would be returned.
//
// TimeBits holds the number of bits to use for timestamp.
// The timeBits is calculated by 63 - nodeBits - sequenceBits.
// The default timeBits is set to 39: 63-11-13.
// The timeBits should be between 36 (inclusive) and 52 (inclusive).
// Otherwise, an error would be returned.
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

	epoch        int64
	timeUnit     int64
	timeBits     uint8
	nodeBits     uint8
	sequenceBits uint8

	timeShift    uint8
	maxTimestamp int64
	nodeShift    uint8
	maxNode      int64
	nodeMask     int64
	maxSequence  int64
	sequenceMask int64

	elapsed  int64
	node     int64
	sequence int64

	checkNode func(node int64) bool
}

func Default() *Snowflake {
	s := &Snowflake{
		epoch:        defaultEpoch,
		timeUnit:     defaultTimeUnit,
		timeBits:     defaultTimeBits,
		nodeBits:     defaultNodeBits,
		sequenceBits: defaultSequenceBits,
		node:         -1,
	}
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
	s.timeBits = 63 - s.nodeBits - s.sequenceBits
	if s.timeBits < 36 || s.timeBits > 52 {
		return nil, ErrInvalidTimeBits
	}
	if s.node >= (1 << s.nodeBits) {
		return nil, ErrOverNodeLimit
	}
	if s.checkNode != nil && !s.checkNode(s.node) {
		return nil, ErrCheckNodeFailed
	}

	s.prepare()

	return s, nil
}

func (s *Snowflake) prepare() {
	s.timeShift = s.nodeBits + s.sequenceBits
	s.maxTimestamp = -1 ^ (-1 << s.timeBits) // (1 << s.timeBits) - 1

	s.nodeShift = s.sequenceBits
	s.maxNode = -1 ^ (-1 << s.nodeBits) // (1 << s.nodeBits) - 1
	s.nodeMask = s.maxNode << s.nodeShift

	s.maxSequence = -1 ^ (-1 << s.sequenceBits) // (1 << s.sequenceBits) - 1
	s.sequenceMask = s.maxSequence

	if s.node == -1 {
		// The default node is set to the lower 8 bits of the private IP address.
		node := int64(internal.Lower8BitPrivateIPv4())
		if node > s.maxNode {
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
func (s *Snowflake) Generate() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.toInternalTimestamp(time.Now())
	if now == s.elapsed {
		s.sequence = (s.sequence + 1) & s.sequenceMask
		if s.sequence == 0 {
			// wait to next time unit: for-loop or sleep
			for now <= s.elapsed {
				now = s.toInternalTimestamp(time.Now())
			}
		}
	} else {
		s.sequence = 0
	}
	s.elapsed = now
	if s.elapsed > s.maxTimestamp {
		return 0, ErrOverTimeLimit
	}

	return (s.elapsed << s.timeShift) | (s.node << s.nodeShift) | (s.sequence), nil
}

func (s *Snowflake) toInternalTimestamp(t time.Time) int64 {
	return (t.UnixMilli() - s.epoch) / s.timeUnit
}

// Timestamp returns an int64 unix timestamp in milliseconds of the snowflake ID time.
func (s *Snowflake) Timestamp(id int64) int64 {
	return (id>>s.timeShift)*s.timeUnit + s.epoch
}

// Node returns an int64 of the snowflake ID node id.
func (s *Snowflake) Node(id int64) int64 {
	return (id & s.nodeMask) >> s.nodeShift
}

// Sequence returns an int64 of the snowflake ID sequence number.
func (s *Snowflake) Sequence(id int64) int64 {
	return id & s.sequenceMask
}

// Compose creates a snowflake ID from its components.
// The time parameter should be the time when the ID was generated.
// The node parameter should be between 0 and 2^s.nodeBits-1 (inclusive).
// The sequence parameter should be between 0 and 2^s.sequenceBits-1 (inclusive).
func (s *Snowflake) Compose(t time.Time, node, sequence int64) (int64, error) {
	elapsed := s.toInternalTimestamp(t)
	if elapsed < 0 || elapsed > s.maxTimestamp {
		return 0, ErrOverTimeLimit
	}
	if node < 0 || node > s.maxNode {
		return 0, ErrOverNodeLimit
	}
	if sequence < 0 || sequence > s.maxSequence {
		return 0, ErrOverSequenceLimit
	}
	return (elapsed << s.timeShift) | (node << s.nodeShift) | (sequence), nil
}

// Decompose returns a set of snowflake ID parts.
func (s *Snowflake) Decompose(id int64) map[string]int64 {
	timestamp := s.Timestamp(id)
	node := s.Node(id)
	sequence := s.Sequence(id)
	return map[string]int64{
		"id":        id,
		"timestamp": timestamp,
		"node":      node,
		"sequence":  sequence,
	}
}
