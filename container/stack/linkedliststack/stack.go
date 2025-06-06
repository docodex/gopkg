// Package linkedliststack implements a singly linked list stack.
package linkedliststack

import (
	"encoding/json"

	"github.com/docodex/gopkg/jsonx"
)

// node is a node of a linked-list stack.
type node[T any] struct {
	next  *node[T] // Next pointer in the singly-linked-list stack of nodes.
	value T        // The value stored with this node.
}

// Stack represents a singly linked-list stack.
type Stack[T any] struct {
	head node[T] // sentinel stack node, only head.next are used
	len  int     // current stack length excluding the sentinel node
}

// New returns an initialized stack.
func New[T any]() *Stack[T] {
	return new(Stack[T]).init()
}

// init initializes or clears stack s.
func (s *Stack[T]) init() *Stack[T] {
	s.head.next = nil
	s.len = 0
	return s
}

// Len returns the number of nodes of stack s.
// The complexity is O(1).
func (s *Stack[T]) Len() int {
	return s.len
}

// Values returns all values in stack (in LIFO order).
func (s *Stack[T]) Values() []T {
	values := make([]T, s.len)
	for i, x := 0, s.head.next; i < s.len; i, x = i+1, x.next {
		values[i] = x.value
	}
	return values
}

// listValues returns all values in stack (in [Stack.Push] order).
func (s *Stack[T]) listValues() []T {
	values := make([]T, s.len)
	for i, x := s.len-1, s.head.next; i >= 0; i, x = i-1, x.next {
		values[i] = x.value
	}
	return values
}

// String returns the string representation of stack.
// Ref: std fmt.Stringer.
func (s *Stack[T]) String() string {
	values, _ := jsonx.MarshalToString(s.listValues())
	return "LinkedListStack: " + values
}

// MarshalJSON marshals stack into valid JSON.
// Ref: std json.Marshaler.
func (s *Stack[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.listValues())
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
	s.Clear()
	for i := range v {
		s.Push(v[i])
	}
	return nil
}

// Push adds the given value v to the top of stack.
func (s *Stack[T]) Push(v T) {
	s.head.next = &node[T]{
		next:  s.head.next,
		value: v,
	}
	s.len++
}

// Pop removes the top element if exists in stack and returns it.
// The ok result indicates whether such element was removed from stack.
func (s *Stack[T]) Pop() (value T, ok bool) {
	if s.len > 0 {
		x := s.head.next
		value = x.value
		ok = true
		s.head.next = x.next
		x.next = nil
		s.len--
	}
	return
}

// Peek returns the top element if exists in stack without removing it.
// The ok result indicates whether such element was found in stack.
func (s *Stack[T]) Peek() (value T, ok bool) {
	if s.len > 0 {
		value = s.head.next.value
		ok = true
	}
	return
}

// Clear removes all elements in stack.
func (s *Stack[T]) Clear() {
	for x, y := &s.head, s.head.next; y != nil; x, y = y, y.next {
		x.next = nil
	}
	s.init()
}
