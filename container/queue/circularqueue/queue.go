// Package circularqueue implements a circular buffer.
//
// In computer science, a circular buffer, circular queue, cyclic buffer or ring buffer
// is a data structure that uses a single, fixed-size buffer as if it were connected
// end-to-end. This structure lends itself easily to buffering data streams.
//
// Reference: https://en.wikipedia.org/wiki/Circular_buffer
package circularqueue

import (
	"encoding/json"

	"github.com/docodex/gopkg/jsonx"
)

// Queue represents a circular queue which holds the elements in a slice.
type Queue[T any] struct {
	values []T // current queue elements
	first  int // first element index
	tail   int // next of last element index
	len    int // current queue length
	cap    int // current queue capacity, cannot be changed after init
}

// New returns an initialized circular queue with the given capacity.
func New[T any](capacity int) *Queue[T] {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}
	return new(Queue[T]).init(capacity)
}

// init initializes or clears queue q.
func (q *Queue[T]) init(capacity int) *Queue[T] {
	q.values = make([]T, capacity)
	q.first = 0
	q.tail = 0
	q.len = 0
	q.cap = capacity
	return q
}

// Empty checks if a queue is empty or not
func (q *Queue[T]) Empty() bool {
	return q.len == 0
}

// Full checks if a queue is full or not
func (q *Queue[T]) Full() bool {
	return q.len == q.cap
}

// Len returns the number of elements of queue q.
// The complexity is O(1).
func (q *Queue[T]) Len() int {
	return q.len
}

// Values returns all values in queue (in FIFO order).
func (q *Queue[T]) Values() []T {
	if q.Empty() {
		return nil
	}
	values := make([]T, 0, q.len)
	if q.first < q.tail {
		values = append(values, q.values[q.first:q.tail]...)
		return values
	}
	values = append(values, q.values[q.first:]...)
	values = append(values, q.values[:q.tail]...)
	return values
}

// String returns the string representation of queue.
// Ref: std fmt.Stringer.
func (q *Queue[T]) String() string {
	values, _ := jsonx.MarshalToString(q.Values())
	return "CircularQueue: " + values
}

// MarshalJSON marshals queue into valid JSON.
// Ref: std json.Marshaler.
func (q *Queue[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.Values())
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
	q.init(max(q.cap, len(v)))
	copy(q.values, v)
	q.first = 0
	q.tail = len(v)
	q.len = len(v)
	return nil
}

// Enqueue adds the value v to the end of queue and return true.
// If queue is full, Enqueue do nothing and return false.
func (q *Queue[T]) Enqueue(v T) bool {
	if q.Full() {
		return false
	}
	q.values[q.tail] = v
	q.tail = (q.tail + 1) % q.cap
	q.len++
	return true
}

// Dequeue removes the first element if exists in queue and returns it.
// The ok result indicates whether such element was removed from queue.
func (q *Queue[T]) Dequeue() (value T, ok bool) {
	if q.Empty() {
		return
	}
	value = q.values[q.first]
	ok = true
	q.first = (q.first + 1) % q.cap
	q.len--
	return
}

// Peek returns the first element if exists in queue without removing it.
// The ok result indicates whether such element was found in queue.
func (q *Queue[T]) Peek() (value T, ok bool) {
	if q.Empty() {
		return
	}
	value = q.values[q.first]
	ok = true
	return
}

// Clear removes all elements in queue.
func (q *Queue[T]) Clear() {
	q.init(q.cap)
}
