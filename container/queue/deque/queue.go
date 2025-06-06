// Package deque implements a double ended queue.
//
// Reference: https://en.wikipedia.org/wiki/Double-ended_queue
package deque

import (
	"encoding/json"

	"github.com/docodex/gopkg/jsonx"
)

// Queue represents a double ended queue which holds the elements in a slice.
type Queue[T any] struct {
	values []T // current queue elements
	first  int // first element index
	tail   int // last element index + 1
}

// New returns an initialized double ended queue.
func New[T any]() *Queue[T] {
	return new(Queue[T]).init()
}

const (
	defaultCapacity  = 128
	defaultInitIndex = defaultCapacity >> 1
)

// init initializes or clears queue q.
func (q *Queue[T]) init() *Queue[T] {
	q.values = make([]T, defaultCapacity)
	q.first = defaultInitIndex
	q.tail = defaultInitIndex
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
	return "DoubleEndedQueue: " + values
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

// insert inserts the values v to position i in queue: prepend v to queue if i equals to
// q.first; append v to queue if i equals to q.tail; or do nothing if i is not equal to
// q.first, also not equal to q.tail.
// At the same time, insert checks and expands the underlying array or moves queue
// elements if necessary.
func (q *Queue[T]) insert(i int, v ...T) {
	switch i {
	case q.first:
		first := q.first - len(v)
		if first >= 0 {
			// prepend
			copy(q.values[first:q.first], v)
			q.first = first
			return
		}
	case q.tail:
		tail := q.tail + len(v)
		if tail <= cap(q.values) {
			// append
			copy(q.values[q.tail:tail], v)
			q.tail = tail
			return
		}
	default:
		// invalid position i, do nothing
		return
	}
	s1 := q.Len()
	s2 := s1 + len(v)
	capacity := max(s2<<1, defaultCapacity)
	if capacity > cap(q.values) {
		// expand & migrate
		v1 := make([]T, capacity)
		first := (capacity - s2) >> 1
		tail := first + s2
		if i == q.first {
			// prepend
			j := first + len(v)
			copy(v1[first:j], v)
			copy(v1[j:tail], q.values[q.first:q.tail])
		} else {
			// append
			j := first + s1
			copy(v1[first:j], q.values[q.first:q.tail])
			copy(v1[j:tail], v)
		}
		q.values = v1
		q.first = first
		q.tail = tail
	} else {
		// move
		first := (cap(q.values) - s2) >> 1
		tail := first + s2
		if i == q.first {
			// prepend
			j := first + len(v)
			copy(q.values[j:tail], q.values[q.first:q.tail])
			copy(q.values[first:j], v)
		} else {
			// append
			j := first + s1
			copy(q.values[first:j], q.values[q.first:q.tail])
			copy(q.values[j:tail], v)
		}
		q.first = first
		q.tail = tail
	}
}

// EnqueueFront adds the value v to the front of queue.
func (q *Queue[T]) EnqueueFront(v T) {
	q.insert(q.first, v)
}

// EnqueueBack adds the value v to the back of queue.
func (q *Queue[T]) EnqueueBack(v T) {
	q.insert(q.tail, v)
}

// checkAndShrink checks and shrinks the underlying array if necessary.
func (q *Queue[T]) checkAndShrink() {
	if cap(q.values) <= defaultCapacity {
		return
	}
	size := q.Len()
	if size<<2 > cap(q.values) {
		return
	}
	// shrink & migrate
	capacity := max(size<<1, defaultCapacity)
	v := make([]T, capacity)
	first := (capacity - size) >> 1
	tail := first + size
	copy(v[first:tail], q.values[q.first:q.tail])
	q.values = v
	q.first = first
	q.tail = tail
}

// DequeueFront removes the first element if exists in queue and returns it.
// The ok result indicates whether such element was removed from queue.
func (q *Queue[T]) DequeueFront() (value T, ok bool) {
	if q.first < q.tail {
		value = q.values[q.first]
		ok = true
		q.first++
		q.checkAndShrink()
	}
	return
}

// DequeueBack removes the last element if exists in queue and returns it.
// The ok result indicates whether such element was removed from queue.
func (q *Queue[T]) DequeueBack() (value T, ok bool) {
	if q.first < q.tail {
		q.tail--
		value = q.values[q.tail]
		ok = true
		q.checkAndShrink()
	}
	return
}

// PeekFront returns the first element if exists in queue without removing it.
// The ok result indicates whether element was found in queue.
func (q *Queue[T]) PeekFront() (value T, ok bool) {
	if q.first < q.tail {
		value = q.values[q.first]
		ok = true
	}
	return
}

// PeekBack returns the last element if exists in queue without removing it.
// The ok result indicates whether element was found in queue.
func (q *Queue[T]) PeekBack() (value T, ok bool) {
	if q.first < q.tail {
		value = q.values[q.tail-1]
		ok = true
	}
	return
}

// Clear removes all elements in queue.
func (q *Queue[T]) Clear() {
	q.init()
}
