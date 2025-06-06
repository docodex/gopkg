// Package linkedlistqueue implements a singly linked list queue.
package linkedlistqueue

import (
	"encoding/json"

	"github.com/docodex/gopkg/jsonx"
)

// node is a node of a linked-list queue.
type node[T any] struct {
	next  *node[T] // Next pointer in the singly-linked-list queue of nodes.
	value T        // The value stored with this node.
}

// Queue represents a singly-linked-list queue.
type Queue[T any] struct {
	head node[T]  // sentinel queue node, only head.next are used
	last *node[T] // the last node in queue, or point to head while queue is empty
	len  int      // current queue length excluding the sentinel node
}

// New returns an initialized queue.
func New[T any]() *Queue[T] {
	return new(Queue[T]).init()
}

// init initializes or clears queue q.
func (q *Queue[T]) init() *Queue[T] {
	q.head.next = nil
	q.last = &q.head
	q.len = 0
	return q
}

// Len returns the number of nodes of queue q.
// The complexity is O(1).
func (q *Queue[T]) Len() int {
	return q.len
}

func (q *Queue[T]) Values() []T {
	values := make([]T, q.len)
	for i, x := 0, q.head.next; i < q.len; i, x = i+1, x.next {
		values[i] = x.value
	}
	return values
}

// String returns the string representation of queue.
// Ref: std fmt.Stringer.
func (q *Queue[T]) String() string {
	values, _ := jsonx.MarshalToString(q.Values())
	return "LinkedListQueue: " + values
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
	q.Clear()
	for i := range v {
		q.Enqueue(v[i])
	}
	return nil
}

// Enqueue adds the value v to the end of queue.
func (q *Queue[T]) Enqueue(v T) {
	q.last.next = &node[T]{
		value: v,
	}
	q.last = q.last.next
	q.len++
}

// Dequeue removes the first element if exists in queue and returns it.
// The ok result indicates whether such element was removed from queue.
func (q *Queue[T]) Dequeue() (value T, ok bool) {
	if q.len > 0 {
		x := q.head.next
		value = x.value
		ok = true
		q.head.next = x.next
		x.next = nil
		q.len--
		if q.len == 0 {
			q.last = &q.head
		}
	}
	return
}

// Peek returns the first element if exists in queue without removing it.
// The ok result indicates whether such element was found in queue.
func (q *Queue[T]) Peek() (value T, ok bool) {
	if q.len > 0 {
		value = q.head.next.value
		ok = true
	}
	return
}

// Clear removes all elements in queue.
func (q *Queue[T]) Clear() {
	for x, y := &q.head, q.head.next; y != nil; x, y = y, y.next {
		x.next = nil
	}
	q.init()
}
