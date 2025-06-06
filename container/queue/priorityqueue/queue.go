// Package priorityqueue implements a priority queue (binary heap).
//
// In computer science, a priority queue is an abstract data type similar to a regular
// queue or stack abstract data type. In a priority queue, each element has an associated
// priority, which determines its order of service. Priority queue serves highest priority
// items first. Priority values have to be instances of an ordered data type, and higher
// priority can be given either to the lesser or to the greater values with respect to the
// given order relation. For example, in Java standard library, PriorityQueue's the least
// elements with respect to the order have the highest priority. This implementation
// detail is without much practical significance, since passing to the opposite order
// relation turns the least values into the greatest, and vice versa.
//
// Reference: https://en.wikipedia.org/wiki/Priority_queue
package priorityqueue

import (
	"cmp"
	"encoding/json"

	"github.com/docodex/gopkg/container"
	"github.com/docodex/gopkg/jsonx"
)

// Element is an element of a priority queue.
type Element[T any] struct {
	Value T   // the value stored with this element
	index int // the index of this element in queue (maintained by queue)
}

// Index returns the index of this element in queue.
func (e *Element[T]) Index() int {
	return e.index
}

// Queue represents an priority queue which holds the elements in a slice.
type Queue[T any] struct {
	elements []*Element[T]     // current queue elements
	less     container.Less[T] // function to compare queue elements
}

// New returns an initialized priority queue with [cmp.Less] as the less function and the given
// values v added.
func New[T cmp.Ordered](v ...T) *Queue[T] {
	return NewFunc(func(a, b T) bool {
		return cmp.Less(a, b)
	}, v...)
}

// NewFunc returns an initialized priority queue with the given function less as the less function
// and the given values v added.
func NewFunc[T any](less container.Less[T], v ...T) *Queue[T] {
	if less == nil {
		less = func(a, b T) bool {
			// just to cover nil less error
			return false
		}
	}
	q := &Queue[T]{
		elements: nil,
		less:     less,
	}
	if len(v) != 0 {
		elements := make([]*Element[T], len(v))
		for i := range v {
			elements[i] = &Element[T]{
				Value: v[i],
				index: i,
			}
		}
		q.elements = elements
		q.init()
	}
	return q
}

// init shift elements in queue to satisfy the property that each node is the minimum-valued
// node in its subtree.
// The complexity is O(n) where n = h.Len().
func (q *Queue[T]) init() {
	// heapify all nodes except leaves
	// for i := range values {
	// 	q.Enqueue(values[i])
	// }
	// the way below is more efficient than enqueue values to queue one by one
	// shift values in queue to satisfy the properties of priority queue.
	for i := ((len(q.elements) - 1) - 1) >> 1; i >= 0; i-- {
		q.shiftDown(i)
	}
}

// swap swaps the elements with indices i and j
func (q *Queue[T]) swap(i, j int) {
	q.elements[i].index = j
	q.elements[j].index = i
	q.elements[i], q.elements[j] = q.elements[j], q.elements[i]
}

// shiftUp shift the element of index i up if necessary.
func (q *Queue[T]) shiftUp(i int) {
	for {
		// parent
		p := (i - 1) >> 1
		// if p is invalid or value of index i is not less than value of index p, break
		if p == i || p < 0 || !q.less(q.elements[i].Value, q.elements[p].Value) {
			break
		}
		// swap values of indices i and p
		q.swap(i, p)
		// loop continue
		i = p
	}
}

// shiftDown shift the element of index i down if necessary, and return true if the shift
// operation done once or more, or return false.
func (q *Queue[T]) shiftDown(i int) bool {
	p := i
	for {
		j := p<<1 + 1 // left child
		// if j is invalid (j < 0 while int overflow), break
		if j >= len(q.elements) || j < 0 {
			break
		}
		if k := j + 1; k < len(q.elements) && q.less(q.elements[k].Value, q.elements[j].Value) {
			j = k // right child: 2*i + 2
		}
		// if value of index j is not less than value of index p, break
		if !q.less(q.elements[j].Value, q.elements[p].Value) {
			break
		}
		// swap values of indices j and p
		q.swap(j, p)
		// loop continue
		p = j
	}
	return p != i
}

// Len returns the number of elements of queue q.
// The complexity is O(1).
func (q *Queue[T]) Len() int {
	return len(q.elements)
}

// Values returns all values in queue (in [Queue.Dequeue] order).
func (q *Queue[T]) Values() []T {
	e1 := make([]*Element[T], 0, len(q.elements))
	for _, e := range q.elements {
		e1 = append(e1, &Element[T]{
			Value: e.Value,
			index: e.index,
		})
	}
	q1 := &Queue[T]{
		elements: e1,
		less:     q.less,
	}
	values := make([]T, 0, len(q.elements))
	for range q.elements {
		v, _ := q1.Dequeue()
		values = append(values, v)
	}
	return values
}

// String returns the string representation of queue.
// Ref: std fmt.Stringer.
func (q *Queue[T]) String() string {
	values, _ := jsonx.MarshalToString(q.Values())
	return "PriorityQueue: " + values
}

// MarshalJSON marshals queue into valid JSON.
// Ref: std json.Marshaler.
func (q *Queue[T]) MarshalJSON() ([]byte, error) {
	values := make([]T, 0, len(q.elements))
	for _, e := range q.elements {
		values = append(values, e.Value)
	}
	return json.Marshal(values)
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
	if len(v) != 0 {
		elements := make([]*Element[T], len(v))
		for i := range v {
			elements[i] = &Element[T]{
				Value: v[i],
				index: i,
			}
		}
		q.elements = elements
		q.init()
	}
	return nil
}

const defaultCapacity = 128

// checkAndExpand checks and expands the underlying array if necessary.
func (q *Queue[T]) checkAndExpand(delta int) {
	size := len(q.elements) + delta
	if size <= cap(q.elements) {
		return
	}
	// expand & migrate
	capacity := max(size<<1, defaultCapacity)
	v := make([]*Element[T], 0, capacity)
	v = append(v, q.elements...)
	q.elements = v
}

// checkAndShrink checks and shrinks the underlying array if necessary.
func (q *Queue[T]) checkAndShrink() {
	if cap(q.elements) <= defaultCapacity {
		return
	}
	if len(q.elements)<<2 > cap(q.elements) {
		return
	}
	// shrink & migrate
	capacity := max(len(q.elements)<<1, defaultCapacity)
	v := make([]*Element[T], 0, capacity)
	v = append(v, q.elements...)
	q.elements = v
}

// Enqueue adds the value v to the end of queue.
func (q *Queue[T]) Enqueue(v T) *Element[T] {
	q.checkAndExpand(1)
	e := &Element[T]{
		Value: v,
		index: len(q.elements),
	}
	q.elements = append(q.elements, e)
	q.shiftUp(e.index)
	return e
}

// Dequeue removes the first element if exists in queue and returns it.
// The ok result indicates whether such element was removed from queue.
func (q *Queue[T]) Dequeue() (value T, ok bool) {
	if len(q.elements) != 0 {
		n := len(q.elements) - 1
		q.swap(0, n)
		value = q.elements[n].Value
		ok = true
		q.elements = q.elements[:n]
		q.shiftDown(0)
		q.checkAndShrink()
	}
	return
}

// Peek returns the first element if exists in queue without removing it.
// The ok result indicates whether such element was found in queue.
func (q *Queue[T]) Peek() (value T, ok bool) {
	if len(q.elements) != 0 {
		value = q.elements[0].Value
		ok = true
	}
	return
}

// Clear removes all elements in queue.
func (q *Queue[T]) Clear() {
	q.elements = nil
}

// Elements returns the underlying elements slice of queue.
func (q *Queue[T]) Elements() []*Element[T] {
	return q.elements
}

// Remove removes and returns the element at index i from queue.
// The complexity is O(log n) where n = h.Len().
func (q *Queue[T]) Remove(i int) (value T, ok bool) {
	if i < 0 || i >= len(q.elements) {
		return
	}
	n := len(q.elements) - 1
	if i != n {
		q.swap(i, n)
	}
	value = q.elements[n].Value
	ok = true
	q.elements = q.elements[:n]
	if i != n && !q.shiftDown(i) {
		q.shiftUp(i)
	}
	q.checkAndShrink()
	return
}

// Fix re-establishes queue ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling [Queue.Remove] followed by a [Queue.Enqueue] of the
// new value.
// The complexity is O(log n) where n = h.Len().
func (q *Queue[T]) Fix(i int) {
	if i < 0 || i >= len(q.elements) {
		return
	}
	if !q.shiftDown(i) {
		q.shiftUp(i)
	}
}

// Update updates the element value to v at index i, and re-establishes heap ordering.
// [Queue.Update] is equivalent to, but less expensive than, calling [Queue.Remove]
// followed by a [Queue.Enqueue] of the new value.
// The complexity is O(log n) where n = h.Len().
func (q *Queue[T]) Update(i int, v T) {
	if i < 0 || i >= len(q.elements) {
		return
	}
	q.elements[i].Value = v
	if !q.shiftDown(i) {
		q.shiftUp(i)
	}
}
