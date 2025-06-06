package container

// Container is the base interface for all container data structures to implement.
type Container[T any] interface {
	// Len returns the number of elements of a container.
	Len() int
	// Values returns a slice of all elements of a container.
	Values() []T
	// String returns the string representation of a container.
	String() string
}

// Compare should return a negative number (-1) when a < b, a positive number (1) when
// a > b and zero (0) when a == b or a and b are incomparable in the sense of a strict
// weak ordering.
//
// See https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings.
type Compare[T any] func(a, b T) int

// Less should return true if a is less than b, otherwise, return false.
type Less[T any] func(a, b T) bool
