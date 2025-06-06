// Package arrayqueue implements an array queue.
package arrayqueue

import (
	"encoding/json"

	"github.com/docodex/gopkg/jsonx"
)

// Queue represents an array queue which holds the elements in a slice.
type Queue[T any] struct {
	values []T // current queue elements
	first  int // first element index
	tail   int // last element index + 1
}

// New returns an initialized queue.
func New[T any]() *Queue[T] {
	return new(Queue[T]).init()
}

// init initializes or clears queue q.
func (q *Queue[T]) init() *Queue[T] {
	q.values = nil
	q.first = 0
	q.tail = 0
	return q
}

// Len returns the number of elements of queue q.
// The complexity is O(1).
func (q *Queue[T]) Len() int {
	return q.tail - q.first
}

// Values returns all values in queue (in FIFO order).
func (q *Queue[T]) Values() []T {
	values := make([]T, q.Len())
	copy(values, q.values[q.first:q.tail])
	return values
}

// String returns the string representation of queue.
// Ref: std fmt.Stringer.
func (q *Queue[T]) String() string {
	values, _ := jsonx.MarshalToString(q.values[q.first:q.tail])
	return "ArrayQueue: " + values
}

// MarshalJSON marshals queue into valid JSON.
// Ref: std json.Marshaler.
func (q *Queue[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.values[q.first:q.tail])
}

// UnmarshalJSON unmarshals a JSON description of queue.
// The input can be assumed to be a valid encoding of a JSON value.
// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
// Ref: std json.Unmarshaler.
func (q *Queue[T]) UnmarshalJSON(data []byte) error {
	var v []T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	q.values = v
	q.first = 0
	q.tail = len(v)
	return nil
}

const defaultCapacity = 128

// checkAndExpandOrMove checks and expands the underlying array or moves queue
// elements if necessary.
func (q *Queue[T]) checkAndExpandOrMove(delta int) {
	if q.tail+delta <= cap(q.values) {
		return
	}
	size := q.Len()
	capacity := max((size+delta)<<1, defaultCapacity)
	if capacity > cap(q.values) {
		// expand & migrate
		v := make([]T, capacity)
		copy(v[:size], q.values[q.first:q.tail])
		q.values = v
	} else {
		// move
		copy(q.values[:size], q.values[q.first:q.tail])
	}
	q.first = 0
	q.tail = size
}

// checkAndShrink checks and shrinks the underlying array if necessary.
func (q *Queue[T]) checkAndShrink() {
	if q.tail <= defaultCapacity {
		return
	}
	size := q.Len()
	if size<<2 > cap(q.values) {
		return
	}
	// shrink & migrate
	v := make([]T, max(size<<1, defaultCapacity))
	copy(v[:size], q.values[q.first:q.tail])
	q.values = v
	q.first = 0
	q.tail = size
}

// Enqueue adds the value v to the end of queue.
func (q *Queue[T]) Enqueue(v T) {
	q.checkAndExpandOrMove(1)
	q.values[q.tail] = v
	q.tail++
}

// Dequeue removes the first element if exists in queue and returns it.
// The ok result indicates whether such element was removed from queue.
func (q *Queue[T]) Dequeue() (value T, ok bool) {
	if q.first < q.tail {
		value = q.values[q.first]
		ok = true
		q.first++
		q.checkAndShrink()
	}
	return
}

// Peek returns the first element if exists in queue without removing it.
// The ok result indicates whether such element was found in queue.
func (q *Queue[T]) Peek() (value T, ok bool) {
	if q.first < q.tail {
		value = q.values[q.first]
		ok = true
	}
	return
}

// Clear removes all elements in queue.
func (q *Queue[T]) Clear() {
	q.init()
}
