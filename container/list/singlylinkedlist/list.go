// Package singlylinkedlist implements a singly linked list.
//
// To iterate over a list (where l is a *List):
//
//	for x := l.Front(); x != nil; x = x.Next() {
//		// do something with x.Value
//	}
//
// or:
//
//	l.Range(func(index int, value T) bool {
//		// do something with index and value
//		return true
//	})
package singlylinkedlist

import (
	"encoding/json"
	"slices"

	"github.com/docodex/gopkg/container"
	"github.com/docodex/gopkg/jsonx"
)

// Node is a node of a linked list.
type Node[T any] struct {
	// The value stored with this node.
	Value T

	// Next pointer in the singly-linked list of nodes.
	// To simplify the implementation, internally a list l is implemented as a ring, such that
	// &l.root is both the next node of the last list node (l.BackNode()) and the previous node of
	// the first list node (l.FrontNode()).
	next *Node[T]

	// The list to which this node belongs.
	list *List[T]
}

// Next returns the next list node or nil.
func (n *Node[T]) Next() *Node[T] {
	if x := n.next; n.list != nil && x != &n.list.root {
		return x
	}
	return nil
}

// List represents a singly linked list.
type List[T any] struct {
	root Node[T]  // sentinel list node, only &root and root.next are used
	last *Node[T] // the last node in list, or point to root while list is empty
	len  int      // current list length excluding the sentinel node
}

// New returns an initialized list with the values v added.
func New[T any](v ...T) *List[T] {
	l := new(List[T]).init()
	l.insert(&l.root, v...)
	return l
}

// init initializes or clears list l.
func (l *List[T]) init() *List[T] {
	l.root.next = &l.root
	l.last = &l.root
	l.len = 0
	return l
}

// insert inserts new nodes with the given values v after at, increments l.len, and returns the
// first node just inserted.
func (l *List[T]) insert(at *Node[T], v ...T) *Node[T] {
	if len(v) == 0 {
		return nil
	}
	x := at
	y := at.next
	for i := range v {
		x.next = &Node[T]{
			Value: v[i],
			list:  l,
		}
		x = x.next
	}
	x.next = y
	if l.last == at {
		l.last = x
	}
	l.len += len(v)
	return at.next
}

// remove removes x (with previous node prev) from its list, decrements l.len, and returns the
// removed node value.
func (l *List[T]) remove(x, prev *Node[T]) (value T, ok bool) {
	if x == &l.root {
		return
	}
	prev.next = x.next
	if l.last == x {
		l.last = prev
	}
	x.next = nil
	x.list = nil
	l.len--
	value = x.Value
	ok = true
	return
}

// Len returns the number of nodes of list l (excluding sentinel nodes).
// The complexity is O(1).
func (l *List[T]) Len() int {
	return l.len
}

// Values returns a slice of all values of list.
func (l *List[T]) Values() []T {
	values := make([]T, l.len)
	for i, x := 0, l.root.next; i < l.len; i, x = i+1, x.next {
		values[i] = x.Value
	}
	return values
}

// String returns the string representation of list.
// Ref: std fmt.Stringer.
func (l *List[T]) String() string {
	values, _ := jsonx.MarshalToString(l.Values())
	return "SinglyLinkedList: " + values
}

// MarshalJSON marshals list into valid JSON.
// Ref: std json.Marshaler.
func (l *List[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Values())
}

// UnmarshalJSON unmarshals a JSON description of list.
// The input can be assumed to be a valid encoding of a JSON value.
// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
// Ref: std json.Unmarshaler.
func (l *List[T]) UnmarshalJSON(data []byte) error {
	var v []T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	l.Clear()
	l.PushBack(v...)
	return nil
}

// FrontNode returns the first node of list l or nil if list is empty.
func (l *List[T]) FrontNode() *Node[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// BackNode returns the last node of list l or nil if list is empty.
func (l *List[T]) BackNode() *Node[T] {
	if l.len == 0 {
		return nil
	}
	return l.last
}

// Front returns the first node value if exists in list.
// The ok result indicates whether such node was found in list.
func (l *List[T]) Front() (value T, ok bool) {
	if x := l.FrontNode(); x != nil {
		value = x.Value
		ok = true
	}
	return
}

// Back returns the last node value if exists in list.
// The ok result indicates whether such node was found in list.
func (l *List[T]) Back() (value T, ok bool) {
	if x := l.BackNode(); x != nil {
		value = x.Value
		ok = true
	}
	return
}

// PushFront inserts new nodes with the given values v at the front of list.
func (l *List[T]) PushFront(v ...T) {
	l.insert(&l.root, v...)
}

// PushBack inserts new nodes with the given values v at the back of list.
func (l *List[T]) PushBack(v ...T) {
	l.insert(l.last, v...)
}

// PopFront removes the first node if exists in list and returns its value.
// The ok result indicates whether such node was removed from list.
func (l *List[T]) PopFront() (value T, ok bool) {
	if x := l.FrontNode(); x != nil {
		value, ok = l.remove(x, &l.root)
	}
	return
}

// PopBack removes the last node if exists in list and returns its value.
// The ok result indicates whether such node was removed from list.
func (l *List[T]) PopBack() (value T, ok bool) {
	if x := l.BackNode(); x != nil {
		y := &l.root
		for y.next != x {
			y = y.next
		}
		value, ok = l.remove(x, y)
	}
	return
}

// Clear removes all nodes in list.
func (l *List[T]) Clear() {
	for x := l.root.next; x != &l.root; {
		y := x.next
		x.next = nil // avoid memory leaks
		x.list = nil
		x = y
	}
	l.init()
}

// indexGet gets the node of index i if exists in list, or nil if index i is invalid.
func (l *List[T]) indexGet(i int) (*Node[T], bool) {
	if i < 0 || i >= l.len {
		return nil, false
	}
	if i == l.len-1 {
		return l.last, true
	}
	j, x := 0, l.root.next
	for j < i {
		j, x = j+1, x.next
	}
	return x, true
}

// Get returns the node value of index i if exists in list.
// The ok result indicates whether such node was found in list.
func (l *List[T]) Get(i int) (value T, ok bool) {
	var x *Node[T]
	x, ok = l.indexGet(i)
	if ok {
		value = x.Value
	}
	return
}

// Set sets the value to v of index i if exists in list.
func (l *List[T]) Set(i int, v T) {
	if x, ok := l.indexGet(i); ok {
		x.Value = v
	}
}

// Add inserts new nodes with the given values v to index i if exists in list, or appends new
// nodes with the given value v to the back of list if index i points to the next index of the
// last element in list.
func (l *List[T]) Add(i int, v ...T) {
	if i == l.len {
		l.PushBack(v...)
		return
	}
	if i == 0 {
		l.PushFront(v...)
		return
	}
	// if i-1 not exists, then i<=0 or l.len<i, skip this insert
	if x, ok := l.indexGet(i - 1); ok {
		l.insert(x, v...)
	}
}

// Del removes the node at index i if exists in list.
func (l *List[T]) Del(i int) {
	if l.len == 0 {
		return
	}
	if i == 0 {
		l.remove(l.root.next, &l.root)
		return
	}
	// if i-1 not exists, then i<=0 or l.len<i, skip this insert
	if x, ok := l.indexGet(i - 1); ok {
		l.remove(x.next, x)
	}
}

// Swap swaps the values with indices i and j if both indices exist in list.
func (l *List[T]) Swap(i, j int) {
	if i == j || i < 0 || i >= l.len || j < 0 || j >= l.len {
		return
	}
	var x, y *Node[T]
	for k, z := 0, l.root.next; x == nil || y == nil; k, z = k+1, z.next {
		switch k {
		case i:
			x = z
		case j:
			y = z
		default:
		}
	}
	x.Value, y.Value = y.Value, x.Value
}

// Sort sorts list values (in-place) with the given cmp.
func (l *List[T]) Sort(cmp container.Compare[T]) {
	if cmp != nil && l.len > 1 {
		values := l.Values()
		slices.SortFunc(values, cmp)
		l.Clear()
		l.insert(&l.root, values...)
	}
}

// Range calls f sequentially for each index i and value v present in list.
// If f returns false, range stops the iteration.
func (l *List[T]) Range(f func(i int, v T) bool) {
	if f == nil {
		return
	}
	for i, x := 0, l.root.next; i < l.len; i, x = i+1, x.next {
		if !f(i, x.Value) {
			break
		}
	}
}

// InsertAfter inserts new nodes with the given values v immediately after mark, and returns the
// first node just inserted.
// If mark is not a node of l, the list is not modified.
func (l *List[T]) InsertAfter(mark *Node[T], v ...T) *Node[T] {
	if mark == nil || mark.list != l {
		return nil
	}
	// if mark.list == l, l must have been initialized when mark was inserted in l
	return l.insert(mark, v...)
}

// RemoveAfter removes node after mark from l if mark is a node of list l.
// It returns the node value just removed.
func (l *List[T]) RemoveAfter(mark *Node[T]) (value T, ok bool) {
	if mark != nil && mark.list == l {
		// if mark.list == l, l must have been initialized when mark was inserted in l
		if x := mark.next; x != nil && x != &l.root {
			value = x.Value
			ok = true
			l.remove(x, mark)
		}
	}
	return
}

// PushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same.
func (l *List[T]) PushFrontList(other *List[T]) {
	if other != nil {
		mark := &l.root
		for i, x := other.Len(), other.FrontNode(); i > 0; i, x = i-1, x.Next() {
			mark = l.insert(mark, x.Value)
		}
	}
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same.
func (l *List[T]) PushBackList(other *List[T]) {
	if other != nil {
		for i, x := other.Len(), other.FrontNode(); i > 0; i, x = i-1, x.Next() {
			l.insert(l.last, x.Value)
		}
	}
}
