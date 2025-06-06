package avltree_test

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/tree/avltree"
	"github.com/stretchr/testify/assert"
)

func TestAVLTree(t *testing.T) {
	testcases := []int{100, 200, 1000, 10000}
	for _, n := range testcases {
		t1 := avltree.New[int, struct{}]()
		nums := rand.Perm(n)
		for _, v := range nums {
			t1.Insert(v, struct{}{})
		}

		bytes, _ := t1.MarshalJSON()
		_ = t1.UnmarshalJSON(bytes)

		var count int
		t1.Range(func(k int, v struct{}) bool {
			count++
			return true
		})
		fmt.Println("count:", count)
		assert.Equal(t, count, len(nums))
		assert.Equal(t, count, t1.Len())
		assert.Equal(t, count, t1.Root().Len())

		rets, _ := t1.InOrder()
		if !sort.IntsAreSorted(rets) {
			t.Error("Error with Push")
		}

		if res := t1.Min(); res == nil || res.Key() != rets[0] {
			t.Errorf("Error with Min, get %d, want: %d", res.Value, rets[0])
		}

		if res := t1.Max(); res == nil || res.Key() != rets[n-1] {
			t.Errorf("Error with Max, get %d, want: %d", res.Value, rets[n-1])
		}

		for i := 0; i < n-1; i++ {
			t1.Remove(nums[i])
			rets, _ = t1.InOrder()
			if !sort.IntsAreSorted(rets) {
				t.Errorf("Error With Remove")
			}
		}
	}
}

func TestAVLInit(t *testing.T) {
	testcases := []int{100, 200, 1000, 10000}
	for _, n := range testcases {
		nums := rand.Perm(n)
		m := make(map[int]struct{}, len(nums))
		for _, v := range nums {
			m[v] = struct{}{}
		}
		v, err := json.Marshal(m)
		assert.Nil(t, err)
		t1 := avltree.New[int, struct{}]()
		err = t1.UnmarshalJSON(v)
		assert.Nil(t, err)

		v, err = t1.MarshalJSON()
		assert.Nil(t, err)
		t2 := avltree.New[int, struct{}]()
		err = t2.UnmarshalJSON(v)
		assert.Nil(t, err)

		rets, _ := t1.InOrder()
		if !sort.IntsAreSorted(rets) {
			t.Error("Error with Push")
		}

		if res := t1.Min(); res == nil || res.Key() != rets[0] {
			t.Errorf("Error with Min, get %d, want: %d", res.Value, rets[0])
		}

		if res := t1.Max(); res == nil || res.Key() != rets[n-1] {
			t.Errorf("Error with Max, get %d, want: %d", res.Value, rets[n-1])
		}

		for i := 0; i < n-1; i++ {
			t1.Remove(nums[i])
			rets, _ = t1.InOrder()
			if !sort.IntsAreSorted(rets) {
				t.Errorf("Error With Remove")
			}
		}

		rets, _ = t2.InOrder()
		if !sort.IntsAreSorted(rets) {
			t.Error("Error with Push")
		}

		if res := t2.Min(); res == nil || res.Key() != rets[0] {
			t.Errorf("Error with Min, get %d, want: %d", res.Value, rets[0])
		}

		if res := t2.Max(); res == nil || res.Key() != rets[n-1] {
			t.Errorf("Error with Max, get %d, want: %d", res.Value, rets[n-1])
		}

		for i := 0; i < n-1; i++ {
			t2.Remove(nums[i])
			rets, _ = t2.InOrder()
			if !sort.IntsAreSorted(rets) {
				t.Errorf("Error With Remove")
			}
		}
	}
}

func TestAVLSearch(t *testing.T) {
	type Entry struct {
		key, value int
	}

	t1 := avltree.New[int, int]()
	for i := range 10 {
		t1.Insert(i, i+10000)
	}
	assert.False(t, t1.Len() == 0)
	assert.Equal(t, 10, t1.Len())
	assert.Equal(t, t1.Len(), t1.Root().Len())

	fmt.Println(t1)
	buf, _ := t1.MarshalJSON()
	fmt.Println(string(buf))

	for i := range 10 {
		x := t1.Search(i)
		assert.Equal(t, i+10000, x.Value)
	}
}

func TestAVLDelete(t *testing.T) {
	t1 := avltree.New[int, int]()
	m := make(map[int]int)
	for i := range 1000 {
		t1.Insert(i, i+10000)
		m[i] = i + 10000
	}
	count := 1000
	for k, v := range m {
		assert.Equal(t, t1.Len(), t1.Root().Len())
		x := t1.Search(k)
		assert.Equal(t, v, x.Value)
		t1.Remove(k)
		assert.Nil(t, t1.Search(k))
		count--
		assert.Equal(t, count, t1.Len())
	}
}

func TestAVLInsert(t *testing.T) {
	t.Run("LLRotation-Test", func(t *testing.T) {
		t1 := avltree.New[int, struct{}]()
		t1.Insert(5, struct{}{})
		t1.Insert(4, struct{}{})
		t1.Insert(3, struct{}{})

		root := t1.Root()
		if root.Key() != 4 {
			t.Errorf("Root should have key = 4, not %v", root.Key())
		}
		if root.Left().Key() != 3 {
			t.Errorf("Left child should have key = 3")
		}
		if root.Right().Key() != 5 {
			t.Errorf("Right child should have key = 5")
		}
	})

	t.Run("LRRotation-Test", func(t *testing.T) {
		t1 := avltree.New[int, struct{}]()
		t1.Insert(5, struct{}{})
		t1.Insert(4, struct{}{})
		t1.Insert(3, struct{}{})

		root := t1.Root()
		if root.Key() != 4 {
			t.Errorf("Root should have key = 4, not %v", root.Key())
		}
		if root.Left().Key() != 3 {
			t.Errorf("Left child should have key = 3")
		}
		if root.Right().Key() != 5 {
			t.Errorf("Right child should have key = 5")
		}
	})

	t.Run("RRRotation-Test", func(t *testing.T) {
		t1 := avltree.New[int, struct{}]()
		t1.Insert(3, struct{}{})
		t1.Insert(4, struct{}{})
		t1.Insert(5, struct{}{})

		root := t1.Root()
		if root.Key() != 4 {
			t.Errorf("Root should have key = 4, not %v", root.Key())
		}
		if root.Left().Key() != 3 {
			t.Errorf("Left child should have key = 3")
		}
		if root.Right().Key() != 5 {
			t.Errorf("Right child should have key = 5")
		}
	})

	t.Run("RLRotation-Test", func(t *testing.T) {
		tree := avltree.New[int, struct{}]()
		tree.Insert(3, struct{}{})
		tree.Insert(5, struct{}{})
		tree.Insert(4, struct{}{})

		root := tree.Root()
		if root.Key() != 4 {
			t.Errorf("Root should have key = 4")
		}
		if root.Left().Key() != 3 {
			t.Errorf("Left child should have key = 3")
		}
		if root.Right().Key() != 5 {
			t.Errorf("Right child should have key = 5")
		}
	})
}

func TestAVLRemove(t *testing.T) {
	t.Run("LLRotation-Test", func(t *testing.T) {
		t1 := avltree.New[int, struct{}]()
		t1.Remove(5)

		t1.Insert(5, struct{}{})
		t1.Insert(4, struct{}{})
		t1.Insert(3, struct{}{})
		t1.Insert(2, struct{}{})

		fmt.Println(t1.String())

		t1.Remove(5)
		t1.Remove(50)

		fmt.Println(t1.String())

		root := t1.Root()
		if root.Key() != 3 {
			t.Errorf("Root should have key = 3")
		}
		if root.Left().Key() != 2 {
			t.Errorf("Left child should have key = 2")
		}
		if root.Right().Key() != 4 {
			t.Errorf("Right child should have key = 4")
		}
	})

	t.Run("LRRotation-Test", func(t *testing.T) {
		t1 := avltree.New[int, struct{}]()

		t1.Insert(10, struct{}{})
		t1.Insert(8, struct{}{})
		t1.Insert(8, struct{}{})
		t1.Insert(6, struct{}{})
		t1.Insert(7, struct{}{})

		t1.Remove(10)
		t1.Remove(5)

		root := t1.Root()
		if root.Key() != 7 {
			t.Errorf("Root should have key = 7")
		}
		if root.Left().Key() != 6 {
			t.Errorf("Left child should have key = 6")
		}
		if root.Right().Key() != 8 {
			t.Errorf("Right child should have key = 8")
		}
	})

	t.Run("RRRotation-Test", func(t *testing.T) {
		t1 := avltree.New[int, struct{}]()

		t1.Insert(2, struct{}{})
		t1.Insert(3, struct{}{})
		t1.Insert(3, struct{}{})
		t1.Insert(4, struct{}{})
		t1.Insert(5, struct{}{})

		t1.Remove(2)
		t1.Remove(15)

		root := t1.Root()
		if root.Key() != 4 {
			t.Errorf("Root should have key = 4")
		}
		if root.Left().Key() != 3 {
			t.Errorf("Left child should have key = 3")
		}
		if root.Right().Key() != 5 {
			t.Errorf("Right child should have key = 5")
		}
	})

	t.Run("RLRotation-Test", func(t *testing.T) {
		t1 := avltree.New[int, struct{}]()

		t1.Insert(7, struct{}{})
		t1.Insert(6, struct{}{})
		t1.Insert(6, struct{}{})
		t1.Insert(9, struct{}{})
		t1.Insert(8, struct{}{})

		t1.Remove(6)

		root := t1.Root()
		if root.Key() != 8 {
			t.Errorf("Root should have key = 8")
		}
		if root.Left().Key() != 7 {
			t.Errorf("Left child should have key = 7")
		}
		if root.Right().Key() != 9 {
			t.Errorf("Right child should have key = 9")
		}
	})

	t.Run("Random Test", func(t *testing.T) {
		nums := []int{100, 500, 1000, 10_000}
		for _, n := range nums {
			t1 := avltree.New[int, struct{}]()
			nums := rand.Perm(n)
			for _, v := range nums {
				t1.Insert(v, struct{}{})
			}

			rets, _ := t1.InOrder()
			if !sort.IntsAreSorted(rets) {
				t.Error("Error with Push")
			}

			if res := t1.Min(); res == nil || res.Key() != rets[0] {
				t.Errorf("Error with Min, get %v, want: %d", res, rets[0])
			}

			if res := t1.Max(); res == nil || res.Key() != rets[n-1] {
				t.Errorf("Error with Max, get %v, want: %d", res, rets[n-1])
			}

			for i := range n {
				t1.Remove(nums[i])
				rets, _ = t1.InOrder()
				if !sort.IntsAreSorted(rets) {
					t.Errorf("Error With Remove")
				}
			}
		}
	})
}

func TestAVLTreeSearch(t *testing.T) {
	t1 := avltree.New[int, string]()
	if actualValue := t1.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
	if actualValue := t1.Search(2); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	t1.Insert(1, "x") // 1->x
	t1.Insert(2, "b") // 1->x, 2->b (in order)
	t1.Insert(1, "a") // 1->a, 2->b (in order, replacement)
	t1.Insert(3, "c") // 1->a, 2->b, 3->c (in order)
	t1.Insert(4, "d") // 1->a, 2->b, 3->c, 4->d (in order)
	t1.Insert(5, "e") // 1->a, 2->b, 3->c, 4->d, 5->e (in order)
	t1.Insert(6, "f") // 1->a, 2->b, 3->c, 4->d, 5->e, 6->f (in order)

	//  AVLTree
	//  │       ┌── 6
	//  │   ┌── 5
	//  └── 4
	//      │   ┌── 3
	//      └── 2
	//          └── 1

	if actualValue := t1.Len(); actualValue != 6 {
		t.Errorf("Got %v expected %v", actualValue, 6)
	}
	fmt.Println(t1.String())

	t1.Insert(7, "h")
	t1.Insert(9, "i")
	t1.Insert(8, "j")
	fmt.Println(t1.String())
}

func TestAVLTreeInsert(t *testing.T) {
	t1 := avltree.New[int, string]()
	t1.Insert(5, "e")
	t1.Insert(6, "f")
	t1.Insert(7, "g")
	t1.Insert(3, "c")
	t1.Insert(4, "d")
	t1.Insert(1, "x")
	t1.Insert(2, "b")
	t1.Insert(1, "a") //overwrite

	if actualValue := t1.Len(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}
	if actualValue, expectedValue := t1.Keys(), []int{1, 2, 3, 4, 5, 6, 7}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := t1.Values(), []string{"a", "b", "c", "d", "e", "f", "g"}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tests1 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "", false},
	}

	for _, test := range tests1 {
		// retrievals
		actualValue, actualFound := t1.Get(test[0].(int))
		if actualValue != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}
}

func TestAVLTreeRemove(t *testing.T) {
	t1 := avltree.New[int, string]()
	t1.Insert(5, "e")
	t1.Insert(6, "f")
	t1.Insert(7, "g")
	t1.Insert(3, "c")
	t1.Insert(4, "d")
	t1.Insert(1, "x")
	t1.Insert(2, "b")
	t1.Insert(1, "a") //overwrite

	fmt.Println(t1.String())
	t1.Remove(5)
	fmt.Println(t1.String())
	t1.Remove(6)
	fmt.Println(t1.String())
	t1.Remove(7)
	fmt.Println(t1.String())
	t1.Remove(8)
	fmt.Println(t1.String())
	t1.Remove(5)
	fmt.Println(t1.String())

	if actualValue, expectedValue := t1.Keys(), []int{1, 2, 3, 4}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := t1.Values(), []string{"a", "b", "c", "d"}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := t1.Len(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 4)
	}

	tests2 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "", false},
		{6, "", false},
		{7, "", false},
		{8, "", false},
	}

	for _, test := range tests2 {
		actualValue, actualFound := t1.Get(test[0].(int))
		if actualValue != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}

	t1.Remove(1)
	t1.Remove(4)
	t1.Remove(2)
	t1.Remove(3)
	t1.Remove(2)
	t1.Remove(2)

	if actualValue, expectedValue := t1.Keys(), []int{}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := t1.Values(), []string{}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if empty, size := t1.Len() == 0, t1.Len(); empty != true || size != -0 {
		t.Errorf("Got %v expected %v", empty, true)
	}
}

func TestAVLTreeRemove2(t *testing.T) {
	t1 := avltree.New[int, struct{}]()

	nums := []int{10, 8, 88, 888, 4, 1<<63 - 1, -(1 << 62), 188, -188, 4, 88, 1 << 32}
	for _, v := range nums {
		t1.Insert(v, struct{}{})
	}

	fmt.Println(t1.String())

	t1.Remove(188)

	if ret, _ := t1.InOrder(); !sort.IntsAreSorted(ret) {
		t.Errorf("Error with Remove: %v", ret)
	}

	t1.Remove(188)
	if ret, _ := t1.InOrder(); !sort.IntsAreSorted(ret) {
		t.Errorf("Error with Remove: %v", ret)
	}

	t1.Remove(1<<63 - 1)
	if ret, _ := t1.InOrder(); !sort.IntsAreSorted(ret) {
		t.Errorf("Error with Remove: %v", ret)
	}

	t1.Remove(4)
	if ret, _ := t1.InOrder(); !sort.IntsAreSorted(ret) {
		t.Errorf("Error with Remove: %v", ret)
	}

	if ret := t1.Max(); ret == nil || ret.Key() != (1<<32) {
		t.Errorf("Error with Remove, max: %v, want: %v", ret.Key(), (1 << 32))
	}

	if ret := t1.Min(); ret == nil || ret.Key() != -(1<<62) {
		t.Errorf("Error with Remove, min: %v, want: %v", ret.Key(), (1 << 32))
	}
}

func TestAVLTreeMinAndMax(t *testing.T) {
	t1 := avltree.New[int, string]()

	if actualValue := t1.Min(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := t1.Max(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	fmt.Println(t1.String())

	t1.Insert(1, "a")
	t1.Insert(5, "e")
	t1.Insert(6, "f")
	t1.Insert(7, "g")
	t1.Insert(3, "c")
	t1.Insert(4, "d")
	t1.Insert(1, "x") // overwrite
	t1.Insert(2, "b")

	fmt.Println(t1.String())

	if actualValue, expectedValue := t1.Min().Key(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := t1.Min().Value, "x"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := t1.Max().Key(), 7; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := t1.Max().Value, "g"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestAVLTreeCeilingAndFloor(t *testing.T) {
	t1 := avltree.New[int, string]()

	if node := t1.Floor(0); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
	if node := t1.Ceiling(0); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	t1.Insert(5, "e")
	t1.Insert(6, "f")
	t1.Insert(7, "g")
	t1.Insert(3, "c")
	t1.Insert(4, "d")
	t1.Insert(1, "x")
	t1.Insert(2, "b")

	fmt.Println(t1.String())

	if node := t1.Floor(4); node == nil || node.Key() != 4 {
		t.Errorf("Got %v expected %v", node.Key(), 4)
	}
	if node := t1.Floor(0); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	if node := t1.Ceiling(4); node == nil || node.Key() != 4 {
		t.Errorf("Got %v expected %v", node.Key(), 4)
	}
	if node := t1.Ceiling(8); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
}

func TestAVLTreeString(t *testing.T) {
	t1 := avltree.New[int, struct{}]()
	t1.Insert(1, struct{}{})
	t1.Insert(2, struct{}{})
	t1.Insert(7, struct{}{})
	t1.Insert(3, struct{}{})
	t1.Insert(5, struct{}{})
	t1.Insert(6, struct{}{})
	t1.Insert(4, struct{}{})
	t1.Insert(8, struct{}{})

	if !strings.HasPrefix(t1.String(), "AVLTree") {
		t.Errorf("String should start with container name")
	}

	fmt.Println(t1.String())
}

func TestTraversal(t *testing.T) {
	t1 := avltree.New[int, struct{}]()
	for range 20 {
		t1.Insert(rand.IntN(100), struct{}{})
	}
	fmt.Println(t1.String())

	v1 := t1.Keys()
	fmt.Println(v1)

	v2, _ := t1.LevelOrder()
	fmt.Println(v2)
	assert.Equal(t, len(v1), len(v2))

	v3, _ := t1.PreOrder()
	fmt.Println(v3)
	assert.Equal(t, len(v1), len(v3))

	v5, _ := t1.InOrder()
	fmt.Println(v5)
	assert.Equal(t, v5, v1)

	v7, _ := t1.PostOrder()
	fmt.Println(v7)
	assert.Equal(t, len(v1), len(v7))

	t2 := avltree.New[int, string]()
	t2.Insert(1, "a")
	t2.Insert(5, "e")
	t2.Insert(6, "f")
	t2.Insert(7, "g")
	t2.Insert(3, "c")
	t2.Insert(4, "d")
	t2.Insert(1, "x") // overwrite
	t2.Insert(2, "b")
	fmt.Println(t2.String())
	keys, values := t2.LevelOrder()
	fmt.Println("LevelOrder:")
	fmt.Println(keys)
	fmt.Println(values)
	fmt.Println("---")
	keys, values = t2.PreOrder()
	fmt.Println("PreOrder:")
	fmt.Println(keys)
	fmt.Println(values)
	fmt.Println("---")
	keys, values = t2.InOrder()
	fmt.Println("InOrder:")
	fmt.Println(keys)
	fmt.Println(values)
	fmt.Println("---")
	keys, values = t2.PostOrder()
	fmt.Println("PostOrder:")
	fmt.Println(keys)
	fmt.Println(values)
}

func TestClear(t *testing.T) {
	t1 := avltree.New[int, string]()
	t1.Insert(1, "a")
	t1.Insert(5, "e")
	t1.Insert(6, "f")
	t1.Insert(7, "g")
	t1.Insert(3, "c")
	t1.Insert(4, "d")
	t1.Insert(1, "x") // overwrite
	t1.Insert(2, "b")
	fmt.Println(t1)
	t1.Clear()
	fmt.Println(t1)
}

func benchmarkSearch(b *testing.B, t1 *avltree.Tree[int, struct{}], size int) {
	for b.Loop() {
		for n := range size {
			t1.Search(n)
		}
	}
}

func benchmarkInsert(b *testing.B, t1 *avltree.Tree[int, struct{}], size int) {
	for b.Loop() {
		for n := range size {
			t1.Insert(n, struct{}{})
		}
	}
}

func benchmarkRemove(b *testing.B, t1 *avltree.Tree[int, struct{}], size int) {
	for b.Loop() {
		for n := range size {
			t1.Remove(n)
		}
	}
}

func BenchmarkAVLTreeSearch100(b *testing.B) {
	b.StopTimer()
	size := 100
	t1 := avltree.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkSearch(b, t1, size)
}

func BenchmarkAVLTreeSearch1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	t1 := avltree.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkSearch(b, t1, size)
}

func BenchmarkAVLTreeSearch10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	t1 := avltree.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkSearch(b, t1, size)
}

func BenchmarkAVLTreeSearch100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	t1 := avltree.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkSearch(b, t1, size)
}

func BenchmarkAVLTreeInsert100(b *testing.B) {
	b.StopTimer()
	size := 100
	t1 := avltree.New[int, struct{}]()
	b.StartTimer()
	benchmarkInsert(b, t1, size)
}

func BenchmarkAVLTreeInsert1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	t1 := avltree.New[int, struct{}]()
	b.StartTimer()
	benchmarkInsert(b, t1, size)
}

func BenchmarkAVLTreeInsert10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	t1 := avltree.New[int, struct{}]()
	b.StartTimer()
	benchmarkInsert(b, t1, size)
}

func BenchmarkAVLTreeInsert100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	t1 := avltree.New[int, struct{}]()
	b.StartTimer()
	benchmarkInsert(b, t1, size)
}

func BenchmarkAVLTreeRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	t1 := avltree.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, t1, size)
}

func BenchmarkAVLTreeRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	t1 := avltree.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, t1, size)
}

func BenchmarkAVLTreeRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	t1 := avltree.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, t1, size)
}

func BenchmarkAVLTreeRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	t1 := avltree.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, t1, size)
}
