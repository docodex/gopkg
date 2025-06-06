// Package binaryheap implements a binary heap which is a binary tree with the
// property that each node is the minimum-valued node in its subtree.
//
// The minimum element in the tree is the root, at index 0.
//
// A heap is a common way to implement a priority queue. To build a priority
// queue, implement a [container.Less] method on T, so [Heap.Push] adds items
// while [Heap.Pop] removes the highest-priority item from queue.
//
// The Less method on T defines this heap as either min or max heap.
//
// Reference: http://en.wikipedia.org/wiki/Binary_heap
package binaryheap

import (
	"cmp"
	"encoding/json"

	"github.com/docodex/gopkg/container"
	"github.com/docodex/gopkg/jsonx"
)

// Heap represents an binary heap which holds the elements in a slice.
type Heap[T any] struct {
	values []T               // current heap elements
	less   container.Less[T] // function to compare heap elements
}

// New returns an initialized heap with [cmp.Less] as the less function and the given values v
// added.
func New[T cmp.Ordered](v ...T) *Heap[T] {
	return NewFunc(func(a, b T) bool {
		return cmp.Less(a, b)
	}, v...)
}

// NewFunc returns an initialized heap with the given function less as the less function and the
// given values v added.
func NewFunc[T any](less container.Less[T], v ...T) *Heap[T] {
	if less == nil {
		less = func(a, b T) bool {
			// just to cover nil less error
			return false
		}
	}
	h := &Heap[T]{
		values: v,
		less:   less,
	}
	h.init()
	return h
}

// init shift values in heap to satisfy the property that each node is the minimum-valued
// node in its subtree.
// The complexity is O(n) where n = h.Len().
func (h *Heap[T]) init() {
	// heapify all nodes except leaves
	// this way is more efficient than push values to heap one by one
	for i := (h.parent(len(h.values) - 1)); i >= 0; i-- {
		h.shiftDown(i)
	}
}

// left returns the left child index
func (h *Heap[T]) left(i int) int {
	return (i << 1) + 1
}

// right returns the right child index
func (h *Heap[T]) right(i int) int {
	return (i << 1) + 2
}

// parent returns the parent index
func (h *Heap[T]) parent(i int) int {
	return (i - 1) >> 1
}

// swap swaps the elements with indices i and j
func (h *Heap[T]) swap(i, j int) {
	h.values[i], h.values[j] = h.values[j], h.values[i]
}

// shiftUp shift the value of index i up if necessary.
func (h *Heap[T]) shiftUp(i int) {
	for {
		p := h.parent(i)
		// p is invalid or no further heapification needed, break
		if p == i || p < 0 || !h.less(h.values[i], h.values[p]) {
			break
		}
		// swap the values of indices i and p
		h.swap(i, p)
		// loop upwards heapification
		i = p
	}
}

// shiftDown shift the value of index i down if necessary, and return true if the shift
// operation done once or more, or return false.
func (h *Heap[T]) shiftDown(i int) bool {
	p := i
	for {
		j := h.left(p) // left child: 2*i + 1
		// if j is invalid (j < 0 while int overflow), break
		if j >= len(h.values) || j < 0 {
			break
		}
		if k := j + 1; k < len(h.values) && h.less(h.values[k], h.values[j]) {
			j = k // right child: 2*i + 2
		}
		// no further heapification needed, break
		if !h.less(h.values[j], h.values[p]) {
			break
		}
		h.swap(p, j)
		// loop downwards heapification
		p = j
	}
	return p != i
}

// Len returns the number of elements of heap h.
// The complexity is O(1).
func (h *Heap[T]) Len() int {
	return len(h.values)
}

// Values returns all values in heap (in [Heap.Pop] order).
func (h *Heap[T]) Values() []T {
	v1 := make([]T, len(h.values))
	copy(v1, h.values)
	h1 := &Heap[T]{
		values: v1,
		less:   h.less,
	}
	values := make([]T, 0, len(h.values))
	for range h.values {
		v, _ := h1.Pop()
		values = append(values, v)
	}
	return values
}

// String returns the string representation of heap.
// Ref: std fmt.Stringer.
func (h *Heap[T]) String() string {
	values, _ := jsonx.MarshalToString(h.Values())
	return "BinaryHeap: " + values
}

// MarshalJSON marshals heap into valid JSON.
// Ref: std json.Marshaler.
func (h *Heap[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.values)
}

// UnmarshalJSON unmarshals a JSON description of heap.
// The input can be assumed to be a valid encoding of a JSON value.
// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
// Ref: std json.Unmarshaler.
func (h *Heap[T]) UnmarshalJSON(data []byte) error {
	var v []T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	h.values = v
	h.init()
	return nil
}

const defaultCapacity = 128

// checkAndExpand checks and expands the underlying array if necessary.
func (h *Heap[T]) checkAndExpand(delta int) {
	size := len(h.values) + delta
	if size <= cap(h.values) {
		return
	}
	// expand & migrate
	capacity := max(size<<1, defaultCapacity)
	v := make([]T, 0, capacity)
	v = append(v, h.values...)
	h.values = v
}

// checkAndShrink checks and shrinks the underlying array if necessary.
func (h *Heap[T]) checkAndShrink() {
	if cap(h.values) <= defaultCapacity {
		return
	}
	if len(h.values)<<2 > cap(h.values) {
		return
	}
	// shrink & migrate
	capacity := max(len(h.values)<<1, defaultCapacity)
	v := make([]T, 0, capacity)
	v = append(v, h.values...)
	h.values = v
}

// Push adds the given value v to heap.
func (h *Heap[T]) Push(v T) {
	h.checkAndExpand(1)
	h.values = append(h.values, v)
	h.shiftUp(len(h.values) - 1)
}

// Pop removes the top element if exists in heap and returns it.
// The ok result indicates whether such element was removed from heap.
func (h *Heap[T]) Pop() (value T, ok bool) {
	if len(h.values) != 0 {
		n := len(h.values) - 1
		h.swap(0, n)
		value = h.values[n]
		ok = true
		h.values = h.values[:n]
		h.shiftDown(0)
		h.checkAndShrink()
	}
	return
}

// Peek returns the top element if exists in heap without removing it.
// The ok result indicates whether such element was found in heap.
func (h *Heap[T]) Peek() (value T, ok bool) {
	if len(h.values) != 0 {
		value = h.values[0]
		ok = true
	}
	return
}

// Clear removes all elements in heap.
func (h *Heap[T]) Clear() {
	h.values = nil
}

// Elements returns the underlying elements slice of heap.
// Note: Do not change the index of any element because index must be maintained by the heap.
// You can update any element value, and call [Heap.Fix] with the index of the element to update.
// Or you can just call [Heap.Update] with the index of the element to update and the new value to
// update to.
// At the same time, you can also remove any element by call [Heap.Remove] with the index of the
// element to remove.
func (h *Heap[T]) Elements() []T {
	return h.values
}

// Remove removes and returns the element at the given index i from heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[T]) Remove(i int) (value T, ok bool) {
	if i >= 0 && i < len(h.values) {
		n := len(h.values) - 1
		if i != n {
			h.swap(i, n)
		}
		value = h.values[n]
		ok = true
		h.values = h.values[:n]
		if i != n && !h.shiftDown(i) {
			h.shiftUp(i)
		}
		h.checkAndShrink()
	}
	return
}

// Fix re-establishes queue ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to, but less
// expensive than, calling [Heap.Remove] followed by a [Heap.Push] of the new value.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[T]) Fix(i int) {
	if i >= 0 && i < len(h.values) {
		if !h.shiftDown(i) {
			h.shiftUp(i)
		}
	}
}

// Update updates the element value to v at index i, and re-establishes heap ordering.
// [Heap.Update] is a shortcut for update an element value at the given index i with the given
// value v followed by a [Heap.Fix] with the given index i. It is equivalent to, but less
// expensive than, calling [Heap.Remove] with the given index i followed by a [Heap.Push] of the
// given value v.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[T]) Update(i int, v T) {
	if i >= 0 && i < len(h.values) {
		h.values[i] = v
		if !h.shiftDown(i) {
			h.shiftUp(i)
		}
	}
}
