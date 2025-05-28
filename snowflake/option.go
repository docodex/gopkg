package snowflake

// Option represents a modification to the default behavior of a Snowflake.
type Option func(s *Snowflake) error

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
