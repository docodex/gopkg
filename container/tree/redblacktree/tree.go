// Package redblacktree implements a red-black tree.
//
// Could be used as backed data-structure by TreeSet and TreeMap.
//
// References:
// - http://en.wikipedia.org/wiki/Red%E2%80%93black_tree
// - https://en.wikipedia.org/wiki/AVL_tree
// - https://en.wikipedia.org/wiki/Binary_search_tree
package redblacktree

import (
	"cmp"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docodex/gopkg/container"
)

type color int8

const (
	red color = iota
	black
)

// Node is a node of a binary tree.
type Node[K comparable, V any] struct {
	// The key used to compare nodes.
	key K

	// The value stored with this node.
	Value V

	// The color of this node, red or black. Root and nil (leaf) node are always black.
	color color

	// Left and right children nodes of this node in tree.
	left, right *Node[K, V]

	// Parent node of this node in tree.
	parent *Node[K, V]
}

// newNode returns a new node with the given key k as the key, the given value v as the value, the
// given color c as the color, and the given node p as the parent.
func newNode[K comparable, V any](k K, v V, c color, p *Node[K, V]) *Node[K, V] {
	return &Node[K, V]{
		key:    k,
		Value:  v,
		color:  c,
		left:   nil,
		right:  nil,
		parent: p,
	}
}

// Key returns the key of node.
func (n *Node[K, V]) Key() K {
	return n.key
}

// Left returns the left child node, or nil if node has no left child.
func (n *Node[K, V]) Left() *Node[K, V] {
	return n.left
}

// Right returns the right child node, or nil if node has no right child.
func (n *Node[K, V]) Right() *Node[K, V] {
	return n.right
}

// Parent returns the parent node, or nil if node has no parent.
func (n *Node[K, V]) Parent() *Node[K, V] {
	return n.parent
}

// Len returns the number of nodes of subtree with node n as the root.
// The complexity is O(n).
func (n *Node[K, V]) Len() int {
	count := 1
	if n.left != nil {
		count += n.left.Len()
	}
	if n.right != nil {
		count += n.right.Len()
	}
	return count
}

// Min returns the node which key is the minimum key of subtree with node n as the root.
func (n *Node[K, V]) Min() *Node[K, V] {
	x := n
	for x.left != nil {
		x = x.left
	}
	return x
}

// Max returns the node which key is the maximum key of subtree with node n as the root.
func (n *Node[K, V]) Max() *Node[K, V] {
	x := n
	for x.right != nil {
		x = x.right
	}
	return x
}

// Tree represents an red-black tree.
type Tree[K comparable, V any] struct {
	root *Node[K, V]          // the root node of tree
	len  int                  // current tree length which is the number of nodes of tree
	cmp  container.Compare[K] // function to compare tree nodes
}

// New returns an initialized tree with [cmp.Compare] as the cmp function.
func New[K cmp.Ordered, V any]() *Tree[K, V] {
	return NewFunc[K, V](func(a, b K) int {
		return cmp.Compare(a, b)
	})
}

// NewFunc returns an initialized tree with the given function cmp as the cmp function.
func NewFunc[K comparable, V any](cmp container.Compare[K]) *Tree[K, V] {
	if cmp == nil {
		cmp = func(a, b K) int {
			// just to cover nil cmp error
			return 0
		}
	}
	return &Tree[K, V]{
		root: nil,
		len:  0,
		cmp:  cmp,
	}
}

// color returns the color of the given node x. Nil (leaf) node is black.
func (t *Tree[K, V]) color(x *Node[K, V]) color {
	if x == nil {
		// color of nil (leaf) node is black
		return black
	}
	return x.color
}

// rightRotate do right rotate operation, nodes x and x.left must not be nil.
func (t *Tree[K, V]) rightRotate(x *Node[K, V]) {
	y := x.left
	// after right rotate, the right child of node y would be moved as the left child of node x
	x.left = y.right
	if x.left != nil {
		x.left.parent = x
	}
	// rotate node x to the right around node y
	y.right = x
	y.parent = x.parent
	if y.parent == nil {
		t.root = y
	} else if y.parent.left == x {
		y.parent.left = y
	} else {
		y.parent.right = y
	}
	x.parent = y
}

// leftRotate do left rotate operation, nodes x and x.right must not be nil.
func (t *Tree[K, V]) leftRotate(x *Node[K, V]) {
	y := x.right
	// after left rotate, the left child of node y would be moved as the right child of node x
	x.right = y.left
	if x.right != nil {
		x.right.parent = x
	}
	// rotate node x to the left around node y
	y.left = x
	y.parent = x.parent
	if y.parent == nil {
		t.root = y
	} else if y.parent.left == x {
		y.parent.left = y
	} else {
		y.parent.right = y
	}
	x.parent = y
}

// search returns the node which key equals to the given key k from subtree with node x as the
// root, or nil if no such node found.
func (t *Tree[K, V]) search(x *Node[K, V], k K) *Node[K, V] {
	for x != nil {
		r := t.cmp(k, x.key)
		if r < 0 {
			x = x.left
		} else if r > 0 {
			x = x.right
		} else {
			return x
		}
	}
	return nil
}

// Root returns the root node of tree, or nil if tree is empty.
func (t *Tree[K, V]) Root() *Node[K, V] {
	return t.root
}

// Len returns the number of nodes of tree t.
// The complexity is O(1).
func (t *Tree[K, V]) Len() int {
	return t.len
}

// Values returns all values in tree (in in-order traversal order).
func (t *Tree[K, V]) Values() []V {
	_, values := t.InOrder()
	return values
}

// Values returns all keys in tree (in in-order traversal order).
func (t *Tree[K, V]) Keys() []K {
	keys, _ := t.InOrder()
	return keys
}

// String returns the string representation of tree.
// Ref: std fmt.Stringer.
func (t *Tree[K, V]) String() string {
	var buf strings.Builder
	buf.WriteString("RedBlackTree\n")
	t.write(&buf, t.root, "", true)
	return buf.String()
}

// write writes the structure of subtree with node x as the root to buffer buf.
func (t *Tree[K, V]) write(buf *strings.Builder, x *Node[K, V], prefix string, tail bool) {
	if x == nil {
		return
	}
	if x.right != nil {
		newPrefix := prefix
		if tail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		t.write(buf, x.right, newPrefix, false)
	}
	buf.WriteString(prefix)
	if tail {
		buf.WriteString("└── ")
	} else {
		buf.WriteString("┌── ")
	}
	color := "b"
	if x.color == red {
		color = "r"
	}
	fmt.Fprintf(buf, "%v:%v[%s]\n", x.key, x.Value, color)
	if x.left != nil {
		newPrefix := prefix
		if tail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		t.write(buf, x.left, newPrefix, true)
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

// Insert inserts a new node with the given key-value pair (k, v) to tree, or update the key and
// value if key k already exists in tree.
func (t *Tree[K, V]) Insert(k K, v V) {
	// find the right position and do insert
	x := t.root
	for x != nil {
		r := t.cmp(k, x.key)
		if r < 0 {
			if x.left == nil {
				// insert new red node as the left child of node x
				x.left = newNode(k, v, red, x) // color of non-root node initialized to red
				x = x.left
				t.len++
				break
			} else {
				x = x.left
			}
		} else if r > 0 {
			if x.right == nil {
				// insert new red node as the right child of node x
				x.right = newNode(k, v, red, x) // color of non-root node initialized to red
				x = x.right
				t.len++
				break
			} else {
				x = x.right
			}
		} else {
			// duplicated key found, update the key and value
			x.key = k
			x.Value = v
			return
		}
	}
	// node x is nil, then t.root must be nil, insert new black node as the root of tree
	if x == nil {
		t.root = newNode(k, v, black, nil)
		t.len++
		return
	}
	// check and do fixup for the inbalance introduced by insert if necessary
	t.insertFixup(x)
}

// fixupInsert checks and does fixup for the inbalance introduced by insert if necessary.
// The given node x (just inserted) must not be nil, and the color of node x is red.
func (t *Tree[K, V]) insertFixup(x *Node[K, V]) {
	// as node x is red, the parent of node x should not be red
	for x.parent != nil && x.parent.color == red {
		// x.parent.parent should not be nil, for x.parent is red and cannot be the root of tree
		// x.parent.parent should be black, for x.parent is red
		if x.parent == x.parent.parent.left {
			x = t.insertFixupLeft(x)
		} else {
			x = t.insertFixupRight(x)
		}
	}
	// root node must be black, update root color for root node may be changed by rotate
	t.root.color = black
}

// insertFixupLeft checks and does fixup for the inbalance introduced by insert if necessary.
// The given node x must not be nil, and the color of node x is red.
// At the same time, both parent and grandparent of node x are not nil, and parent is the left
// child of grandparent of node x.
func (t *Tree[K, V]) insertFixupLeft(x *Node[K, V]) *Node[K, V] {
	// if uncle is red
	if x.parent.parent.right != nil && x.parent.parent.right.color == red {
		// both parent and uncle are red, exchange the color of them and the grandparent
		// and do fixup with grandparent node
		x.parent.color = black
		x.parent.parent.right.color = black
		x.parent.parent.color = red
		return x.parent.parent
	}
	// as uncle is not red, then it must be nil to satisfy the balance of tree
	if x == x.parent.right {
		// parent is the left child of grandparent, and x is the right child of parent
		// so, do left rotation on parent first
		x = x.parent
		t.leftRotate(x)
	}
	// uncle is nil, exchange the color of parent and grandparent
	x.parent.color = black
	x.parent.parent.color = red
	// parent is the left child of grandparent, do right rotate on grandparent
	t.rightRotate(x.parent.parent)
	return x
}

// insertFixupRight checks and does fixup for the inbalance introduced by insert if necessary.
// The given node x must not be nil, and the color of node x is red.
// At the same time, both parent and grandparent of node x are not nil, and parent is the right
// child of grandparent of node x.
func (t *Tree[K, V]) insertFixupRight(x *Node[K, V]) *Node[K, V] {
	// uncle is red
	if x.parent.parent.left != nil && x.parent.parent.left.color == red {
		// both parent and uncle are red, exchange the color of them and the grandparent
		// and do fixup with grandparent node
		x.parent.color = black
		x.parent.parent.left.color = black
		x.parent.parent.color = red
		return x.parent.parent
	}
	// as uncle is not red, then it must be nil to satisfy the balance of tree
	if x == x.parent.left {
		// parent is the right child of grandparent, and x is the left child of parent
		// so, do right rotation on parent first
		x = x.parent
		t.rightRotate(x)
	}
	// uncle is nil or black, exchange the color of parent and grandparent
	x.parent.color = black
	x.parent.parent.color = red
	// parent is the right child of grandparent, do left rotate on grandparent
	t.leftRotate(x.parent.parent)
	return x
}

// Remove removes the node which key equals to the given key k from tree.
func (t *Tree[K, V]) Remove(k K) {
	if x := t.search(t.root, k); x != nil {
		t.remove(x)
		t.len--
	}
}

// remove removes the node x from tree, checks and does fixup for the inbalance introduced by
// remove if necessary.
// The given node x must not be nil.
// 1. x has tow children, replace the key and value of x with the key and value of its next node
// x in in-order traversal order, then call remove(x)
// 2. x has only one child: left or right, then the subtree (left or right) must have only one
// node x (x.left or x.right), and x must be red for the balance of tree, replace the key and
// value of node x with the key and value of node x, then remove x directly
// 3. x has no children: if x is red, remove it directly; or, check the color of sibling and
// parent of node x
func (t *Tree[K, V]) remove(x *Node[K, V]) {
	// if node x has two children, find the next node y in in-order traversal order
	// replace x.key with y.key, replace x.Value with y.Value, and transfer remove(x) to remove(y)
	if x.left != nil && x.right != nil {
		y := x.right
		for y.left != nil {
			y = y.left
		}
		x.key = y.key
		x.Value = y.Value
		x = y
	}
	// now, node x has at most one child: left or child
	// if x has left child, then it has no right child, and x.left must be red for balance
	// replace the key and value of x with the key and value of x.left, remove x.left and return
	if x.left != nil {
		x.key = x.left.key
		x.Value = x.left.Value
		x.left.parent = nil
		x.left = nil
		return
	}
	// if x has right child, then it has no left child, and x.right must be red for balance
	// replace the key and value of x with the key and value of x.right, remove x.right and return
	if x.right != nil {
		x.key = x.right.key
		x.Value = x.right.Value
		x.right.parent = nil
		x.right = nil
		return
	}
	// now, x has no children: neither left, nor right
	// if x.parent is nil, then x must be the root of tree, remove x and return
	if x.parent == nil {
		t.root = nil
		return
	}
	// now, x.parent is not nil, then x must not be the root of tree, remove x
	if x == x.parent.left {
		x.parent.left = nil
	} else {
		x.parent.right = nil
	}
	p := x.parent
	x.parent = nil
	// now, x.parent has only one child: left or right
	// if x (just removed) was black, removing it could break up the balance of tree
	// check and do fixup for the inbalance introduced by remove if necessary
	if x.color == black {
		t.removeFixup(p)
	}
}

// removeFixup checks and does fixup for the inbalance introduced by remove if necessary.
// The given node p is the parent node of the node just removed, node p must not be nil, and has
// only one child.
func (t *Tree[K, V]) removeFixup(p *Node[K, V]) {
	// node x is defined as the node replacing the position of the node just removed
	var x *Node[K, V]
	for x != t.root && t.color(x) == black {
		// 1. x is not the root of tree
		// 2. x has dual-black color (one black is derived from the node just removed)
		// 3. p must not be nil:
		// 3.1 if x is nil (first round of loop), p is the parent of the node just removed
		// 3.2 if x is not nil, p is the parent of x
		// 4. sibling of x is not nil as x is dual-black (for balance)
		if x != nil {
			p = x.parent
		}
		if x == p.left {
			x = t.removeFixupLeft(p)
		} else {
			x = t.removeFixupRight(p)
		}
	}
	// now x is the root of tree (p == nil), or x has red-black color, set it to single-black
	if x != nil {
		x.color = black
	}
}

// removeFixupLeft checks and does fixup for the inbalance introduced by remove if necessary.
// The given node p must not be nil, and has a left color child (could be nil) with dual-black.
// So, the right child of node p is not nil (for balance).
func (t *Tree[K, V]) removeFixupLeft(p *Node[K, V]) *Node[K, V] {
	if p.right.color == red {
		// if sibling is red, then parent p is black, and sibling must have two black children
		// exchange the color of parent p and sibline, and do left lotate
		p.color = red
		p.right.color = black
		t.leftRotate(p)
	}
	// now, current sibling is black
	if t.color(p.right.right) == black && t.color(p.right.left) == black {
		// sibling is black, and has two black children, while x (current node) is daul-black
		// consider moving a black color from sibling and x to parent p
		// then, parent p has an additional black color
		// so, p has daul-black or red-black color
		// focus on parent p and continue the loop
		p.right.color = red
		return p
	}
	if t.color(p.right.right) == black && t.color(p.right.left) == red {
		// sibling is black, and further nephew is black (or nil), closer nephew is red
		// exchange the color of sibling and closer nephew, then do right rotate on sibling
		p.right.color = red
		p.right.left.color = black
		t.rightRotate(p.right)
	}
	// now, sibling is black, further nephew is red
	// 1. exchange the color of parent p and sibling
	// 2. set further nephew to black
	// 3. do left rotate
	// 4. break up the loop on parent p
	p.right.color = p.color
	p.color = black
	p.right.right.color = black
	t.leftRotate(p)
	return t.root
}

// removeFixupRight checks and does fixup for the inbalance introduced by remove if necessary.
// The given node p must not be nil, and has a right child (could be nil) with dual-black color.
// So, the left child of node p is not nil (for balance).
func (t *Tree[K, V]) removeFixupRight(p *Node[K, V]) *Node[K, V] {
	if p.left.color == red {
		// if sibling is red, then parent p is black, and sibling must have two black children
		// exchange the color of parent p and sibline, and do right lotate
		p.color = red
		p.left.color = black
		t.rightRotate(p)
	}
	// now, current sibling is black
	if t.color(p.left.left) == black && t.color(p.left.right) == black {
		// sibling is black, and has two black children, while x (current node) is daul-black
		// consider moving a black color from sibling and x to parent p
		// then, parent p has an additional black color
		// so, p has daul-black or red-black color
		// focus on parent p and continue the loop
		p.left.color = red
		return p
	}
	if t.color(p.left.left) == black && t.color(p.left.right) == red {
		// sibling is black, and further nephew is black (or nil), closer nephew is red
		// exchange the color of sibling and closer nephew, then do left rotate on sibling
		p.left.color = red
		p.left.right.color = black
		t.leftRotate(p.left)
	}
	// now, sibling is black, further nephew is red
	// 1. exchange the color of parent p and sibling
	// 2. set further nephew to black
	// 3. do left rotate
	// 4. break up the loop on parent p
	p.left.color = p.color
	p.color = black
	p.left.left.color = black
	t.rightRotate(p)
	return t.root
}

// Search returns the node which key equals to the given key k, or nil if no such node found.
func (t *Tree[K, V]) Search(k K) *Node[K, V] {
	return t.search(t.root, k)
}

// Get returns the value which key equals to the given key k.
// The ok result indicates whether such value was found in tree.
func (t *Tree[K, V]) Get(k K) (value V, ok bool) {
	if x := t.Search(k); x != nil {
		value = x.Value
		ok = true
	}
	return
}

// Min returns the node which key is the minimum key of tree.
func (t *Tree[K, V]) Min() *Node[K, V] {
	if t.root == nil {
		return nil
	}
	return t.root.Min()
}

// Max returns the node which key is the maximum key of tree.
func (t *Tree[K, V]) Max() *Node[K, V] {
	if t.root == nil {
		return nil
	}
	return t.root.Max()
}

// Prev returns the previous node (in in-order traversal order) of the given node x, or nil if no
// such node found.
func (t *Tree[K, V]) Prev(x *Node[K, V]) *Node[K, V] {
	if x == nil {
		return nil
	}
	if x.left != nil {
		return x.left.Max()
	}
	p := x.parent
	for p != nil && x == p.left {
		x = p
		p = x.parent
	}
	return p
}

// Next returns the next node (in in-order traversal order) of the given node x, or nil if no such
// node found.
func (t *Tree[K, V]) Next(x *Node[K, V]) *Node[K, V] {
	if x == nil {
		return nil
	}
	if x.right != nil {
		return x.right.Min()
	}
	p := x.parent
	for p != nil && x == p.right {
		x = p
		p = x.parent
	}
	return p
}

// Floor returns the floor node of the given key k, or nil if no such node found.
//
// A floor node is defined as the largest node which key is smaller than or equal to the given key
// k.
// A floor node may not be found, either because the tree is empty, or because all keys in the
// tree is larger than the given key k.
func (t *Tree[K, V]) Floor(k K) *Node[K, V] {
	var floor *Node[K, V]
	for x := t.root; x != nil; {
		r := t.cmp(k, x.key)
		if r < 0 {
			x = x.left
		} else if r > 0 {
			floor = x
			x = x.right
		} else {
			return x
		}
	}
	return floor
}

// Ceiling returns the ceiling node of key k, or nil if no such node found.
//
// A ceiling node is defined as the smallest node which key is larger than or equal to the given
// key k.
// A ceiling node may not be found, either because the tree is empty, or because all keys in the
// tree is smaller than the given key k.
func (t *Tree[K, V]) Ceiling(k K) *Node[K, V] {
	var ceiling *Node[K, V]
	for x := t.root; x != nil; {
		r := t.cmp(k, x.key)
		if r < 0 {
			ceiling = x
			x = x.left
		} else if r > 0 {
			x = x.right
		} else {
			return x
		}
	}
	return ceiling
}

// LevelOrder performs level-order traversal for tree, and returns a pair of slices (keys, values)
// as the result.
func (t *Tree[K, V]) LevelOrder() ([]K, []V) {
	keys := make([]K, 0, t.len)
	values := make([]V, 0, t.len)
	var q []*Node[K, V] // queue
	if t.root != nil {
		q = append(q, t.root)
	}
	for len(q) != 0 {
		x := q[0]
		q = q[1:]
		keys = append(keys, x.key)
		values = append(values, x.Value)
		if x.left != nil {
			q = append(q, x.left)
		}
		if x.right != nil {
			q = append(q, x.right)
		}
	}
	return keys, values
}

// PreOrder performs pre-order traversal for tree, and returns a pair of slices (keys, values) as
// the result.
func (t *Tree[K, V]) PreOrder() ([]K, []V) {
	keys := make([]K, 0, t.len)
	values := make([]V, 0, t.len)
	var s []*Node[K, V] // stack
	x := t.root
	for x != nil || len(s) != 0 {
		for x != nil {
			keys = append(keys, x.key)
			values = append(values, x.Value)
			s = append(s, x)
			x = x.left
		}
		// now, x == nil, and len(s) != 0
		x = s[len(s)-1]
		s = s[:len(s)-1]
		x = x.right
	}
	return keys, values
}

// InOrder performs in-order traversal for tree, and returns a pair of slices (keys, values) as
// the result.
func (t *Tree[K, V]) InOrder() ([]K, []V) {
	keys := make([]K, 0, t.len)
	values := make([]V, 0, t.len)
	var s []*Node[K, V] // stack
	x := t.root
	for x != nil || len(s) != 0 {
		for x != nil {
			s = append(s, x)
			x = x.left
		}
		// now, x == nil, and len(s) != 0
		x = s[len(s)-1]
		s = s[:len(s)-1]
		keys = append(keys, x.key)
		values = append(values, x.Value)
		x = x.right
	}
	return keys, values
}

// PostOrder performs post-order traversal for tree, and returns a pair of slices (keys, values)
// as the result.
func (t *Tree[K, V]) PostOrder() ([]K, []V) {
	keys := make([]K, 0, t.len)
	values := make([]V, 0, t.len)
	var s []*Node[K, V]  // stack
	var last *Node[K, V] // last visited node
	x := t.root
	for x != nil || len(s) != 0 {
		for x != nil {
			s = append(s, x)
			x = x.left
		}
		// now, x == nil, and len(s) != 0
		x = s[len(s)-1] // peek
		if x.right == nil || x.right == last {
			// there is no right child or return from right subtree, visit current node
			keys = append(keys, x.key)
			values = append(values, x.Value)
			last = x
			x = nil
			s = s[:len(s)-1] // pop
		} else {
			// right child exists, process right subtree
			x = x.right
		}
	}
	return keys, values
}

// PostOrderByReverse performs post-order traversal for tree by reverse the result slice, and
// returns a pair of slices (keys, values) as the result.
func (t *Tree[K, V]) PostOrderByReverse() ([]K, []V) {
	keys := make([]K, 0, t.len)
	values := make([]V, 0, t.len)
	// visit priority: root node -> right subtree -> left subtree
	var s []*Node[K, V] // stack
	x := t.root
	for x != nil || len(s) != 0 {
		for x != nil {
			keys = append(keys, x.key)
			values = append(values, x.Value)
			s = append(s, x)
			x = x.right
		}
		// now, x == nil, and len(s) != 0
		x = s[len(s)-1]
		s = s[:len(s)-1]
		x = x.left
	}
	// reverse: left subtree -> right subtree -> root node
	for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
		keys[i], keys[j] = keys[j], keys[i]
		values[i], values[j] = values[j], values[i]
	}
	return keys, values
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
		if x.left != nil {
			q = append(q, x.left)
		}
		if x.right != nil {
			q = append(q, x.right)
		}
		x.left = nil
		x.right = nil
		x.parent = nil
	}
	t.root = nil
	t.len = 0
}

// Range calls f sequentially for each key-value pair (k, v) present in tree in in-order traversal
// order.
// If f returns false, range stops the iteration.
func (t *Tree[K, V]) Range(f func(k K, v V) bool) {
	if f == nil {
		return
	}
	var s []*Node[K, V] // stack
	x := t.root
	for x != nil || len(s) != 0 {
		for x != nil {
			s = append(s, x)
			x = x.left
		}
		x = s[len(s)-1]
		s = s[:len(s)-1]
		if !f(x.key, x.Value) {
			break
		}
		x = x.right
	}
}
