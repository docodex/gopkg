package btree

import (
	"math/rand/v2"
	"testing"
)

func (t *Tree[K, V]) btreeInvariants(tb testing.TB, x *Node[K, V], isRoot bool, lo, hi *K) int {
	tb.Helper()
	if x == nil {
		return 0
	}
	n := len(x.Entries)
	if isRoot {
		if t.root != nil && n == 0 {
			tb.Fatalf("non-empty root must have >= 1 entry")
		}
	} else {
		if n < t.minSize || n > t.maxSize {
			tb.Fatalf("entry count out of range: n=%d, want [%d,%d]", n, t.minSize, t.maxSize)
		}
	}
	// keys strictly increasing
	for i := 1; i < n; i++ {
		if t.cmp(x.Entries[i-1].key, x.Entries[i].key) >= 0 {
			tb.Fatalf("keys not strictly increasing in node: %v then %v", x.Entries[i-1].key, x.Entries[i].key)
		}
	}
	// slot intervals
	if lo != nil && n > 0 && t.cmp(*lo, x.Entries[0].key) >= 0 {
		tb.Fatalf("left slot violated: lo=%v, first=%v", *lo, x.Entries[0].key)
	}
	if hi != nil && n > 0 && t.cmp(x.Entries[n-1].key, *hi) >= 0 {
		tb.Fatalf("right slot violated: last=%v, hi=%v", x.Entries[n-1].key, *hi)
	}

	if len(x.children) == 0 {
		return 1
	}
	if len(x.children) != n+1 {
		tb.Fatalf("children count mismatch: %d entries, %d children", n, len(x.children))
	}

	// recurse into children; track slot bounds
	var firstDepth int
	for i, c := range x.children {
		if c.parent != x {
			tb.Fatalf("child %d has wrong parent pointer", i)
		}
		var childLo, childHi *K
		if i > 0 {
			childLo = &x.Entries[i-1].key
		} else {
			childLo = lo
		}
		if i < n {
			childHi = &x.Entries[i].key
		} else {
			childHi = hi
		}
		d := t.btreeInvariants(tb, c, false, childLo, childHi)
		if i == 0 {
			firstDepth = d
		} else if d != firstDepth {
			tb.Fatalf("unequal leaf depths under node: child0=%d childN=%d", firstDepth, d)
		}
	}
	return firstDepth + 1
}

func (t *Tree[K, V]) assertBTreeInvariants(tb testing.TB) {
	tb.Helper()
	if t.root == nil {
		if t.len != 0 {
			tb.Fatalf("nil root but len=%d", t.len)
		}
		return
	}
	t.btreeInvariants(tb, t.root, true, nil, nil)
}

func TestBTreeInvariants_AfterInsertsAndRemoves(t *testing.T) {
	orders := []int{3, 4, 5, 7, 10}
	sizes := []int{50, 200, 1000}
	for _, m := range orders {
		for _, n := range sizes {
			tr := New[int, struct{}](m)
			nums := rand.Perm(n)
			for i, v := range nums {
				tr.Insert(v, struct{}{})
				if i%31 == 0 {
					tr.assertBTreeInvariants(t)
				}
			}
			tr.assertBTreeInvariants(t)

			for i := 0; i < n/2; i++ {
				tr.Remove(nums[i])
				if i%23 == 0 {
					tr.assertBTreeInvariants(t)
				}
			}
			tr.assertBTreeInvariants(t)
		}
	}
}

func TestBTreeInvariants_Randomized(t *testing.T) {
	r := rand.New(rand.NewPCG(5, 6))
	for _, m := range []int{3, 5, 8} {
		tr := New[int, struct{}](m)
		alive := make(map[int]struct{})
		for i := range 3000 {
			switch r.IntN(3) {
			case 0, 1:
				k := r.IntN(1500)
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
			if i%150 == 0 {
				tr.assertBTreeInvariants(t)
				if tr.Len() != len(alive) {
					t.Fatalf("len mismatch (m=%d): tree=%d expected=%d", m, tr.Len(), len(alive))
				}
			}
		}
		tr.assertBTreeInvariants(t)
	}
}
