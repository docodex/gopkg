// package skiplist implements a skiplist data structure.
//
// In computer science, a skip list (or skiplist) is a probabilistic data structure that allows
// O(log n) average complexity for search as well as O(log n) average complexity for insertion
// within an ordered sequence of n elements. Thus it can get the best features of a sorted array
// (for searching) while maintaining a linked list-like structure that allows insertion, which is
// not possible with a static array. Fast search is made possible by maintaining a linked
// hierarchy of subsequences, with each successive subsequence skipping over fewer elements than
// the previous one. Searching starts in the sparsest subsequence until two consecutive elements
// have been found, one smaller and one larger than or equal to the element searched for. Via the
// linked hierarchy, these two elements link to elements of the next sparsest subsequence, where
// searching is continued until finally searching in the full sequence. The elements that are
// skipped over may be chosen probabilistically or deterministically, with the former being more
// common.
//
// Reference: https://en.wikipedia.org/wiki/Skip_list
package skiplist

import (
	"cmp"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/bytedance/gopkg/lang/fastrand"
	"github.com/docodex/gopkg/container"
)

// Element represents a key-value pair of a node.
type Element[K comparable, V any] struct {
	// The key used to compare elements.
	key K

	// The value stored with this element.
	Value V
}

// Key returns the key of element.
func (e *Element[K, V]) Key() K {
	return e.key
}

// Node is a node of a skiplist.
type Node[K comparable, V any] struct {
	// The element stored with this node.
	Element *Element[K, V]

	// Forward nodes in each level (forwards[0] is the forward node in base level).
	// To simplify the implementation, internally a skiplist l is implemented as a ring, such that
	// &l.root is both the forward node of the last skiplist node (l.MaxNode()) and the backward
	// node of the first skiplist node (l.MinNode()).
	forwards []*Node[K, V]

	// Number of nodes crossed (from this node to the forward node) in each level.
	// The length of spans is always equal to the length of forwards.
	spans []int

	// Backward node of this node in base level.
	// To simplify the implementation, internally a skiplist l is implemented as a ring, such that
	// &l.root is both the forward node of the last skiplist node (l.MaxNode()) and the backward
	// node of the first skiplist node (l.MinNode()).
	backward *Node[K, V]

	// The skiplist to which this node belongs.
	list *Skiplist[K, V]
}

// Next returns the next list node or nil.
func (n *Node[K, V]) Next() *Node[K, V] {
	if len(n.forwards) != 0 {
		if x := n.forwards[0]; n.list != nil && x != &n.list.root {
			return x
		}
	}
	return nil
}

// Prev returns the previous list node or nil.
func (n *Node[K, V]) Prev() *Node[K, V] {
	if x := n.backward; n.list != nil && x != &n.list.root {
		return x
	}
	return nil
}

const maxLevel = 32 // Should be enough for 2^64 elements

// Skiplist represents a skiplist.
type Skiplist[K comparable, V any] struct {
	root  Node[K, V]           // sentinel skiplist node
	len   int                  // current skiplist length excluding the sentinel node
	level int                  // current max level in skiplist
	cmp   container.Compare[K] // function to compare skiplist nodes
}

// New returns an initialized skiplist with [cmp.Compare] as the cmp function.
func New[K cmp.Ordered, V any]() *Skiplist[K, V] {
	return NewFunc[K, V](func(a, b K) int {
		return cmp.Compare(a, b)
	})
}

// NewFunc returns an initialized skiplist with the given function cmp as the cmp function.
func NewFunc[K comparable, V any](cmp container.Compare[K]) *Skiplist[K, V] {
	return new(Skiplist[K, V]).init(cmp)
}

// init initializes or clears skiplist l.
func (l *Skiplist[K, V]) init(cmp container.Compare[K]) *Skiplist[K, V] {
	l.root.forwards = make([]*Node[K, V], maxLevel)
	l.root.spans = make([]int, maxLevel)
	for i := range maxLevel {
		l.root.forwards[i] = &l.root
		l.root.spans[i] = 0 // spans initialized to l.len
	}
	l.root.backward = &l.root
	l.len = 0
	l.level = 1
	if cmp == nil {
		cmp = func(a, b K) int {
			// just to cover nil cmp error
			return 0
		}
	}
	l.cmp = cmp
	return l
}

// Len returns the number of nodes of skiplist t.
// The complexity is O(1).
func (l *Skiplist[K, V]) Len() int {
	return l.len
}

// Values returns all values in skiplist.
func (l *Skiplist[K, V]) Values() []V {
	values := make([]V, 0, l.len)
	l.Range(func(k K, v V) bool {
		values = append(values, v)
		return true
	})
	return values
}

// Values returns all keys in skiplist.
func (l *Skiplist[K, V]) Keys() []K {
	keys := make([]K, 0, l.len)
	l.Range(func(k K, v V) bool {
		keys = append(keys, k)
		return true
	})
	return keys
}

// String returns the string representation of skiplist.
// Ref: std fmt.Stringer.
func (l *Skiplist[K, V]) String() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "Skiplist(len:%d|level:%d): [", l.len, l.level)
	l.Range(func(k K, v V) bool {
		fmt.Fprintf(&buf, "(%v,%v)", k, v)
		return true
	})
	buf.WriteString("]")
	return buf.String()
}

// MarshalJSON marshals skiplist into valid JSON.
// Ref: std json.Marshaler.
func (l *Skiplist[K, V]) MarshalJSON() ([]byte, error) {
	m := make(map[K]V, l.len)
	l.Range(func(k K, v V) bool {
		m[k] = v
		return true
	})
	return json.Marshal(m)
}

// UnmarshalJSON unmarshals a JSON description of skiplist.
// The input can be assumed to be a valid encoding of a JSON value.
// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
// Ref: std json.Unmarshaler.
func (l *Skiplist[K, V]) UnmarshalJSON(data []byte) error {
	var m map[K]V
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	l.Clear()
	for k, v := range m {
		l.Insert(k, v)
	}
	return nil
}

// Insert inserts a new node with the given key-value pair (k, v) to skiplist, or update the key
// and value if the given key k already exists in skiplist.
func (l *Skiplist[K, V]) Insert(k K, v V) {
	update := make([]*Node[K, V], maxLevel) // previous nodes of target position in each level
	rank := make([]int, maxLevel)           // nodes crossed by (distance to root)
	x := &l.root
	for i := l.level - 1; i >= 0; i-- {
		// store rank crossed from root to reach the insertion position in each level
		if i == l.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		// find the right position to insert in each level
		for x.forwards[i] != &l.root {
			r := l.cmp(x.forwards[i].Element.key, k)
			if r > 0 {
				break
			}
			if r == 0 {
				// duplicated key found, update the key and value
				x.forwards[i].Element.key = k
				x.forwards[i].Element.Value = v
				return
			}
			rank[i] += x.spans[i]
			x = x.forwards[i]
		}
		update[i] = x
	}
	// get random level n, if n > l.level, then extend skiplist level
	n := l.randomLevel()
	if n > l.level {
		for i := l.level; i < n; i++ {
			rank[i] = 0
			update[i] = &l.root
			// span of extended levels initialized to l.len
			update[i].spans[i] = l.len
		}
		l.level = n
	}
	// create and insert new node
	x = &Node[K, V]{
		Element: &Element[K, V]{
			key:   k,
			Value: v,
		},
		forwards: make([]*Node[K, V], n),
		spans:    make([]int, n),
		backward: nil,
		list:     l,
	}
	for i := range n {
		x.forwards[i] = update[i].forwards[i]
		update[i].forwards[i] = x
		// update span covered by update[i] as x is inserted here
		delta := rank[0] - rank[i]
		x.spans[i] = update[i].spans[i] - delta
		update[i].spans[i] = delta + 1
	}
	// if n < l.level, increment span for untouched levels
	for i := n; i < l.level; i++ {
		update[i].spans[i]++
	}
	// set backward node (just base level needed)
	x.backward = update[0]
	x.forwards[0].backward = x
	l.len++
}

const threshold = math.MaxUint32 >> 2 // P = 0.25

// randomLevel returns a random level number which is not greater than the max level.
func (l *Skiplist[K, V]) randomLevel() int {
	level := 1
	for fastrand.Uint32() < threshold && level < maxLevel {
		level++
	}
	return level
}

// Remove removes the node which key equals to the given key k from skiplist and returns the
// element of that node.
func (l *Skiplist[K, V]) Remove(k K) *Element[K, V] {
	update := make([]*Node[K, V], maxLevel) // previous nodes of target node in each level
	x := &l.root
	for i := l.level - 1; i >= 0; i-- {
		for x.forwards[i] != &l.root && l.cmp(x.forwards[i].Element.key, k) < 0 {
			x = x.forwards[i]
		}
		update[i] = x
	}
	var e *Element[K, V]
	x = x.forwards[0]
	if x != &l.root && l.cmp(x.Element.key, k) == 0 {
		e = x.Element
		l.remove(x, update)
	}
	return e
}

// RemoveByRank removes the node which rank equals to the given rank from skiplist and returns the
// element of that node.
func (l *Skiplist[K, V]) RemoveByRank(rank int) *Element[K, V] {
	update := make([]*Node[K, V], maxLevel) // previous nodes of target node in each level
	x := &l.root
	for i := l.level - 1; i >= 0; i-- {
		for x.forwards[i] != &l.root && rank-x.spans[i] > 0 {
			rank -= x.spans[i]
			x = x.forwards[i]
		}
		update[i] = x
	}
	var e *Element[K, V]
	rank-- // x.spans[0] should be 1 for non-empty skiplist
	x = x.forwards[0]
	if x != &l.root {
		e = x.Element
		l.remove(x, update)
	}
	return e
}

// RemoveRange removes the nodes which keys is within the given range [k1, k2) in which k1 is
// inclusive and k2 is exclusive from skiplist and returns the elements of those nodes.
func (l *Skiplist[K, V]) RemoveRange(k1, k2 K) []*Element[K, V] {
	update := make([]*Node[K, V], maxLevel) // previous nodes of target range in each level
	x := &l.root
	for i := l.level - 1; i >= 0; i-- {
		for x.forwards[i] != &l.root && l.cmp(x.forwards[i].Element.key, k1) < 0 {
			x = x.forwards[i]
		}
		update[i] = x
	}
	var removed []*Element[K, V]
	x = x.forwards[0] // x is the first node within range [k1, k2)
	for x != &l.root && l.cmp(x.Element.key, k2) < 0 {
		y := x.forwards[0]
		removed = append(removed, x.Element)
		l.remove(x, update)
		x = y
	}
	return removed
}

// RemoveRangeByRank removes the nodes which rank is within the given range [rank1, rank2) which
// rank1 is inclusive and rank2 is exclusive from skiplist and returns the elements of those
// nodes.
func (l *Skiplist[K, V]) RemoveRangeByRank(rank1, rank2 int) []*Element[K, V] {
	update := make([]*Node[K, V], maxLevel) // previous nodes of target range in each level
	rank := 0
	x := &l.root
	for i := l.level - 1; i >= 0; i-- {
		for x.forwards[i] != &l.root && rank+x.spans[i] < rank1 {
			rank += x.spans[i]
			x = x.forwards[i]
		}
		update[i] = x
	}
	var removed []*Element[K, V]
	rank++            // x.spans[0] should be 1 for non-empty skiplist
	x = x.forwards[0] // x is the first node within range [rank1, rank2)
	for rank < rank2 && x != &l.root {
		y := x.forwards[0]
		removed = append(removed, x.Element)
		l.remove(x, update)
		x = y
		rank++
	}
	return removed
}

// remove removes the nodes x from skiplist.
// The given node x must not be nil, while the given nodes slice update is the previous node of
// node x in each level.
func (l *Skiplist[K, V]) remove(x *Node[K, V], update []*Node[K, V]) {
	for i := 0; i < l.level; i++ {
		if update[i].forwards[i] == x {
			update[i].forwards[i] = x.forwards[i]
			update[i].spans[i] += x.spans[i] - 1
		} else {
			update[i].spans[i]--
		}
	}
	x.forwards[0].backward = x.backward
	clear(x.forwards)
	x.forwards = nil
	x.backward = nil
	x.list = nil
	for l.level > 1 && l.root.forwards[l.level-1] == &l.root {
		l.level--
	}
	l.len--
}

// Get returns the element which key equals to the given key k.
// Get also returns the rank of the returned element in skiplist.
func (l *Skiplist[K, V]) Get(k K) (*Element[K, V], int) {
	rank := 0
	x := &l.root
	for i := l.level - 1; i >= 0; i-- {
		for x.forwards[i] != &l.root && l.cmp(x.forwards[i].Element.key, k) <= 0 {
			rank += x.spans[i]
			x = x.forwards[i]
		}
		if x != &l.root && l.cmp(x.Element.key, k) == 0 {
			return x.Element, rank
		}
	}
	return nil, 0
}

// GetByRank returns the element which rank equals to the given rank.
func (l *Skiplist[K, V]) GetByRank(rank int) *Element[K, V] {
	x := &l.root
	for i := l.level - 1; i >= 0; i-- {
		for x.forwards[i] != &l.root && rank-x.spans[i] >= 0 {
			rank -= x.spans[i]
			x = x.forwards[i]
		}
		if rank == 0 && x != &l.root {
			return x.Element
		}
	}
	return nil
}

// GetRange returns the elements which keys is within the given range [k1, k2) in which k1 is
// inclusive and k2 is exclusive.
func (l *Skiplist[K, V]) GetRange(k1, k2 K) []*Element[K, V] {
	x := &l.root
	for i := l.level - 1; i >= 0; i-- {
		for x.forwards[i] != &l.root && l.cmp(x.forwards[i].Element.key, k1) < 0 {
			x = x.forwards[i]
		}
	}
	var elements []*Element[K, V]
	x = x.forwards[0] // x is the first node within range [k1, k2)
	for x != &l.root && l.cmp(x.Element.key, k2) < 0 {
		elements = append(elements, x.Element)
		x = x.forwards[0]
	}
	return elements
}

// GetRangeByRank removes the nodes which rank is within the given range [rank1, rank2) which
// rank1 is inclusive and rank2 is exclusive.
func (l *Skiplist[K, V]) GetRangeByRank(rank1, rank2 int) []*Element[K, V] {
	rank := 0
	x := &l.root
	for i := l.level - 1; i >= 0; i-- {
		for x.forwards[i] != &l.root && rank+x.spans[i] < rank1 {
			rank += x.spans[i]
			x = x.forwards[i]
		}
	}
	var elements []*Element[K, V]
	rank++            // x.spans[0] should be 1 for non-empty skiplist
	x = x.forwards[0] // x is the first node within range [rank1, rank2)
	for rank < rank2 && x != &l.root {
		elements = append(elements, x.Element)
		x = x.forwards[0]
		rank++
	}
	return elements
}

// MinNode returns the node which key is the minimum key of skiplist.
func (l *Skiplist[K, V]) MinNode() *Node[K, V] {
	if l.len == 0 {
		return nil
	}
	return l.root.forwards[0]
}

// MaxNode returns the node which key is the maximum key of skiplist.
func (l *Skiplist[K, V]) MaxNode() *Node[K, V] {
	if l.len == 0 {
		return nil
	}
	return l.root.backward
}

// Min returns the element which key is the minimum key of skiplist.
func (l *Skiplist[K, V]) Min() *Element[K, V] {
	if x := l.MinNode(); x != nil {
		return x.Element
	}
	return nil
}

// Max returns the element which key is the maximum key of skiplist.
func (l *Skiplist[K, V]) Max() *Element[K, V] {
	if x := l.MaxNode(); x != nil {
		return x.Element
	}
	return nil
}

// MinNodeInRange returns the node which key is the minimum key within the given range [k1, k2) in
// which k1 is inclusive and k2 is exclusive.
func (l *Skiplist[K, V]) MinNodeInRange(k1, k2 K) *Node[K, V] {
	x := &l.root
	for i := l.level - 1; i >= 0; i-- {
		for x.forwards[i] != &l.root && l.cmp(x.forwards[i].Element.key, k1) < 0 {
			x = x.forwards[i]
		}
	}
	x = x.forwards[0] // x.Element.key should be greater than or equal to k1
	if x != &l.root && l.cmp(x.Element.key, k2) < 0 {
		return x
	}
	return nil
}

// MaxNodeInRange returns the node which key is the maximum key within the given range [k1, k2) in
// which k1 is inclusive and k2 is exclusive.
func (l *Skiplist[K, V]) MaxNodeInRange(k1, k2 K) *Node[K, V] {
	x := &l.root
	for i := l.level - 1; i >= 0; i-- {
		for x.forwards[i] != &l.root && l.cmp(x.forwards[i].Element.key, k2) < 0 {
			x = x.forwards[i]
		}
	}
	// x.Element.key should be less than k2
	if x != &l.root && l.cmp(x.Element.key, k1) >= 0 {
		return x
	}
	return nil
}

// MinInRange returns the element which key is the minimum key within the given range [k1, k2) in
// which k1 is inclusive and k2 is exclusive.
func (l *Skiplist[K, V]) MinInRange(k1, k2 K) *Element[K, V] {
	if x := l.MinNodeInRange(k1, k2); x != nil {
		return x.Element
	}
	return nil
}

// MaxInRange returns the element which key is the maximum key within the given range [k1, k2) in
// which k1 is inclusive and k2 is exclusive.
func (l *Skiplist[K, V]) MaxInRange(k1, k2 K) *Element[K, V] {
	if x := l.MaxNodeInRange(k1, k2); x != nil {
		return x.Element
	}
	return nil
}

// Clear removes all nodes in skiplist.
func (l *Skiplist[K, V]) Clear() {
	for x := l.root.forwards[0]; x != &l.root; {
		y := x.forwards[0]
		clear(x.forwards)
		x.forwards = nil
		x.backward = nil
		x.list = nil
		x = y
	}
	l.init(l.cmp)
}

// Range calls f sequentially for each key-value pair (k, v) present in skiplist.
// If f returns false, range stops the iteration.
func (l *Skiplist[K, V]) Range(f func(k K, v V) bool) {
	if f == nil {
		return
	}
	for x := l.root.forwards[0]; x != &l.root; x = x.forwards[0] {
		if !f(x.Element.key, x.Element.Value) {
			break
		}
	}
}
