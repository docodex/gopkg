// Package arraystack implements an array stack.
package arraystack

import (
	"encoding/json"

	"github.com/docodex/gopkg/jsonx"
)

// Stack represents an array stack which holds the elements in a slice.
type Stack[T any] struct {
	values []T // current stack elements
}

// New returns an initialized stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{values: nil}
}

// Len returns the number of elements of stack s.
// The complexity is O(1).
func (s *Stack[T]) Len() int {
	return len(s.values)
}

// Values returns all values in stack (in LIFO order).
func (s *Stack[T]) Values() []T {
	values := make([]T, len(s.values))
	for i, j := len(s.values)-1, 0; i >= 0; i, j = i-1, j+1 {
		values[j] = s.values[i]
	}
	return values
}

// String returns the string representation of stack.
// Ref: std fmt.Stringer.
func (s *Stack[T]) String() string {
	values, _ := jsonx.MarshalToString(s.Values())
	return "ArrayStack: " + values
}

// MarshalJSON marshals stack into valid JSON.
// Ref: std json.Marshaler.
func (s *Stack[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.values)
}

// UnmarshalJSON unmarshals a JSON description of stack.
// The input can be assumed to be a valid encoding of a JSON value.
// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
// Ref: std json.Unmarshaler.
func (s *Stack[T]) UnmarshalJSON(data []byte) error {
	var v []T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s.values = v
	return nil
}

const defaultCapacity = 128

// checkAndExpand checks and expands the underlying array if necessary.
func (s *Stack[T]) checkAndExpand(delta int) {
	size := len(s.values) + delta
	if size <= cap(s.values) {
		return
	}
	// expand & migrate
	capacity := max(size<<1, defaultCapacity)
	v := make([]T, 0, capacity)
	v = append(v, s.values...)
	s.values = v
}

// checkAndShrink checks and shrinks the underlying array if necessary.
func (s *Stack[T]) checkAndShrink() {
	if cap(s.values) <= defaultCapacity {
		return
	}
	if len(s.values)<<2 > cap(s.values) {
		return
	}
	// shrink & migrate
	capacity := max(len(s.values)<<1, defaultCapacity)
	v := make([]T, 0, capacity)
	v = append(v, s.values...)
	s.values = v
}

// Push adds the given value v to the top of stack.
func (s *Stack[T]) Push(v T) {
	s.checkAndExpand(1)
	s.values = append(s.values, v)
}

// Pop removes the top element if exists in stack and returns it.
// The ok result indicates whether such element was removed from stack.
func (s *Stack[T]) Pop() (value T, ok bool) {
	if len(s.values) != 0 {
		last := len(s.values) - 1
		value = s.values[last]
		ok = true
		s.values = s.values[:last]
		s.checkAndShrink()
	}
	return
}

// Peek returns the top element if exists in stack without removing it.
// The ok result indicates whether such element was found in stack.
func (s *Stack[T]) Peek() (value T, ok bool) {
	if len(s.values) != 0 {
		value = s.values[len(s.values)-1]
		ok = true
	}
	return
}

// Clear removes all elements in stack.
func (s *Stack[T]) Clear() {
	s.values = nil
}
