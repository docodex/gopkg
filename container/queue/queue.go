// Package queue provides an abstract Queue interface.
//
// In computer science, a queue is a collection of entities that are maintained in
// a sequence and can be modified by the addition of entities at one end of the
// sequence and the removal of entities from the other end of the sequence. By
// convention, the end of the sequence at which elements are added is called the
// back, tail, or rear of queue, and the end at which elements are removed is
// called the head or front of queue, analogously to the words used when people
// line up to wait for goods or services.
// The operation of adding an element to the rear of queue is known as enqueue, and
// the operation of removing an element from the front is known as dequeue. Other
// operations may also be allowed, often including a peek or front operation that
// returns the value of the next element to be dequeued without remove it.
//
// Reference: https://en.wikipedia.org/wiki/Queue_(abstract_data_type)
package queue

import "github.com/docodex/gopkg/container"

type Queue[T any] interface {
	container.Container[T]

	// MarshalJSON marshals queue into valid JSON.
	// Ref: std json.Marshaler.
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals a JSON description of queue.
	// The input can be assumed to be a valid encoding of a JSON value.
	// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
	// Ref: std json.Unmarshaler.
	UnmarshalJSON(data []byte) error

	// Dequeue removes the first element if exists in queue and returns it.
	// The ok result indicates whether such element was removed from queue.
	Dequeue() (value T, ok bool)
	// Peek returns the first element if exists in queue without removing it.
	// The ok result indicates whether such element was found in queue.
	Peek() (value T, ok bool)

	// Clear removes all elements in queue.
	Clear()
}
