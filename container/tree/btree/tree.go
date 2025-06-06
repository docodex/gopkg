// Package btree implements a B tree.
//
// According to Knuth's definition, a B-tree of order m is a tree which satisfies the
// following properties:
// - Every node has at most m children.
// - Every node, except for the root and the leaves, has at least ⌈m/2⌉ children.
// - The root node has at least two children unless it is a leaf.
// - A non-leaf node with k children contains k−1 keys.
// - All leaves appear on the same level.
//
// Reference: https://en.wikipedia.org/wiki/B-tree
package btree

import (
	"cmp"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docodex/gopkg/container"
)

// Entry represents a key-value pair of a node.
type Entry[K comparable, V any] struct {
	// The key used to compare entries.
	key K

	// The value stored with this entry.
	Value V
}

// Key returns the key of entry.
func (e *Entry[K, V]) Key() K {
	return e.key
}

// Node is a node of a B-tree.
type Node[K comparable, V any] struct {
	// The entries stored with this node.
	Entries []*Entry[K, V]

	// Children nodes of this node in tree.
	children []*Node[K, V]

	// Parent node of this node in tree.
	parent *Node[K, V]
}

// Children returns the children nodes, or nil if node has no child.
func (n *Node[K, V]) Children() []*Node[K, V] {
	return n.children
}

// Parent returns the parent node, or nil if node has no parent.
func (n *Node[K, V]) Parent() *Node[K, V] {
	return n.parent
}

// Height returns the height of subtree with node n as the root.
func (n *Node[K, V]) Height() int {
	h := 1
	for x := n; len(x.children) != 0; {
		h++
		x = x.children[0]
	}
	return h
}

// Len returns the number of entries of subtree with node n as the root.
// The complexity is O(n).
func (n *Node[K, V]) Len() int {
	count := len(n.Entries)
	for i := 0; i < len(n.children); i++ {
		count += n.children[i].Len()
	}
	return count
}

// MinNode returns the left-most (min) node which entries contains the minimum key of subtree with
// node n as the root.
func (n *Node[K, V]) MinNode() *Node[K, V] {
	x := n
	for len(x.children) != 0 {
		x = x.children[0]
	}
	return x
}

// MaxNode returns the right-most (max) node which entries contains the maximum key of subtree
// with node n as the root.
func (n *Node[K, V]) MaxNode() *Node[K, V] {
	x := n
	for len(x.children) != 0 {
		x = x.children[len(x.children)-1]
	}
	return x
}

// Min returns the entry which key is the minimum key of subtree with node n as the root.
func (n *Node[K, V]) Min() *Entry[K, V] {
	x := n.MinNode() // x must not be nil
	return x.Entries[0]
}

// Max returns the entry which key is the maximum key of subtree with node n as the root.
func (n *Node[K, V]) Max() *Entry[K, V] {
	x := n.MaxNode() // x must not be nil
	return x.Entries[len(x.Entries)-1]
}

// Tree represents a B-tree.
type Tree[K comparable, V any] struct {
	root *Node[K, V]          // the root node of tree
	m    int                  // order (maximum number of children for nodes)
	mid  int                  // (m-1)/2 or m/2, middle index of entries used for splitting
	len  int                  // current tree length which is the number of values stored in tree
	cmp  container.Compare[K] // function to compare tree nodes

	// minSize: m-1, maximum number of entries for nodes
	// maxSize: ceil(m/2)-1, minimum number of entries for nodes (except for the root and leaves)
	minSize, maxSize int
}

// New returns an initialized tree with [cmp.Compare] as the cmp function.
func New[K cmp.Ordered, V any](order int) *Tree[K, V] {
	return NewFunc[K, V](order, func(a, b K) int {
		return cmp.Compare(a, b)
	})
}

// NewFunc returns an initialized tree with the given function cmp as the cmp function.
func NewFunc[K comparable, V any](order int, cmp container.Compare[K]) *Tree[K, V] {
	if cmp == nil {
		cmp = func(a, b K) int {
			// just to cover nil cmp error
			return 0
		}
	}
	m := max(order, 3) // order m must be greater than 2
	return &Tree[K, V]{
		root:    nil,
		m:       m,
		mid:     (m - 1) / 2,
		len:     0,
		cmp:     cmp,
		minSize: (m+1)/2 - 1,
		maxSize: m - 1,
	}
}

// search returns the node which entries contains the given key k and the corresponding index in
// subtree with node x as the root, or nil and -1 if no such node found.
func (t *Tree[K, V]) search(x *Node[K, V], k K) (node *Node[K, V], index int) {
	index = -1
	for x != nil {
		i, ok := t.searchEntries(x, k)
		if ok {
			node = x
			index = i
			return
		}
		if len(x.children) == 0 {
			return
		}
		x = x.children[i]
	}
	return
}

// searchEntries searches the given key k only within the given node x among its entries.
func (t *Tree[K, V]) searchEntries(x *Node[K, V], k K) (index int, ok bool) {
	i, j := 0, len(x.Entries)-1
	for i <= j {
		mid := (j + i) / 2
		val := t.cmp(k, x.Entries[mid].key)
		switch {
		case val < 0:
			j = mid - 1
		case val > 0:
			i = mid + 1
		case val == 0:
			return mid, true
		}
	}
	return i, false
}

// Root returns the root node of tree, or nil if tree is empty.
func (t *Tree[K, V]) Root() *Node[K, V] {
	return t.root
}

// Height returns the height of tree.
func (t *Tree[K, V]) Height() int {
	if t.root == nil {
		return 0
	}
	return t.root.Height()
}

// Len returns the number of entries of tree t.
// The complexity is O(1).
func (t *Tree[K, V]) Len() int {
	return t.len
}

// Values returns all values in tree (in in-order traversal order).
func (t *Tree[K, V]) Values() []V {
	entries := t.InOrder()
	values := make([]V, 0, len(entries))
	for i := range entries {
		values = append(values, entries[i].Value)
	}
	return values
}

// Values returns all keys in tree (in in-order traversal order).
func (t *Tree[K, V]) Keys() []K {
	entries := t.InOrder()
	keys := make([]K, 0, len(entries))
	for i := range entries {
		keys = append(keys, entries[i].key)
	}
	return keys
}

// String returns the string representation of tree.
// Ref: std fmt.Stringer.
func (t *Tree[K, V]) String() string {
	var buf strings.Builder
	buf.WriteString("BTree\n")
	t.write(&buf, t.root, 0)
	return buf.String()
}

// write writes the structure of subtree with node x as the root to buffer buf.
func (t *Tree[K, V]) write(buf *strings.Builder, x *Node[K, V], level int) {
	if x == nil {
		return
	}
	for i := 0; i <= len(x.Entries); i++ {
		if i < len(x.children) {
			t.write(buf, x.children[i], level+1)
		}
		if i < len(x.Entries) {
			buf.WriteString(strings.Repeat("    ", level))
			fmt.Fprintf(buf, "%v:%v\n", x.Entries[i].key, x.Entries[i].Value)
		}
	}
}

// MarshalJSON marshals tree into valid JSON.
// Ref: std json.Marshaler.
func (t *Tree[K, V]) MarshalJSON() ([]byte, error) {
	m := make(map[K]V, t.len)
	t.Range(func(k K, v V) bool {
		m[k] = v
		return true
	})
	return json.Marshal(m)
}

// UnmarshalJSON unmarshals a JSON description of tree.
// The input can be assumed to be a valid encoding of a JSON value.
// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
// Ref: std json.Unmarshaler.
func (t *Tree[K, V]) UnmarshalJSON(data []byte) error {
	var m map[K]V
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	t.Clear()
	for k, v := range m {
		t.Insert(k, v)
	}
	return nil
}

// Insert inserts a new entry with the given key-value pair (k, v) to tree, or update the key and
// value if key k already exists in tree.
func (t *Tree[K, V]) Insert(k K, v V) {
	e := &Entry[K, V]{key: k, Value: v}
	if t.root == nil {
		t.root = &Node[K, V]{
			Entries:  []*Entry[K, V]{e},
			parent:   nil,
			children: nil,
		}
		t.len++
		return
	}
	if t.insert(t.root, e) {
		t.len++
	}
}

// insert inserts the given entry e to subtree with node x as the root.
// The given node x and entry e must not be nil.
func (t *Tree[K, V]) insert(x *Node[K, V], e *Entry[K, V]) (done bool) {
	i, ok := t.searchEntries(x, e.key)
	if ok {
		x.Entries[i] = e
		return false
	}
	// if x is a leaf node
	if len(x.children) == 0 {
		x.Entries = append(x.Entries, nil)
		copy(x.Entries[i+1:len(x.Entries)], x.Entries[i:len(x.Entries)-1])
		x.Entries[i] = e
		t.checkAndSplit(x)
		return true
	}
	// x is an internal node (not leaf)
	return t.insert(x.children[i], e)
}

// checkAndSplit checks the length of entries of the given node x, and split node if necessary.
func (t *Tree[K, V]) checkAndSplit(x *Node[K, V]) {
	if len(x.Entries) <= t.maxSize {
		return
	}
	if x == t.root {
		t.splitRoot()
		return
	}
	t.split(x)
}

// splitRoot splits entries and children nodes of root into left and right child nodes.
func (t *Tree[K, V]) splitRoot() {
	// split entries of root into left and right nodes
	left := &Node[K, V]{
		Entries:  append([]*Entry[K, V]{}, t.root.Entries[:t.mid]...),
		children: nil,
		parent:   nil,
	}
	right := &Node[K, V]{
		Entries:  append([]*Entry[K, V]{}, t.root.Entries[t.mid+1:]...),
		children: nil,
		parent:   nil,
	}
	// split children nodes of root into left and right nodes
	if len(t.root.children) != 0 {
		left.children = append([]*Node[K, V]{}, t.root.children[:t.mid+1]...)
		for i := range left.children {
			left.children[i].parent = left
		}
		right.children = append([]*Node[K, V]{}, t.root.children[t.mid+1:]...)
		for i := range right.children {
			right.children[i].parent = right
		}
	}
	// new root is a node with one entry and two children (left and right)
	root := &Node[K, V]{
		Entries:  []*Entry[K, V]{t.root.Entries[t.mid]},
		children: []*Node[K, V]{left, right},
		parent:   nil,
	}
	left.parent = root
	right.parent = root
	t.root = root
}

// split splits entries and children nodes of the given node x.
// The given node x must not be nil, also must not be the root of tree.
func (t *Tree[K, V]) split(x *Node[K, V]) {
	// parent p must not be nil as node x is not the root of tree
	p := x.parent
	// split entries of node x into left and right nodes
	left := &Node[K, V]{
		Entries:  append([]*Entry[K, V]{}, x.Entries[:t.mid]...),
		children: nil,
		parent:   p,
	}
	right := &Node[K, V]{
		Entries:  append([]*Entry[K, V]{}, x.Entries[t.mid+1:]...),
		children: nil,
		parent:   p,
	}
	// split entries and nodes of node x into left and right nodes
	if len(x.children) != 0 {
		left.children = append([]*Node[K, V]{}, x.children[:t.mid+1]...)
		for i := range left.children {
			left.children[i].parent = left
		}
		right.children = append([]*Node[K, V]{}, x.children[t.mid+1:]...)
		for i := range right.children {
			right.children[i].parent = right
		}
	}
	// insert middle entry in node x into parent
	i, _ := t.searchEntries(p, x.Entries[t.mid].key)
	p.Entries = append(p.Entries, nil)
	copy(p.Entries[i+1:len(p.Entries)], p.Entries[i:len(p.Entries)-1])
	p.Entries[i] = x.Entries[t.mid]
	// remove node x from children of parent p
	// insert left and right into children of parent p
	p.children = append(p.children, nil)
	copy(p.children[i+2:len(p.children)], p.children[i+1:len(p.children)-1])
	p.children[i] = left
	p.children[i+1] = right
	// check and split on parent p
	t.checkAndSplit(p)
}

// Remove removes the entry (and node) which key equals to the given key k from tree.
func (t *Tree[K, V]) Remove(k K) {
	if x, i := t.Search(k); x != nil {
		t.remove(x, i)
		t.len--
	}
}

// remove removes an entry at index i from entries of the given node x, checks and does fixup for
// the inbalance introduced by remove if necessary.
// The given node x must not be nil.
func (t *Tree[K, V]) remove(x *Node[K, V], i int) {
	// if x is an internal node (not leaf)
	if len(x.children) != 0 {
		// largest node in the left subtree must not be nil to satisfy the properties of B-tree
		y := x.children[i].MaxNode() // y must be a leaf node
		j := len(y.Entries) - 1      // lagest entry index in node y
		// replace the entry to be removed in node x with the lagest entry in node y
		x.Entries[i] = y.Entries[j]
		// transfer remove(x, i) to remove(y, j)
		x = y
		i = j
	}
	// now, x must be a leaf node
	k := x.Entries[i].key
	t.removeEntry(x, i)
	t.removeFixup(x, k)
	if len(t.root.Entries) == 0 {
		t.root = nil
	}
}

// removeFixup checks and does fixup for the inbalance introduced by remove if necessary.
// The given node x is the node just removed entry, and the given key k is the key of entry just
// removed from entries of node x.
// At the same time, the given node x is a leaf node or its number of children already matches its
// number of entries after remove.
func (t *Tree[K, V]) removeFixup(x *Node[K, V], k K) {
	if x == nil || len(x.Entries) >= t.minSize {
		return
	}

	// try to borrow from left sibling
	x1, i1 := t.leftSibling(x, k)
	if x1 != nil && len(x1.Entries) > t.minSize {
		// rotate right
		x.Entries = append([]*Entry[K, V]{x.parent.Entries[i1]}, x.Entries...)
		j := len(x1.Entries) - 1
		x.parent.Entries[i1] = x1.Entries[j]
		t.removeEntry(x1, j)
		if len(x1.children) != 0 {
			j := len(x1.children) - 1
			x1.children[j].parent = x
			x.children = append([]*Node[K, V]{x1.children[j]}, x.children...)
			t.removeChild(x1, j)
		}
		return
	}

	// try to borrow from right sibling
	x2, i2 := t.rightSibling(x, k)
	i3 := i2 - 1
	if x2 != nil && len(x2.Entries) > t.minSize {
		// rotate left
		x.Entries = append(x.Entries, x.parent.Entries[i3])
		x.parent.Entries[i3] = x2.Entries[0]
		t.removeEntry(x2, 0)
		if len(x2.children) != 0 {
			x2.children[0].parent = x
			x.children = append(x.children, x2.children[0])
			t.removeChild(x2, 0)
		}
		return
	}

	// merge with siblings
	if x1 != nil {
		// merge with left sibling
		entries := append(x1.Entries, x.parent.Entries[i1])
		x.Entries = append(entries, x.Entries...)
		k = x.parent.Entries[i1].key
		t.removeEntry(x.parent, i1)
		for i := range x1.children {
			x1.children[i].parent = x
		}
		x.children = append(x1.children, x.children...)
		t.removeChild(x.parent, i1)
	} else if x2 != nil {
		// merge with right sibling
		x.Entries = append(x.Entries, x.parent.Entries[i3])
		x.Entries = append(x.Entries, x2.Entries...)
		k = x.parent.Entries[i3].key
		t.removeEntry(x.parent, i3)
		for i := range x2.children {
			x2.children[i].parent = x
		}
		x.children = append(x.children, x2.children...)
		t.removeChild(x.parent, i2)
	}
	// update the root of tree if root becomes empty by merge
	if x.parent == t.root && len(t.root.Entries) == 0 {
		x.parent = nil
		t.root = x
		return
	}

	// parent might underflow
	// check and do fixup for the inbalance introduced by remove on parent if necessary
	t.removeFixup(x.parent, k)
}

// removeEntry just removes an entry at index i from entries of the given node x.
// The given node x must not be nil.
func (t *Tree[K, V]) removeEntry(x *Node[K, V], i int) {
	if i >= 0 && i < len(x.Entries) {
		last := len(x.Entries) - 1
		copy(x.Entries[i:last], x.Entries[i+1:len(x.Entries)])
		x.Entries[last] = nil
		x.Entries = x.Entries[:last]
	}
}

// removeChild just removes a child at index i from children of the given node x.
// The given node x must not be nil.
func (t *Tree[K, V]) removeChild(x *Node[K, V], i int) {
	if i >= 0 && i < len(x.children) {
		last := len(x.children) - 1
		copy(x.children[i:last], x.children[i+1:len(x.children)])
		x.children[last] = nil
		x.children = x.children[:last]
	}
}

// leftSibling returns the left sibling of node x and the corresponding child index (in parent) if
// exist, or (nil, -1) would be returned.
// The given node x must not be nil.
func (t *Tree[K, V]) leftSibling(x *Node[K, V], k K) (node *Node[K, V], index int) {
	index = -1
	if x.parent != nil {
		i, _ := t.searchEntries(x.parent, k)
		i--
		if i >= 0 && i < len(x.parent.children) {
			node = x.parent.children[i]
			index = i
		}
	}
	return
}

// leftSibling returns the right sibling of node x and the corresponding child index (in parent)
// if exist, or (nil, -1) would be returned.
// The given node x must not be nil.
func (t *Tree[K, V]) rightSibling(x *Node[K, V], k K) (node *Node[K, V], index int) {
	index = -1
	if x.parent != nil {
		i, _ := t.searchEntries(x.parent, k)
		i++
		if i > 0 && i < len(x.parent.children) {
			node = x.parent.children[i]
			index = i
		}
	}
	return
}

// Search returns the node which entries contains the given key k and the corresponding index in
// tree, or nil and -1 if no such node found.
func (t *Tree[K, V]) Search(k K) (node *Node[K, V], index int) {
	return t.search(t.root, k)
}

// Get returns the value which key equals to the given key k.
// The ok result indicates whether such value was found in tree.
func (t *Tree[K, V]) Get(k K) (value V, ok bool) {
	if x, i := t.Search(k); x != nil {
		value = x.Entries[i].Value
		ok = true
	}
	return
}

// Min returns the entry which key is the minimum key of tree.
func (t *Tree[K, V]) Min() (entry *Entry[K, V]) {
	if t.root != nil {
		entry = t.root.Min()
	}
	return
}

// Max returns the entry which key is the maximum key of tree.
func (t *Tree[K, V]) Max() (entry *Entry[K, V]) {
	if t.root != nil {
		entry = t.root.Max()
	}
	return
}

// LevelOrder performs level-order traversal for tree, and returns a slice of entries as the
// result.
func (t *Tree[K, V]) LevelOrder() []*Entry[K, V] {
	entries := make([]*Entry[K, V], 0, t.len)
	var q []*Node[K, V] // queue
	if t.root != nil {
		q = append(q, t.root)
	}
	for len(q) != 0 {
		x := q[0]
		q = q[1:]
		entries = append(entries, x.Entries...)
		if len(x.children) != 0 {
			q = append(q, x.children...)
		}
	}
	return entries
}

// InOrder performs in-order traversal for tree, and returns a slice of entries as the result.
func (t *Tree[K, V]) InOrder() []*Entry[K, V] {
	entries := make([]*Entry[K, V], 0, t.len)
	t.inOrder(t.root, &entries)
	return entries
}

// inOrder performs in-order traversal for subtree with node x as the root, and append the result
// to a slice entries.
func (t *Tree[K, V]) inOrder(x *Node[K, V], entries *[]*Entry[K, V]) {
	if x != nil {
		// if x is a leaf node, append its entries
		if len(x.children) == 0 {
			*entries = append(*entries, x.Entries...)
			return
		}
		// x is not a leaf node, traverse each child before the corresponding entry
		for i := 0; i < len(x.Entries); i++ {
			t.inOrder(x.children[i], entries)
			*entries = append(*entries, x.Entries[i])
		}
		// traverse the last child
		t.inOrder(x.children[len(x.Entries)], entries)
	}
}

// Clear removes all nodes in tree.
func (t *Tree[K, V]) Clear() {
	var q []*Node[K, V] // queue
	if t.root != nil {
		q = append(q, t.root)
	}
	for len(q) != 0 {
		x := q[0]
		q = q[1:]
		if len(x.children) != 0 {
			q = append(q, x.children...)
		}
		clear(x.children)
		x.children = nil
		x.parent = nil
	}
	t.root = nil
	t.len = 0
}

// item represents an item in stack.
type item[K comparable, V any] struct {
	node  *Node[K, V] // node to process
	index int         // index of the entry/child to process next within the node
}

// Range calls f sequentially for each entry present in tree in in-order traversal order.
// If f returns false, range stops the iteration.
func (t *Tree[K, V]) Range(f func(k K, v V) bool) {
	if f == nil {
		return
	}
	var s []*item[K, V] // stack
	x := t.root
	for x != nil || len(s) != 0 {
		for x != nil {
			s = append(s, &item[K, V]{node: x, index: 0})
			if len(x.children) == 0 {
				x = nil
			} else {
				x = x.children[0]
			}
		}
		// now, x == nil, and len(s) != 0
		e := s[len(s)-1]
		s = s[:len(s)-1]
		x = e.node
		if len(x.children) == 0 {
			stop := false
			for i := range x.Entries {
				if !f(x.Entries[i].key, x.Entries[i].Value) {
					stop = true
					break
				}
			}
			if stop {
				// x = nil
				break
			}
			x = nil
			continue
		}
		i := e.index
		if i < len(x.Entries) && !f(x.Entries[i].key, x.Entries[i].Value) {
			// x = nil
			break
		}
		i++
		// if didn't finish all the entries of node x, update the index and push it back to stack
		if i < len(x.Entries) {
			e.index = i
			s = append(s, e)
		}
		// if didn't finish all the children of node x, push the next child to stack
		if i < len(x.children) {
			x = x.children[i]
		} else {
			x = nil
		}
	}
}
