// Package container defines the shared interfaces for all container data
// structures in this module.
//
// # Zero-value policy
//
// Unless a type's godoc explicitly states otherwise, concrete containers
// defined under github.com/docodex/gopkg/container are NOT usable in their
// zero value and MUST be constructed with their New, NewFunc, or
// NewWith*-family constructors. Calling methods on a zero-value
// List/Stack/Queue/Set/Map/BidiMap will panic (nil-slice indexing,
// nil-map assignment, or division by zero for circularqueue).
//
// The ring types (container/ring/doublylinkedring, singlylinkedring) are
// the exception: they support zero-value use and lazily self-initialize.
//
// # Concurrency
//
// Containers are NOT safe for concurrent use by default. Types that offer
// a *WithLock constructor (hashset, treeset, hashmap, treemap, hashbidimap,
// treebidimap) take a sync.RWMutex and are safe for concurrent use; all
// other types require the caller to synchronize.
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
