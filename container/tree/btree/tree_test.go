package btree_test

import (
	"cmp"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/tree/btree"
	"github.com/stretchr/testify/assert"
)

func TestBTree(t *testing.T) {
	testcases := []int{10, 20, 100, 200, 1000, 10000}
	for _, n := range testcases {
		t1 := btree.New[int, struct{}](5)
		nums := rand.Perm(n)
		for _, v := range nums {
			t1.Insert(v, struct{}{})
		}
		fmt.Println("height:", t1.Height())

		k1 := t1.Keys()
		var k2 []int
		t1.Range(func(k int, v struct{}) bool {
			k2 = append(k2, k)
			return true
		})
		assert.Equal(t, k1, k2)

		entries := t1.LevelOrder()
		var levelK1 []int
		for _, e := range entries {
			levelK1 = append(levelK1, e.Key())
		}

		bytes, _ := t1.MarshalJSON()
		t2 := btree.New[int, struct{}](5)
		_ = t2.UnmarshalJSON(bytes)
		k3 := t2.Keys()
		assert.Equal(t, k1, k3)

		var count int
		var k4 []int
		t2.Range(func(k int, v struct{}) bool {
			count++
			k4 = append(k4, k)
			return true
		})
		assert.Equal(t, k1, k4)
		fmt.Println("count:", count)
		assert.Equal(t, count, len(nums))
		assert.Equal(t, count, len(entries))
		assert.Equal(t, count, t2.Len())
		assert.Equal(t, count, t2.Root().Len())

		entries = t2.LevelOrder()
		var levelK2 []int
		for _, e := range entries {
			levelK2 = append(levelK2, e.Key())
		}

		if t2.Len() <= 20 {
			fmt.Println("k1:", k1)
			fmt.Println("k2:", k2)
			fmt.Println("k3:", k3)
			fmt.Println("k4:", k4)
			fmt.Println("level-k1:", levelK1)
			fmt.Println("level-k2:", levelK2)
			fmt.Println("-----")
			fmt.Println(t1)
			fmt.Println("-----")
			fmt.Println(t2)
			fmt.Println("-----")
		}

		entries = t2.InOrder()
		rets := make([]int, 0, len(entries))
		for _, e := range entries {
			rets = append(rets, e.Key())
		}
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
			entries := t2.InOrder()
			rets := make([]int, 0, len(entries))
			for _, e := range entries {
				rets = append(rets, e.Key())
			}
			if !sort.IntsAreSorted(rets) {
				t.Errorf("Error With Remove")
			}
		}
	}
}

func TestBTreeGet1(t *testing.T) {
	t1 := btree.New[int, string](3)
	t1.Insert(1, "a")
	t1.Insert(2, "b")
	t1.Insert(3, "c")
	t1.Insert(4, "d")
	t1.Insert(5, "e")
	t1.Insert(6, "f")
	t1.Insert(7, "g")

	fmt.Println(t1.String())

	tests := [][]any{
		{0, "", false},
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "", false},
	}

	for _, test := range tests {
		if value, found := t1.Get(test[0].(int)); value != test[1] || found != test[2] {
			t.Errorf("Got %v,%v expected %v,%v", value, found, test[1], test[2])
		}
	}
}

func TestBTreeGet2(t *testing.T) {
	t1 := btree.New[int, string](3)
	t1.Insert(7, "g")
	t1.Insert(9, "i")
	t1.Insert(10, "j")
	t1.Insert(6, "f")
	t1.Insert(3, "c")
	t1.Insert(4, "d")
	t1.Insert(5, "e")
	t1.Insert(8, "h")
	t1.Insert(2, "b")
	t1.Insert(1, "a")

	fmt.Println(t1)

	tests := [][]any{
		{0, "", false},
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "h", true},
		{9, "i", true},
		{10, "j", true},
		{11, "", false},
	}

	for _, test := range tests {
		if value, found := t1.Get(test[0].(int)); value != test[1] || found != test[2] {
			t.Errorf("Got %v,%v expected %v,%v", value, found, test[1], test[2])
		}
	}
}

func TestBTreeGet3(t *testing.T) {
	t1 := btree.New[int, string](3)

	if actualValue := t1.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	if actualValue, _ := t1.Search(2); actualValue != nil {
		t.Errorf("Got %v expected nil", actualValue)
	}

	t1.Insert(1, "x") // 1->x
	t1.Insert(2, "b") // 1->x, 2->b (in order)
	t1.Insert(1, "a") // 1->a, 2->b (in order, replacement)
	t1.Insert(3, "c") // 1->a, 2->b, 3->c (in order)
	t1.Insert(4, "d") // 1->a, 2->b, 3->c, 4->d (in order)
	t1.Insert(5, "e") // 1->a, 2->b, 3->c, 4->d, 5->e (in order)
	t1.Insert(6, "f") // 1->a, 2->b, 3->c, 4->d, 5->e, 6->f (in order)
	t1.Insert(7, "g") // 1->a, 2->b, 3->c, 4->d, 5->e, 6->f, 7->g (in order)

	// BTree
	//         1
	//     2
	//         3
	// 4
	//         5
	//     6
	//         7

	if actualValue := t1.Len(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}

	if actualValue, _ := t1.Search(2); actualValue == nil || actualValue.Len() != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}

	if actualValue, _ := t1.Search(4); actualValue == nil || actualValue.Len() != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}

	if actualValue, _ := t1.Search(8); actualValue != nil {
		t.Errorf("Got %v expected nil", actualValue)
	}
}

func TestBTreeInsert1(t *testing.T) {
	// https://upload.wikimedia.org/wikipedia/commons/3/33/B_tree_insertion_example.png
	t1 := btree.New[int, int](3)
	assertValidTree(t, t1, 0)

	t1.Insert(1, 0)
	assertValidTree(t, t1, 1)
	assertValidTreeNode(t, t1.Root(), 1, 0, []int{1}, false)

	t1.Insert(2, 1)
	assertValidTree(t, t1, 2)
	assertValidTreeNode(t, t1.Root(), 2, 0, []int{1, 2}, false)

	t1.Insert(3, 2)
	assertValidTree(t, t1, 3)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{2}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{3}, true)

	t1.Insert(4, 2)
	assertValidTree(t, t1, 4)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{2}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 2, 0, []int{3, 4}, true)

	t1.Insert(5, 2)
	assertValidTree(t, t1, 5)
	assertValidTreeNode(t, t1.Root(), 2, 3, []int{2, 4}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, t1.Root().Children()[2], 1, 0, []int{5}, true)

	t1.Insert(6, 2)
	assertValidTree(t, t1, 6)
	assertValidTreeNode(t, t1.Root(), 2, 3, []int{2, 4}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, t1.Root().Children()[2], 2, 0, []int{5, 6}, true)

	t1.Insert(7, 2)
	assertValidTree(t, t1, 7)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{4}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 2, []int{6}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[1], 1, 0, []int{7}, true)
}

func TestBTreeInsert2(t *testing.T) {
	t1 := btree.New[int, int](4)
	assertValidTree(t, t1, 0)

	t1.Insert(0, 0)
	assertValidTree(t, t1, 1)
	assertValidTreeNode(t, t1.Root(), 1, 0, []int{0}, false)

	t1.Insert(2, 2)
	assertValidTree(t, t1, 2)
	assertValidTreeNode(t, t1.Root(), 2, 0, []int{0, 2}, false)

	t1.Insert(1, 1)
	assertValidTree(t, t1, 3)
	assertValidTreeNode(t, t1.Root(), 3, 0, []int{0, 1, 2}, false)

	t1.Insert(1, 1)
	assertValidTree(t, t1, 3)
	assertValidTreeNode(t, t1.Root(), 3, 0, []int{0, 1, 2}, false)

	t1.Insert(3, 3)
	assertValidTree(t, t1, 4)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{1}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 2, 0, []int{2, 3}, true)

	t1.Insert(4, 4)
	assertValidTree(t, t1, 5)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{1}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 3, 0, []int{2, 3, 4}, true)

	t1.Insert(5, 5)
	assertValidTree(t, t1, 6)
	assertValidTreeNode(t, t1.Root(), 2, 3, []int{1, 3}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, t1.Root().Children()[2], 2, 0, []int{4, 5}, true)
}

func TestBTreeInsert3(t *testing.T) {
	// http://www.geeksforgeeks.org/b-t1-set-1-insert-2/
	t1 := btree.New[int, int](6)
	assertValidTree(t, t1, 0)

	t1.Insert(10, 0)
	assertValidTree(t, t1, 1)
	assertValidTreeNode(t, t1.Root(), 1, 0, []int{10}, false)

	t1.Insert(20, 1)
	assertValidTree(t, t1, 2)
	assertValidTreeNode(t, t1.Root(), 2, 0, []int{10, 20}, false)

	t1.Insert(30, 2)
	assertValidTree(t, t1, 3)
	assertValidTreeNode(t, t1.Root(), 3, 0, []int{10, 20, 30}, false)

	t1.Insert(40, 3)
	assertValidTree(t, t1, 4)
	assertValidTreeNode(t, t1.Root(), 4, 0, []int{10, 20, 30, 40}, false)

	t1.Insert(50, 4)
	assertValidTree(t, t1, 5)
	assertValidTreeNode(t, t1.Root(), 5, 0, []int{10, 20, 30, 40, 50}, false)

	t1.Insert(60, 5)
	assertValidTree(t, t1, 6)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{30}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 3, 0, []int{40, 50, 60}, true)

	t1.Insert(70, 6)
	assertValidTree(t, t1, 7)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{30}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 4, 0, []int{40, 50, 60, 70}, true)

	t1.Insert(80, 7)
	assertValidTree(t, t1, 8)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{30}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 5, 0, []int{40, 50, 60, 70, 80}, true)

	t1.Insert(90, 8)
	assertValidTree(t, t1, 9)
	assertValidTreeNode(t, t1.Root(), 2, 3, []int{30, 60}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 2, 0, []int{40, 50}, true)
	assertValidTreeNode(t, t1.Root().Children()[2], 3, 0, []int{70, 80, 90}, true)
}

func TestBTreeInsert4(t *testing.T) {
	t1 := btree.New[int, *struct{}](3)
	assertValidTree(t, t1, 0)

	t1.Insert(6, nil)
	assertValidTree(t, t1, 1)
	assertValidTreeNode(t, t1.Root(), 1, 0, []int{6}, false)

	t1.Insert(5, nil)
	assertValidTree(t, t1, 2)
	assertValidTreeNode(t, t1.Root(), 2, 0, []int{5, 6}, false)

	t1.Insert(4, nil)
	assertValidTree(t, t1, 3)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{5}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{6}, true)

	t1.Insert(3, nil)
	assertValidTree(t, t1, 4)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{5}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 2, 0, []int{3, 4}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{6}, true)

	t1.Insert(2, nil)
	assertValidTree(t, t1, 5)
	assertValidTreeNode(t, t1.Root(), 2, 3, []int{3, 5}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{2}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{4}, true)
	assertValidTreeNode(t, t1.Root().Children()[2], 1, 0, []int{6}, true)

	t1.Insert(1, nil)
	assertValidTree(t, t1, 6)
	assertValidTreeNode(t, t1.Root(), 2, 3, []int{3, 5}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 2, 0, []int{1, 2}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{4}, true)
	assertValidTreeNode(t, t1.Root().Children()[2], 1, 0, []int{6}, true)

	t1.Insert(0, nil)
	assertValidTree(t, t1, 7)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{3}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 2, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[1], 1, 0, []int{6}, true)

	t1.Insert(-1, nil)
	assertValidTree(t, t1, 8)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{3}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 2, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[0], 2, 0, []int{-1, 0}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[1], 1, 0, []int{6}, true)

	t1.Insert(-2, nil)
	assertValidTree(t, t1, 9)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{3}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 2, 3, []int{-1, 1}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[0], 1, 0, []int{-2}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[1], 1, 0, []int{0}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[2], 1, 0, []int{2}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[1], 1, 0, []int{6}, true)

	t1.Insert(-3, nil)
	assertValidTree(t, t1, 10)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{3}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 2, 3, []int{-1, 1}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[0], 2, 0, []int{-3, -2}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[1], 1, 0, []int{0}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[2], 1, 0, []int{2}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[1], 1, 0, []int{6}, true)

	t1.Insert(-4, nil)
	assertValidTree(t, t1, 11)
	assertValidTreeNode(t, t1.Root(), 2, 3, []int{-1, 3}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 2, []int{-3}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 2, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[2], 1, 2, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[0], 1, 0, []int{-4}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[1], 1, 0, []int{-2}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, t1.Root().Children()[2].Children()[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, t1.Root().Children()[2].Children()[1], 1, 0, []int{6}, true)
}

func TestBTreeRemove1(t *testing.T) {
	// empty
	t1 := btree.New[int, int](3)
	t1.Remove(1)
	assertValidTree(t, t1, 0)
}

func TestBTreeRemove2(t *testing.T) {
	// leaf node (no underflow)
	t1 := btree.New[int, *struct{}](3)
	t1.Insert(1, nil)
	t1.Insert(2, nil)

	t1.Remove(1)
	assertValidTree(t, t1, 1)
	assertValidTreeNode(t, t1.Root(), 1, 0, []int{2}, false)

	t1.Remove(2)
	assertValidTree(t, t1, 0)
}

func TestBTreeRemove3(t *testing.T) {
	// merge with right (underflow)
	{
		t1 := btree.New[int, *struct{}](3)
		t1.Insert(1, nil)
		t1.Insert(2, nil)
		t1.Insert(3, nil)

		t1.Remove(1)
		assertValidTree(t, t1, 2)
		assertValidTreeNode(t, t1.Root(), 2, 0, []int{2, 3}, false)
	}
	// merge with left (underflow)
	{
		t1 := btree.New[int, *struct{}](3)
		t1.Insert(1, nil)
		t1.Insert(2, nil)
		t1.Insert(3, nil)

		t1.Remove(3)
		assertValidTree(t, t1, 2)
		assertValidTreeNode(t, t1.Root(), 2, 0, []int{1, 2}, false)
	}
}

func TestBTreeRemove4(t *testing.T) {
	// rotate left (underflow)
	t1 := btree.New[int, *struct{}](3)
	t1.Insert(1, nil)
	t1.Insert(2, nil)
	t1.Insert(3, nil)
	t1.Insert(4, nil)

	assertValidTree(t, t1, 4)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{2}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 2, 0, []int{3, 4}, true)

	t1.Remove(1)
	assertValidTree(t, t1, 3)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{3}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{2}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{4}, true)
}

func TestBTreeRemove5(t *testing.T) {
	// rotate right (underflow)
	t1 := btree.New[int, *struct{}](3)
	t1.Insert(1, nil)
	t1.Insert(2, nil)
	t1.Insert(3, nil)
	t1.Insert(0, nil)

	assertValidTree(t, t1, 4)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{2}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 2, 0, []int{0, 1}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{3}, true)

	t1.Remove(3)
	assertValidTree(t, t1, 3)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{1}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{2}, true)
}

func TestBTreeRemove6(t *testing.T) {
	// root height reduction after a series of underflows on right side
	// use simulator: https://www.cs.usfca.edu/~galles/visualization/BTree.html
	t1 := btree.New[int, *struct{}](3)
	t1.Insert(1, nil)
	t1.Insert(2, nil)
	t1.Insert(3, nil)
	t1.Insert(4, nil)
	t1.Insert(5, nil)
	t1.Insert(6, nil)
	t1.Insert(7, nil)

	assertValidTree(t, t1, 7)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{4}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 2, []int{6}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[1], 1, 0, []int{7}, true)

	t1.Remove(7)
	assertValidTree(t, t1, 6)
	assertValidTreeNode(t, t1.Root(), 2, 3, []int{2, 4}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, t1.Root().Children()[2], 2, 0, []int{5, 6}, true)
}

func TestBTreeRemove7(t *testing.T) {
	// root height reduction after a series of underflows on left side
	// use simulator: https://www.cs.usfca.edu/~galles/visualization/BTree.html
	t1 := btree.New[int, *struct{}](3)
	t1.Insert(1, nil)
	t1.Insert(2, nil)
	t1.Insert(3, nil)
	t1.Insert(4, nil)
	t1.Insert(5, nil)
	t1.Insert(6, nil)
	t1.Insert(7, nil)

	assertValidTree(t, t1, 7)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{4}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 2, []int{6}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[1], 1, 0, []int{7}, true)

	t1.Remove(1) // series of underflows
	assertValidTree(t, t1, 6)
	assertValidTreeNode(t, t1.Root(), 2, 3, []int{4, 6}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 2, 0, []int{2, 3}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[2], 1, 0, []int{7}, true)

	// clear all remaining
	t1.Remove(2)
	assertValidTree(t, t1, 5)
	assertValidTreeNode(t, t1.Root(), 2, 3, []int{4, 6}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{3}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[2], 1, 0, []int{7}, true)

	t1.Remove(3)
	assertValidTree(t, t1, 4)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{6}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 2, 0, []int{4, 5}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{7}, true)

	t1.Remove(4)
	assertValidTree(t, t1, 3)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{6}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 0, []int{7}, true)

	t1.Remove(5)
	assertValidTree(t, t1, 2)
	assertValidTreeNode(t, t1.Root(), 2, 0, []int{6, 7}, false)

	t1.Remove(6)
	assertValidTree(t, t1, 1)
	assertValidTreeNode(t, t1.Root(), 1, 0, []int{7}, false)

	t1.Remove(7)
	assertValidTree(t, t1, 0)
}

func TestBTreeRemove8(t *testing.T) {
	// use simulator: https://www.cs.usfca.edu/~galles/visualization/BTree.html
	t1 := btree.New[int, *struct{}](3)
	t1.Insert(1, nil)
	t1.Insert(2, nil)
	t1.Insert(3, nil)
	t1.Insert(4, nil)
	t1.Insert(5, nil)
	t1.Insert(6, nil)
	t1.Insert(7, nil)
	t1.Insert(8, nil)
	t1.Insert(9, nil)

	assertValidTree(t, t1, 9)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{4}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 2, 3, []int{6, 8}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[1], 1, 0, []int{7}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[2], 1, 0, []int{9}, true)

	t1.Remove(1)
	assertValidTree(t, t1, 8)
	assertValidTreeNode(t, t1.Root(), 1, 2, []int{6}, false)
	assertValidTreeNode(t, t1.Root().Children()[0], 1, 2, []int{4}, true)
	assertValidTreeNode(t, t1.Root().Children()[1], 1, 2, []int{8}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[0], 2, 0, []int{2, 3}, true)
	assertValidTreeNode(t, t1.Root().Children()[0].Children()[1], 1, 0, []int{5}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[0], 1, 0, []int{7}, true)
	assertValidTreeNode(t, t1.Root().Children()[1].Children()[1], 1, 0, []int{9}, true)
}

func TestBTreeRemove9(t *testing.T) {
	const max = 1000
	orders := []int{3, 4, 5, 6, 7, 8, 9, 10, 20, 100, 500, 1000, 5000, 10000}
	for _, order := range orders {

		t1 := btree.New[int, int](order)

		{
			for i := 1; i <= max; i++ {
				t1.Insert(i, i)
			}
			assertValidTree(t, t1, max)

			for i := 1; i <= max; i++ {
				if _, found := t1.Get(i); !found {
					t.Errorf("Not found %v", i)
				}
			}

			for i := 1; i <= max; i++ {
				t1.Remove(i)
			}
			assertValidTree(t, t1, 0)
		}

		{
			for i := max; i > 0; i-- {
				t1.Insert(i, i)
			}
			assertValidTree(t, t1, max)

			for i := max; i > 0; i-- {
				if _, found := t1.Get(i); !found {
					t.Errorf("Not found %v", i)
				}
			}

			for i := max; i > 0; i-- {
				t1.Remove(i)
			}
			assertValidTree(t, t1, 0)
		}
	}
}

func TestBTreeHeight(t *testing.T) {
	t1 := btree.New[int, int](3)
	if actualValue, expectedValue := t1.Height(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	t1.Insert(1, 0)
	if actualValue, expectedValue := t1.Height(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	t1.Insert(2, 1)
	if actualValue, expectedValue := t1.Height(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	t1.Insert(3, 2)
	if actualValue, expectedValue := t1.Height(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	t1.Insert(4, 2)
	if actualValue, expectedValue := t1.Height(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	t1.Insert(5, 2)
	if actualValue, expectedValue := t1.Height(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	t1.Insert(6, 2)
	if actualValue, expectedValue := t1.Height(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	t1.Insert(7, 2)
	if actualValue, expectedValue := t1.Height(), 3; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	t1.Remove(1)
	t1.Remove(2)
	t1.Remove(3)
	t1.Remove(4)
	t1.Remove(5)
	t1.Remove(6)
	t1.Remove(7)
	if actualValue, expectedValue := t1.Height(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeLeftAndRight(t *testing.T) {
	t1 := btree.New[int, string](3)

	if actualValue := t1.Min(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue.Key(), nil)
	}
	if actualValue := t1.Max(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue.Key(), nil)
	}

	t1.Insert(1, "a")
	t1.Insert(5, "e")
	t1.Insert(6, "f")
	t1.Insert(7, "g")
	t1.Insert(3, "c")
	t1.Insert(4, "d")
	t1.Insert(1, "x") // overwrite
	t1.Insert(2, "b")

	if actualValue, expectedValue := t1.Min(), 1; actualValue == nil || actualValue.Key() != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := t1.Min(), "x"; actualValue == nil || actualValue.Value != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := t1.Max(), 7; actualValue == nil || actualValue.Key() != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := t1.Max(), "g"; actualValue == nil || actualValue.Value != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIteratorValuesAndKeys(t *testing.T) {
	t1 := btree.New[int, string](4)
	t1.Insert(4, "d")
	t1.Insert(5, "e")
	t1.Insert(6, "f")
	t1.Insert(3, "c")
	t1.Insert(1, "a")
	t1.Insert(7, "g")
	t1.Insert(2, "b")
	t1.Insert(1, "x") // override
	fmt.Println(t1)
	fmt.Println(t1.Keys())
	fmt.Println(t1.Values())
	if actualValue, expectedValue := t1.Keys(), []int{1, 2, 3, 4, 5, 6, 7}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := t1.Values(), []string{"x", "b", "c", "d", "e", "f", "g"}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := t1.Len(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}
}

func assertValidTree[K comparable, V any](t *testing.T, t1 *btree.Tree[K, V], expectedSize int) {
	if actualValue, expectedValue := t1.Len(), expectedSize; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for tree size", actualValue, expectedValue)
	}
}

func assertValidTreeNode[K comparable, V any](
	t *testing.T,
	x *btree.Node[K, V],
	expectedEntries int,
	expectedChildren int,
	keys []K,
	hasParent bool) {
	if actualValue, expectedValue := x.Parent() != nil, hasParent; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for hasParent", actualValue, expectedValue)
	}
	if actualValue, expectedValue := len(x.Entries), expectedEntries; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for entries size", actualValue, expectedValue)
	}
	if actualValue, expectedValue := len(x.Children()), expectedChildren; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for children size", actualValue, expectedValue)
	}
	for i, key := range keys {
		if actualValue, expectedValue := x.Entries[i].Key(), key; actualValue != expectedValue {
			t.Errorf("Got %v expected %v for key", actualValue, expectedValue)
		}
	}
}

type Entry[K comparable, V any] struct {
	// The key used to compare entries.
	Key K

	// The value stored with this entry.
	Value V
}

func searchEntries[K cmp.Ordered, V any](entries []*Entry[K, V], k K) (index int, ok bool) {
	i, j := 0, len(entries)-1
	for i <= j {
		mid := (j + i) / 2
		val := cmp.Compare(k, entries[mid].Key)
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

func TestSearchEntries(t *testing.T) {
	{
		entries := []*Entry[int, int]{}
		tests := [][]any{
			{0, 0, false},
		}
		for _, test := range tests {
			index, found := searchEntries(entries, test[0].(int))
			if actualValue, expectedValue := index, test[1]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
			if actualValue, expectedValue := found, test[2]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
	{
		entries := []*Entry[int, int]{{2, 0}, {4, 1}, {6, 2}}
		tests := [][]any{
			{0, 0, false},
			{1, 0, false},
			{2, 0, true},
			{3, 1, false},
			{4, 1, true},
			{5, 2, false},
			{6, 2, true},
			{7, 3, false},
		}
		for _, test := range tests {
			index, found := searchEntries(entries, test[0].(int))
			if actualValue, expectedValue := index, test[1]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
			if actualValue, expectedValue := found, test[2]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
}

func TestBTreeSearch(t *testing.T) {
	{
		t1 := btree.New[int, int](3)
		tests := [][]any{
			{0, false, -1},
		}
		for _, test := range tests {
			node, index := t1.Search(test[0].(int))
			if actualValue, expectedValue := node != nil, test[1]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
			if actualValue, expectedValue := index, test[2]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
	{
		t1 := btree.New[int, int](3)
		t1.Insert(2, 0)
		t1.Insert(4, 1)
		t1.Insert(6, 2)
		tests := [][]any{
			{0, 0, false},
			{1, 0, false},
			{2, 0, true},
			{3, 0, false},
			{4, 1, true},
			{5, 0, false},
			{6, 2, true},
			{7, 0, false},
		}
		for _, test := range tests {
			value, ok := t1.Get(test[0].(int))
			if actualValue, expectedValue := value, test[1]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
			if actualValue, expectedValue := ok, test[2]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
}

func TestBTreeSerialization(t *testing.T) {
	t1 := btree.New[string, string](3)
	t1.Insert("c", "3")
	t1.Insert("b", "2")
	t1.Insert("a", "1")

	var err error
	assert := func() {
		if actualValue, expectedValue := t1.Len(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if actualValue, expectedValue := t1.Keys(), []string{"a", "b", "c"}; !slices.Equal(actualValue, expectedValue) {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if actualValue, expectedValue := t1.Values(), []string{"1", "2", "3"}; !slices.Equal(actualValue, expectedValue) {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	bytes, err := t1.MarshalJSON()
	assert()

	fmt.Println(string(bytes))

	t2 := btree.New[string, string](3)
	fmt.Println("t2:", t2)
	err = t2.UnmarshalJSON(bytes)
	assert()
	fmt.Println("t2:", t2)

	bytes, err = json.Marshal([]any{"a", "b", "c", t2})
	if err != nil {
		t.Errorf("Got error %v", err)
	}
	fmt.Println(string(bytes))

	t3 := btree.New[string, int](3)
	err = json.Unmarshal([]byte(`{"a":1,"b":2}`), t3)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
	if actualValue, expectedValue := t3.Len(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := t3.Keys(), []string{"a", "b"}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := t3.Values(), []int{1, 2}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeString(t *testing.T) {
	t1 := btree.New[string, int](3)
	t1.Insert("a", 1)
	if !strings.HasPrefix(t1.String(), "BTree") {
		t.Errorf("String should start with container name")
	}
}

func TestBTreeClear(t *testing.T) {
	t1 := btree.New[int, string](3)

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

func benchmarkGet(b *testing.B, t1 *btree.Tree[int, struct{}], size int) {
	for b.Loop() {
		for n := range size {
			t1.Get(n)
		}
	}
}

func benchmarkInsert(b *testing.B, t1 *btree.Tree[int, struct{}], size int) {
	for b.Loop() {
		for n := range size {
			t1.Insert(n, struct{}{})
		}
	}
}

func benchmarkRemove(b *testing.B, t1 *btree.Tree[int, struct{}], size int) {
	for b.Loop() {
		for n := range size {
			t1.Remove(n)
		}
	}
}

func BenchmarkBTreeGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	t1 := btree.New[int, struct{}](128)
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkGet(b, t1, size)
}

func BenchmarkBTreeGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	t1 := btree.New[int, struct{}](128)
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkGet(b, t1, size)
}

func BenchmarkBTreeGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	t1 := btree.New[int, struct{}](128)
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkGet(b, t1, size)
}

func BenchmarkBTreeGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	t1 := btree.New[int, struct{}](128)
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkGet(b, t1, size)
}

func BenchmarkBTreeInsert100(b *testing.B) {
	b.StopTimer()
	size := 100
	t1 := btree.New[int, struct{}](128)
	b.StartTimer()
	benchmarkInsert(b, t1, size)
}

func BenchmarkBTreeInsert1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	t1 := btree.New[int, struct{}](128)
	b.StartTimer()
	benchmarkInsert(b, t1, size)
}

func BenchmarkBTreeInsert10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	t1 := btree.New[int, struct{}](128)
	b.StartTimer()
	benchmarkInsert(b, t1, size)
}

func BenchmarkBTreeInsert100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	t1 := btree.New[int, struct{}](128)
	b.StartTimer()
	benchmarkInsert(b, t1, size)
}

func BenchmarkBTreeRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	t1 := btree.New[int, struct{}](128)
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, t1, size)
}

func BenchmarkBTreeRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	t1 := btree.New[int, struct{}](128)
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, t1, size)
}

func BenchmarkBTreeRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	t1 := btree.New[int, struct{}](128)
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, t1, size)
}

func BenchmarkBTreeRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	t1 := btree.New[int, struct{}](128)
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, t1, size)
}
