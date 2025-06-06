// Package stack provides an abstract Stack interface.
//
// In computer science, a stack is an abstract data type that serves as a collection of
// elements, with two principal operations: push, which adds an element to the collection,
// and pop, which removes the most recently added element that was not yet removed. The
// order in which elements come off a stack gives rise to its alternative name, LIFO (for
// last in, first out). Additionally, a peek operation may give access to the top without
// modifying the stack.
//
// Reference: https://en.wikipedia.org/wiki/Stack_%28abstract_data_type%29
package stack

import "github.com/docodex/gopkg/container"

type Stack[T any] interface {
	container.Container[T]

	// MarshalJSON marshals stack into valid JSON.
	// Ref: std json.Marshaler.
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals a JSON description of stack.
	// The input can be assumed to be a valid encoding of a JSON value.
	// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
	// Ref: std json.Unmarshaler.
	UnmarshalJSON(data []byte) error

	// Push adds the given value v to the top of stack.
	Push(v T)
	// Pop removes the top element if exists in stack and returns it.
	// The ok result indicates whether such element was removed from stack.
	Pop() (value T, ok bool)
	// Peek returns the top element if exists in stack without removing it.
	// The ok result indicates whether such element was found in stack.
	Peek() (value T, ok bool)

	// Clear removes all elements in stack.
	Clear()
}
