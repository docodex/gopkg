package redblacktree

import (
	"math/rand/v2"
	"testing"
)

func keyOf[K comparable, V any](n *Node[K, V]) any {
	if n == nil {
		return "<nil>"
	}
	return n.key
}

// rbInvariants verifies the classic red-black tree properties rooted at x.
//
// Properties checked:
//  1. Every red node has a black parent (no two consecutive red nodes).
//  2. Every path from x down to a leaf (nil) contains the same number of
//     black nodes ("black-height" is uniform).
//  3. BST ordering: left.key < x.key < right.key under the tree's cmp.
//  4. Parent pointers are consistent.
//
// Returns the black-height of x (count of black nodes on any path from x
// to a descendant nil, not counting x itself).
func (t *Tree[K, V]) rbInvariants(tb testing.TB, x, parent *Node[K, V]) int {
	tb.Helper()
	if x == nil {
		// nil leaves are considered black.
		return 1
	}
	if x.parent != parent {
		tb.Fatalf("parent pointer mismatch at key=%v: got %v, want %v", x.key, keyOf(x.parent), keyOf(parent))
	}
	if x.color == red {
		if parent != nil && parent.color == red {
			tb.Fatalf("red-red violation: parent %v and child %v both red", parent.key, x.key)
		}
	}
	if l := x.left; l != nil && t.cmp(l.key, x.key) >= 0 {
		tb.Fatalf("BST order violated: left %v >= parent %v", l.key, x.key)
	}
	if r := x.right; r != nil && t.cmp(r.key, x.key) <= 0 {
		tb.Fatalf("BST order violated: right %v <= parent %v", r.key, x.key)
	}
	lh := t.rbInvariants(tb, x.left, x)
	rh := t.rbInvariants(tb, x.right, x)
	if lh != rh {
		tb.Fatalf("black-height mismatch at key=%v: left=%d right=%d", x.key, lh, rh)
	}
	if x.color == black {
		return lh + 1
	}
	return lh
}

func (t *Tree[K, V]) assertRBInvariants(tb testing.TB) {
	tb.Helper()
	if t.root != nil && t.root.color != black {
		tb.Fatalf("root must be black, got red")
	}
	t.rbInvariants(tb, t.root, nil)
}

func TestRBInvariants_AfterInsertsAndRemoves(t *testing.T) {
	sizes := []int{100, 500, 2000, 5000}
	for _, n := range sizes {
		tr := New[int, struct{}]()
		nums := rand.Perm(n)
		for i, v := range nums {
			tr.Insert(v, struct{}{})
			if i%97 == 0 {
				tr.assertRBInvariants(t)
			}
		}
		tr.assertRBInvariants(t)

		for i := 0; i < n/2; i++ {
			tr.Remove(nums[i])
			if i%53 == 0 {
				tr.assertRBInvariants(t)
			}
		}
		tr.assertRBInvariants(t)
	}
}

func TestRBInvariants_Randomized(t *testing.T) {
	r := rand.New(rand.NewPCG(3, 4))
	tr := New[int, struct{}]()
	alive := make(map[int]struct{})
	for i := range 5000 {
		switch r.IntN(3) {
		case 0, 1:
			k := r.IntN(2000)
			tr.Insert(k, struct{}{})
			alive[k] = struct{}{}
		case 2:
			if len(alive) == 0 {
				continue
			}
			for k := range alive {
				tr.Remove(k)
				delete(alive, k)
				break
			}
		}
		if i%200 == 0 {
			tr.assertRBInvariants(t)
			if tr.Len() != len(alive) {
				t.Fatalf("len mismatch: tree=%d expected=%d", tr.Len(), len(alive))
			}
		}
	}
	tr.assertRBInvariants(t)
}
