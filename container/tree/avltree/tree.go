// Package avltree implements an AVL balanced binary tree.
//
// References:
// - https://en.wikipedia.org/wiki/AVL_tree
// - https://en.wikipedia.org/wiki/Binary_search_tree
package avltree

import (
	"cmp"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docodex/gopkg/container"
)

// Node is a node of a binary tree.
type Node[K comparable, V any] struct {
	// The key used to compare nodes.
	key K

	// The value stored with this node.
	Value V

	// Left and right children nodes of this node in tree.
	left, right *Node[K, V]

	// The height of this node, represents the height of subtree with current node as the root.
	// Height of leaf node is the default value of int (0), while height of nil node is -1.
	height int
}

// Left returns the key of node.
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

// Tree represents an avl tree.
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

// height returns the height of node x
func (t *Tree[K, V]) height(x *Node[K, V]) int {
	if x == nil {
		// height of nil node is -1
		return -1
	}
	return x.height
}

// updateHeight updates the height of node x
func (t *Tree[K, V]) updateHeight(x *Node[K, V]) {
	if x != nil {
		// as height of nil node is -1, height of leaf node is 0: (-1) + 1
		x.height = max(t.height(x.left), t.height(x.right)) + 1
	}
}

// rightRotate do right rotate operation
func (t *Tree[K, V]) rightRotate(x *Node[K, V]) *Node[K, V] {
	y := x.left
	// after right rotate, the right child of node y would be moved as the left child of node x
	x.left = y.right
	// rotate node x to the right around node y
	y.right = x
	// update the height of nodes x and y
	t.updateHeight(x)
	t.updateHeight(y)
	// return the root of the subtree after rotation
	return y
}

// leftRotate do left rotate operation
func (t *Tree[K, V]) leftRotate(x *Node[K, V]) *Node[K, V] {
	y := x.right
	// after left rotate, the left child of node y would be moved as the right child of node x
	x.right = y.left
	// rotate node x to the left around node y
	y.left = x
	// update the height of nodes x and y
	t.updateHeight(x)
	t.updateHeight(y)
	// return the root of the subtree after rotation
	return y
}

// balanceFactor returns the balance factor of node x.
func (t *Tree[K, V]) balanceFactor(x *Node[K, V]) int {
	if x == nil {
		// balance factor of nil node is 0
		return 0
	}
	// node balance factor = left subtree height - right subtree height
	// so, balance factor f of avl tree nodes should satisfy: -1 <= f <= 1
	return t.height(x.left) - t.height(x.right)
}

// rotate performs rotation operation to restore balance to the subtree
func (t *Tree[K, V]) rotate(x *Node[K, V]) *Node[K, V] {
	f := t.balanceFactor(x)
	// left-leaning subtree
	if f > 1 {
		if t.balanceFactor(x.left) >= 0 {
			// right rotate
			return t.rightRotate(x)
		} else {
			// left rotation (x.left) first, then right rotation (x)
			x.left = t.leftRotate(x.left)
			return t.rightRotate(x)
		}
	}
	// right-leaning subtree
	if f < -1 {
		if t.balanceFactor(x.right) <= 0 {
			// left rotate
			return t.leftRotate(x)
		} else {
			// right rotation (x.right) first, then left rotation (x)
			x.right = t.rightRotate(x.right)
			return t.leftRotate(x)
		}
	}
	// balanced tree, no rotation needed
	return x
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
	buf.WriteString("AVLTree\n")
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
	fmt.Fprintf(buf, "%v:%v\n", x.key, x.Value)
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
	t.root = t.insert(t.root, k, v)
}

// insert inserts a new node with the given key-value pair (k, v) to subtree with node x as the
// root.
func (t *Tree[K, V]) insert(x *Node[K, V], k K, v V) *Node[K, V] {
	if x == nil {
		t.len++
		return &Node[K, V]{
			key:    k,
			Value:  v,
			left:   nil,
			right:  nil,
			height: 0,
		}
	}
	// find the right position and do insert
	r := t.cmp(k, x.key)
	if r < 0 {
		x.left = t.insert(x.left, k, v)
	} else if r > 0 {
		x.right = t.insert(x.right, k, v)
	} else {
		// duplicated value found, update the key and value
		x.key = k
		x.Value = v
		return x
	}
	// update the height of node x
	t.updateHeight(x)
	// perform rotation operation to restore balance to the subtree
	x = t.rotate(x)
	// return the root node of the subtree
	return x
}

// Remove removes the node which key equals to the given key k from tree.
func (t *Tree[K, V]) Remove(k K) {
	t.root = t.remove(t.root, k)
}

// remove removes the node with the given key k from subtree with node x as the root.
func (t *Tree[K, V]) remove(x *Node[K, V], k K) *Node[K, V] {
	if x == nil {
		return nil
	}
	// find and remove the node
	r := t.cmp(k, x.key)
	if r < 0 {
		x.left = t.remove(x.left, k)
	} else if r > 0 {
		x.right = t.remove(x.right, k)
	} else {
		if x.left == nil || x.right == nil {
			y := x.left
			if x.right != nil {
				y = x.right
			}
			// 1. if y == nil, node x has no children, remove the node and return, let x = nil
			// 2. if y != nil, node x has only 1 child, remove the node and return y, let x = y
			t.len--
			return y
		} else {
			// node x has both 2 children
			// remove the next node in in-order traversal and replace the current node with it
			y := x.right
			for y.left != nil {
				y = y.left
			}
			x.right = t.remove(x.right, y.key)
			x.key = y.key
			x.Value = y.Value
		}
	}
	// update the height of node x
	t.updateHeight(x)
	// perform rotation operation to restore balance to the subtree
	x = t.rotate(x)
	// return the root node of the subtree
	return x
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
	t.preOrder(t.root, &keys, &values)
	return keys, values
}

// preOrder performs pre-order traversal for subtree with node x as the root, and append the
// result to a pair of slices (keys, values).
func (t *Tree[K, V]) preOrder(x *Node[K, V], keys *[]K, values *[]V) {
	if x != nil {
		// visit priority: root node -> left subtree -> right subtree
		*keys = append(*keys, x.key)
		*values = append(*values, x.Value)
		t.preOrder(x.left, keys, values)
		t.preOrder(x.right, keys, values)
	}
}

// InOrder performs in-order traversal for tree, and returns a pair of slices (keys, values) as
// the result.
func (t *Tree[K, V]) InOrder() ([]K, []V) {
	keys := make([]K, 0, t.len)
	values := make([]V, 0, t.len)
	t.inOrder(t.root, &keys, &values)
	return keys, values
}

// inOrder performs in-order traversal for subtree with node x as the root, and append the result
// to a pair of slices (keys, values).
func (t *Tree[K, V]) inOrder(x *Node[K, V], keys *[]K, values *[]V) {
	if x != nil {
		// visit priority: left subtree -> root node -> right subtree
		t.inOrder(x.left, keys, values)
		*keys = append(*keys, x.key)
		*values = append(*values, x.Value)
		t.inOrder(x.right, keys, values)
	}
}

// PostOrder performs post-order traversal for tree, and returns a pair of slices (keys, values)
// as the result.
func (t *Tree[K, V]) PostOrder() ([]K, []V) {
	keys := make([]K, 0, t.len)
	values := make([]V, 0, t.len)
	t.postOrder(t.root, &keys, &values)
	return keys, values
}

// postOrder performs post-order traversal for subtree with node x as the root, and append the
// result to a pair of slices (keys, values).
func (t *Tree[K, V]) postOrder(x *Node[K, V], keys *[]K, values *[]V) {
	if x != nil {
		// visit priority: left subtree -> right subtree -> root node
		t.postOrder(x.left, keys, values)
		t.postOrder(x.right, keys, values)
		*keys = append(*keys, x.key)
		*values = append(*values, x.Value)
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
		if x.left != nil {
			q = append(q, x.left)
		}
		if x.right != nil {
			q = append(q, x.right)
		}
		x.left = nil
		x.right = nil
	}
	t.root = nil
	t.len = 0
}

// Range calls f sequentially for each key-value pair (k, v) present in tree in in-order traversal
// order. If f returns false, range stops the iteration.
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
		// now, x == nil, and len(s) != 0
		x = s[len(s)-1]
		s = s[:len(s)-1]
		if !f(x.key, x.Value) {
			break
		}
		x = x.right
	}
}
