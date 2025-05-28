package snowflake

import "time"

// Option represents a modification to the default behavior of a Snowflake.
type Option func(s *Snowflake) error

func WithEpoch(epoch time.Time) Option {
	return func(s *Snowflake) error {
		if epoch.After(time.Now()) {
			return ErrInvalidEpochTime
		}
		s.epoch = epoch.UnixMilli()
		return nil
	}
}

func WithNodeBits(nodeBits uint8) Option {
	return func(s *Snowflake) error {
		if nodeBits < 1 || nodeBits > 26 {
			return ErrInvalidNodeBits
		}
		s.nodeBits = nodeBits
		return nil
	}
}

func WithSequenceBits(sequenceBits uint8) Option {
	return func(s *Snowflake) error {
		if sequenceBits < 1 || sequenceBits > 26 {
			return ErrInvalidSequenceBits
		}
		s.sequenceBits = sequenceBits
		return nil
	}
}

func WithNode(node int64) Option {
	return func(s *Snowflake) error {
		if node < 0 {
			return ErrOverNodeLimit
		}
		s.node = node
		return nil
	}
}

func WithCheckNode(checkNode func(node int64) bool) Option {
	return func(s *Snowflake) error {
		if checkNode != nil {
			s.checkNode = checkNode
		}
		return nil
	}
}
