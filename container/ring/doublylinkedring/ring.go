// Package doublylinkedring implements a doubly linked circular list.
package doublylinkedring

import "github.com/docodex/gopkg/jsonx"

// A Ring is an element of a doubly-linked circular list, or ring.
// Rings do not have a beginning or end; a pointer to any ring element
// serves as reference to the entire ring. Empty rings are represented
// as nil Ring pointers. The zero value for a Ring is a one-element
// ring with a zero value 'Value'.
type Ring[T any] struct {
	next, prev *Ring[T] // Pointers to next and previous ring elements.
	Value      T        // The value stored with this ring element.
}

// New returns an initialized ring.
func New[T any](v T, rest ...T) *Ring[T] {
	r := &Ring[T]{
		Value: v,
	}
	x := r
	for i := range rest {
		x.next = &Ring[T]{
			prev:  x,
			Value: rest[i],
		}
		x = x.next
	}
	x.next = r
	r.prev = x
	return r
}

func (r *Ring[T]) init() {
	r.next = r
	r.prev = r
}

// Next returns the next ring element. r must not be empty.
func (r *Ring[T]) Next() *Ring[T] {
	if r.next == nil {
		r.init()
	}
	return r.next
}

// Prev returns the previous ring element. r must not be empty.
func (r *Ring[T]) Prev() *Ring[T] {
	if r.prev == nil {
		r.init()
	}
	return r.prev
}

// Move moves n % r.Len() elements backward (n < 0) or forward (n >= 0)
// in the ring and returns that ring element. r must not be empty.
func (r *Ring[T]) Move(n int) *Ring[T] {
	switch {
	case n < 0:
		if r.prev == nil {
			r.init()
		}
		for ; n < 0; n++ {
			r = r.prev
		}
	case n > 0:
		if r.next == nil {
			r.init()
		}
		for ; n > 0; n-- {
			r = r.next
		}
	}
	return r
}

// Link connects ring r with ring s such that r.Next()
// becomes s and returns the original value for r.Next().
// r must not be empty.
//
// If r and s point to the same ring, linking
// them removes the elements between r and s from the ring.
// The removed elements form a subring and the result is a
// reference to that subring (if no elements were removed,
// the result is still the original value for r.Next(),
// and not nil).
//
// If r and s point to different rings, linking
// them creates a single ring with the elements of s inserted
// after r. The result points to the element following the
// last element of s after insertion.
func (r *Ring[T]) Link(s *Ring[T]) *Ring[T] {
	x := r.Next()
	if s != nil {
		y := s.Prev()
		// Note: Cannot use multiple assignment because
		// evaluation order of LHS is not specified.
		r.next = s
		s.prev = r
		x.prev = y
		y.next = x
	}
	return x
}

// Unlink removes n % r.Len() elements from the ring r, starting
// at r.Next(). If n % r.Len() == 0, r remains unchanged.
// The result is the removed subring. r must not be empty.
func (r *Ring[T]) Unlink(n int) *Ring[T] {
	if n <= 0 {
		return nil
	}
	return r.Link(r.Move(n + 1))
}

// Len returns the number of elements of ring r.
// The complexity is O(n).
func (r *Ring[T]) Len() int {
	if r.next == nil {
		r.init()
	}
	n := 1
	for x := r.next; x != r; x = x.next {
		n++
	}
	return n
}

// Values returns all values in ring (start with node r).
func (r *Ring[T]) Values() []T {
	size := r.Len()
	values := make([]T, size)
	for i, x := 0, r; i < size; i, x = i+1, x.next {
		values[i] = x.Value
	}
	return values
}

// String returns the string representation of ring.
// Ref: std fmt.Stringer.
func (r *Ring[T]) String() string {
	values, _ := jsonx.MarshalToString(r.Values())
	return "DoublyLinkedRing: " + values
}

// Add inserts new elements with values v after r.
func (r *Ring[T]) Add(v ...T) *Ring[T] {
	if len(v) == 0 {
		return r
	}
	if r.next == nil {
		r.init()
	}
	x := r
	y := r.next
	for i := range v {
		x.next = &Ring[T]{
			prev:  x,
			Value: v[i],
		}
		x = x.next
	}
	x.next = y
	y.prev = x
	return r.next
}

// Range calls f sequentially for each value v present in ring.
// If f returns false, range stops the iteration.
func (r *Ring[T]) Range(f func(v T) bool) {
	if f == nil {
		return
	}
	if !f(r.Value) {
		return
	}
	if r.next == nil {
		r.init()
	}
	for x := r.next; x != r; x = x.next {
		if !f(x.Value) {
			break
		}
	}
}

// Delete deletes a ring by set all ring elements next and previous pointers to nil.
func Delete[T any](r *Ring[T]) {
	if r == nil {
		return
	}
	if r.next == nil {
		r.init()
	}
	for x := r.next; x != r; {
		n := x.next
		x.next = nil // avoid memory leaks
		x.prev = nil // avoid memory leaks
		x = n
	}
	r.next = nil // avoid memory leaks
	r.prev = nil // avoid memory leaks
}
