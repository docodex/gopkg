package skiplist_test

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/skiplist"
	"github.com/stretchr/testify/assert"
)

func TestSkiplist(t *testing.T) {
	l := skiplist.New[string, int]()

	l.Insert("leo", 95)
	x := l.MinNode()
	assert.NotNil(t, x)
	fmt.Println(x.Element.Key(), x.Element.Value)
	x = x.Next()
	assert.Nil(t, x)
	fmt.Println(l)
	fmt.Println("-----------------------------")

	l.Insert("jack", 88)
	x = l.MinNode()
	assert.NotNil(t, x)
	fmt.Println(x.Element.Key(), x.Element.Value)
	x = x.Next()
	assert.NotNil(t, x)
	fmt.Println(x.Element.Key(), x.Element.Value)
	x = x.Next()
	assert.Nil(t, x)
	fmt.Println(l)
	fmt.Println("-----------------------------")

	l.Insert("lily", 100)
	x = l.MinNode()
	assert.NotNil(t, x)
	fmt.Println(x.Element.Key(), x.Element.Value)
	x = x.Next()
	assert.NotNil(t, x)
	fmt.Println(x.Element.Key(), x.Element.Value)
	x = x.Next()
	fmt.Println(x.Element.Key(), x.Element.Value)
	x = x.Next()
	assert.Nil(t, x)
	fmt.Println(l)
	fmt.Println("-----------------------------")

	y, r := l.Get("jack")
	assert.NotNil(t, y)
	fmt.Println(y.Key(), y.Value, r)
	fmt.Println("-----------------------------")

	l.Remove("leo")
	x = l.MinNode()
	assert.NotNil(t, x)
	fmt.Println(x.Element.Key(), x.Element.Value)
	x = x.Next()
	assert.NotNil(t, x)
	fmt.Println(x.Element.Key(), x.Element.Value)
	x = x.Next()
	assert.Nil(t, x)
	fmt.Println(l)
	fmt.Println("-----------------------------")
}

func TestInsert(t *testing.T) {
	l := skiplist.New[int, int]()

	m := make(map[int]int)
	for i := range 100 {
		key := rand.Int() % 100
		l.Insert(key, i)
		m[key] = i
	}
	for key, v := range m {
		e, _ := l.Get(key)
		assert.NotNil(t, e)
		assert.Equal(t, v, e.Value)
	}
	assert.Equal(t, len(m), l.Len())
	fmt.Println(l)
}

func TestRemove(t *testing.T) {
	l := skiplist.New[int, int]()

	m := make(map[int]int)
	for i := range 1000 {
		key := rand.Int() % 1000
		l.Insert(key, i)
		m[key] = i
	}
	assert.Equal(t, len(m), l.Len())

	for range 300 {
		key := rand.Int() % 1000
		l.Remove(key)
		delete(m, key)
		key2 := rand.Int() % 10440
		l.Insert(key2, key)
		m[key2] = key
	}

	for key, v := range m {
		e, _ := l.Get(key)
		assert.NotNil(t, e)
		assert.Equal(t, v, e.Value)
	}
	assert.Equal(t, len(m), l.Len())

	{
		l.Insert(10000, 20000)
		e, _ := l.Get(10000)
		assert.NotNil(t, e)
		assert.Equal(t, 10000, e.Key())
		assert.Equal(t, 20000, e.Value)
	}
	{
		e1 := l.Remove(10000)
		assert.NotNil(t, e1)
		assert.Equal(t, 10000, e1.Key())
		assert.Equal(t, 20000, e1.Value)
	}
	{
		e, _ := l.Get(10000)
		assert.Nil(t, e)
	}
}

func TestSkiplist_Range(t *testing.T) {
	list := skiplist.New[int, int]()
	for i := range 10 {
		list.Insert(i, i*10)
	}
	keys := list.Keys()
	for i := range 10 {
		assert.Equal(t, i, keys[i])
	}
	i := 0
	list.Range(func(k, v int) bool {
		assert.Equal(t, i, k)
		assert.Equal(t, i*10, v)
		i++
		return true
	})
}

func TestSkiplist1(t *testing.T) {
	testcases := []int{100, 200, 1000, 10000}
	for _, n := range testcases {
		l := skiplist.New[int, struct{}]()
		nums := rand.Perm(n)
		for _, v := range nums {
			l.Insert(v, struct{}{})
		}

		fmt.Println(l.Len())
		bytes, _ := l.MarshalJSON()
		_ = l.UnmarshalJSON(bytes)
		fmt.Println(l.Len())
		fmt.Println("-------------")

		var count int
		l.Range(func(k int, v struct{}) bool {
			count++
			return true
		})
		fmt.Println("count:", count)
		assert.Equal(t, count, len(nums))
		assert.Equal(t, count, l.Len())

		rets := l.Keys()
		if !sort.IntsAreSorted(rets) {
			t.Error("Error with Push")
		}

		if res := l.Min(); res == nil || res.Key() != rets[0] {
			t.Errorf("Error with Min, get %d, want: %d", res.Value, rets[0])
		}

		if res := l.Max(); res == nil || res.Key() != rets[n-1] {
			t.Errorf("Error with Max, get %d, want: %d", res.Value, rets[n-1])
		}

		for i := 0; i < n-1; i++ {
			l.Remove(nums[i])
			rets = l.Keys()
			if !sort.IntsAreSorted(rets) {
				t.Errorf("Error With Remove")
			}
		}
	}
}

func TestSkiplistInit(t *testing.T) {
	testcases := []int{100, 200, 1000, 10000}
	for _, n := range testcases {
		nums := rand.Perm(n)
		m := make(map[int]struct{}, len(nums))
		for _, v := range nums {
			m[v] = struct{}{}
		}
		v, err := json.Marshal(m)
		assert.Nil(t, err)
		l := skiplist.New[int, struct{}]()
		err = l.UnmarshalJSON(v)
		assert.Nil(t, err)

		v, err = l.MarshalJSON()
		assert.Nil(t, err)
		l1 := skiplist.New[int, struct{}]()
		err = l1.UnmarshalJSON(v)
		assert.Nil(t, err)

		rets := l.Keys()
		if !sort.IntsAreSorted(rets) {
			t.Error("Error with Push")
		}

		if res := l.Min(); res == nil || res.Key() != rets[0] {
			t.Errorf("Error with Min, get %d, want: %d", res.Value, rets[0])
		}

		if res := l.Max(); res == nil || res.Key() != rets[n-1] {
			t.Errorf("Error with Max, get %d, want: %d", res.Value, rets[n-1])
		}

		for i := 0; i < n-1; i++ {
			l.Remove(nums[i])
			rets = l.Keys()
			if !sort.IntsAreSorted(rets) {
				t.Errorf("Error With Remove")
			}
		}

		rets = l1.Keys()
		if !sort.IntsAreSorted(rets) {
			t.Error("Error with Push")
		}

		if res := l1.Min(); res == nil || res.Key() != rets[0] {
			t.Errorf("Error with Min, get %d, want: %d", res.Value, rets[0])
		}

		if res := l1.Max(); res == nil || res.Key() != rets[n-1] {
			t.Errorf("Error with Max, get %d, want: %d", res.Value, rets[n-1])
		}

		for i := 0; i < n-1; i++ {
			l1.Remove(nums[i])
			rets = l1.Keys()
			if !sort.IntsAreSorted(rets) {
				t.Errorf("Error With Remove")
			}
		}
	}
}

func TestSkiplistGet(t *testing.T) {
	type Entry struct {
		key, value int
	}

	l := skiplist.New[int, int]()
	n := 10
	for i := range n {
		l.Insert(i, i+10000)
	}
	assert.False(t, l.Len() == 0)
	assert.Equal(t, n, l.Len())

	fmt.Println(l)
	buf, _ := l.MarshalJSON()
	fmt.Println(string(buf))

	for i := range n {
		x, rank := l.Get(i)
		assert.Equal(t, i+10000, x.Value)
		assert.Equal(t, rank, i+1)
		prev := l.GetByRank(rank - 1)
		if prev != nil {
			fmt.Printf("%d:%d, %d:%d\n", x.Key(), x.Value, prev.Key(), prev.Value)
			assert.Equal(t, prev.Key(), i-1)
			assert.Equal(t, prev.Value, i-1+10000)
		}
		next := l.GetByRank(rank + 1)
		if next != nil {
			fmt.Printf("%d:%d, %d:%d\n", x.Key(), x.Value, next.Key(), next.Value)
			assert.Equal(t, next.Key(), i+1)
			assert.Equal(t, next.Value, i+1+10000)
		}
	}

	es := l.GetRange(4, 8)
	for i := range 4 {
		assert.NotNil(t, es[i])
		assert.Equal(t, i+4, es[i].Key())
		assert.Equal(t, i+4+10000, es[i].Value)
	}

	es = l.GetRange(4, 14)
	assert.Equal(t, 6, len(es))
	for i := range len(es) {
		assert.NotNil(t, es[i])
		assert.Equal(t, i+4, es[i].Key())
		assert.Equal(t, i+4+10000, es[i].Value)
	}

	es = l.GetRangeByRank(5, 10)
	for i := range 5 {
		assert.NotNil(t, es[i])
		assert.Equal(t, i+4, es[i].Key())
		assert.Equal(t, i+4+10000, es[i].Value)
	}

	es = l.GetRangeByRank(5, 20)
	assert.Equal(t, 6, len(es))
	for i := range len(es) {
		assert.NotNil(t, es[i])
		assert.Equal(t, i+4, es[i].Key())
		assert.Equal(t, i+4+10000, es[i].Value)
	}
}

func TestSkiplistDelete(t *testing.T) {
	l := skiplist.New[int, int]()
	m := make(map[int]int)
	for i := range 1000 {
		l.Insert(i, i+10000)
		m[i] = i + 10000
	}
	count := 1000
	for k, v := range m {
		x, _ := l.Get(k)
		assert.Equal(t, v, x.Value)
		l.Remove(k)
		x, _ = l.Get(k)
		assert.Nil(t, x)
		count--
		assert.Equal(t, count, l.Len())
	}
}

func TestSkiplistInsert1(t *testing.T) {
	l := skiplist.New[int, struct{}]()
	l.Insert(3, struct{}{})
	l.Insert(4, struct{}{})
	l.Insert(5, struct{}{})

	x := l.MinNode()
	if x.Element.Key() != 3 {
		t.Errorf("Min should have key = 3, not %v", x.Element.Key())
	}
	y := x.Next()
	assert.NotNil(t, y)
	if y.Element.Key() != 4 {
		t.Errorf("Rank 2 should have key = 4")
	}
	assert.True(t, y.Prev() == x)
	z := l.MaxNode()
	assert.NotNil(t, z)
	if z.Element.Key() != 5 {
		t.Errorf("Max should have key = 5")
	}
	assert.True(t, z.Prev() == y)
}

func TestSkiplistRemove1(t *testing.T) {
	nums := []int{100, 500, 1000, 10_000}
	for _, n := range nums {
		l := skiplist.New[int, struct{}]()
		nums := rand.Perm(n)
		for _, v := range nums {
			l.Insert(v, struct{}{})
		}

		rets := l.Keys()
		if !sort.IntsAreSorted(rets) {
			t.Error("Error with Push")
		}

		if res := l.Min(); res == nil || res.Key() != rets[0] {
			t.Errorf("Error with Min, get %v, want: %d", res, rets[0])
		}

		if res := l.Max(); res == nil || res.Key() != rets[n-1] {
			t.Errorf("Error with Max, get %v, want: %d", res, rets[n-1])
		}

		for i := range n {
			l.Remove(nums[i])
			rets = l.Keys()
			if !sort.IntsAreSorted(rets) {
				t.Errorf("Error With Remove")
			}
		}
	}
}

func TestSkiplistSearch(t *testing.T) {
	l := skiplist.New[int, string]()
	if actualValue := l.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
	if actualValue, _ := l.Get(2); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	l.Insert(1, "x") // 1->x
	l.Insert(2, "b") // 1->x, 2->b
	l.Insert(1, "a") // 1->a, 2->b ( replacement)
	l.Insert(3, "c") // 1->a, 2->b, 3->c
	l.Insert(4, "d") // 1->a, 2->b, 3->c, 4->d
	l.Insert(5, "e") // 1->a, 2->b, 3->c, 4->d, 5->e
	l.Insert(6, "f") // 1->a, 2->b, 3->c, 4->d, 5->e, 6->f

	if actualValue := l.Len(); actualValue != 6 {
		t.Errorf("Got %v expected %v", actualValue, 6)
	}
	fmt.Println(l.String())

	l.Insert(7, "h")
	l.Insert(9, "i")
	l.Insert(8, "j")
	fmt.Println(l.String())
}

func TestSkiplistInsert(t *testing.T) {
	t1 := skiplist.New[int, string]()
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
		{1, "a", 1},
		{2, "b", 2},
		{3, "c", 3},
		{4, "d", 4},
		{5, "e", 5},
		{6, "f", 6},
		{7, "g", 7},
		{8, "", 0},
	}

	for _, test := range tests1 {
		// retrievals
		e, rank := t1.Get(test[0].(int))
		if e == nil {
			if rank != test[2] {
				t.Errorf("Got %v expected %v", e, test[1])
			}
		} else {
			if e.Value != test[1] || rank != test[2] {
				t.Errorf("Got %v expected %v", e, test[1])
			}
		}
	}
}

func TestSkiplistRemove(t *testing.T) {
	t1 := skiplist.New[int, string]()
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
		{1, "a", 1},
		{2, "b", 2},
		{3, "c", 3},
		{4, "d", 4},
		{5, "", 0},
		{6, "", 0},
		{7, "", 0},
		{8, "", 0},
	}

	for _, test := range tests2 {
		e, rank := t1.Get(test[0].(int))
		if e == nil {
			if rank != test[2] {
				t.Errorf("Got %v expected %v", e, test[1])
			}
		} else {
			if e.Value != test[1] || rank != test[2] {
				t.Errorf("Got %v expected %v", e, test[1])
			}
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

func TestSkiplistRemove2(t *testing.T) {
	t1 := skiplist.New[int, struct{}]()

	nums := []int{10, 8, 88, 888, 4, 1<<63 - 1, -(1 << 62), 188, -188, 4, 88, 1 << 32}
	for _, v := range nums {
		t1.Insert(v, struct{}{})
	}

	fmt.Println(t1.String())

	t1.Remove(188)

	if ret := t1.Keys(); !sort.IntsAreSorted(ret) {
		t.Errorf("Error with Remove: %v", ret)
	}

	t1.Remove(188)
	if ret := t1.Keys(); !sort.IntsAreSorted(ret) {
		t.Errorf("Error with Remove: %v", ret)
	}

	t1.Remove(1<<63 - 1)
	if ret := t1.Keys(); !sort.IntsAreSorted(ret) {
		t.Errorf("Error with Remove: %v", ret)
	}

	t1.Remove(4)
	if ret := t1.Keys(); !sort.IntsAreSorted(ret) {
		t.Errorf("Error with Remove: %v", ret)
	}

	if ret := t1.Max(); ret == nil || ret.Key() != (1<<32) {
		t.Errorf("Error with Remove, max: %v, want: %v", ret.Key(), (1 << 32))
	}

	if ret := t1.Min(); ret == nil || ret.Key() != -(1<<62) {
		t.Errorf("Error with Remove, min: %v, want: %v", ret.Key(), (1 << 32))
	}
}

func TestSkiplistRemove3(t *testing.T) {
	t1 := skiplist.New[int, string]()
	t1.Insert(5, "e")
	t1.Insert(6, "f")
	t1.Insert(7, "g")
	t1.Insert(3, "c")
	t1.Insert(4, "d")
	t1.Insert(1, "x")
	t1.Insert(2, "b")
	t1.Insert(1, "a") //overwrite

	fmt.Println(t1.String())
	t1.RemoveByRank(5)
	fmt.Println(t1.String())
	t1.RemoveByRank(5)
	fmt.Println(t1.String())
	t1.RemoveByRank(5)
	fmt.Println(t1.String())
	t1.RemoveByRank(5)
	fmt.Println(t1.String())
	t1.RemoveByRank(5)
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
		{1, "a", 1},
		{2, "b", 2},
		{3, "c", 3},
		{4, "d", 4},
		{5, "", 0},
		{6, "", 0},
		{7, "", 0},
		{8, "", 0},
	}

	for _, test := range tests2 {
		e, rank := t1.Get(test[0].(int))
		if e == nil {
			if rank != test[2] {
				t.Errorf("Got %v expected %v", e, test[1])
			}
		} else {
			if e.Value != test[1] || rank != test[2] {
				t.Errorf("Got %v expected %v", e, test[1])
			}
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

func TestSkiplistRemove4(t *testing.T) {
	t1 := skiplist.New[int, string]()
	t1.Insert(5, "e")
	t1.Insert(6, "f")
	t1.Insert(7, "g")
	t1.Insert(3, "c")
	t1.Insert(4, "d")
	t1.Insert(1, "x")
	t1.Insert(2, "b")
	t1.Insert(1, "a") //overwrite

	fmt.Println(t1.String())
	t1.RemoveRange(5, 9)
	fmt.Println(t1.String())
	t1.RemoveRange(5, 9)
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
		{1, "a", 1},
		{2, "b", 2},
		{3, "c", 3},
		{4, "d", 4},
		{5, "", 0},
		{6, "", 0},
		{7, "", 0},
		{8, "", 0},
	}

	for _, test := range tests2 {
		e, rank := t1.Get(test[0].(int))
		if e == nil {
			if rank != test[2] {
				t.Errorf("Got %v expected %v", e, test[1])
			}
		} else {
			if e.Value != test[1] || rank != test[2] {
				t.Errorf("Got %v expected %v", e, test[1])
			}
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

func TestSkiplistRemove5(t *testing.T) {
	t1 := skiplist.New[int, string]()
	t1.Insert(5, "e")
	t1.Insert(6, "f")
	t1.Insert(7, "g")
	t1.Insert(3, "c")
	t1.Insert(4, "d")
	t1.Insert(1, "x")
	t1.Insert(2, "b")
	t1.Insert(1, "a") //overwrite

	fmt.Println(t1.String())
	t1.RemoveRangeByRank(5, 9)
	fmt.Println(t1.String())
	t1.RemoveRangeByRank(5, 9)
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
		{1, "a", 1},
		{2, "b", 2},
		{3, "c", 3},
		{4, "d", 4},
		{5, "", 0},
		{6, "", 0},
		{7, "", 0},
		{8, "", 0},
	}

	for _, test := range tests2 {
		e, rank := t1.Get(test[0].(int))
		if e == nil {
			if rank != test[2] {
				t.Errorf("Got %v expected %v", e, test[1])
			}
		} else {
			if e.Value != test[1] || rank != test[2] {
				t.Errorf("Got %v expected %v", e, test[1])
			}
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

func TestSkiplistMinAndMax(t *testing.T) {
	t1 := skiplist.New[int, string]()

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

func TestSkiplistCeilingAndFloor(t *testing.T) {
	l := skiplist.New[int, string]()

	if node := l.MinNodeInRange(0, 100); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
	if node := l.MaxNodeInRange(0, 100); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	l.Insert(5, "e")
	l.Insert(6, "f")
	l.Insert(7, "g")
	l.Insert(3, "c")
	l.Insert(4, "d")
	l.Insert(1, "x")
	l.Insert(2, "b")

	fmt.Println(l.String())

	if node := l.MinInRange(4, 10); node == nil || node.Key() != 4 {
		t.Errorf("Got %v expected %v", node.Key(), 4)
	}
	if node := l.MinInRange(8, 10); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	if node := l.MaxInRange(0, 5); node == nil || node.Key() != 4 {
		t.Errorf("Got %v expected %v", node.Key(), 4)
	}
	if node := l.MaxInRange(0, 1); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
}

func TestSkiplistString(t *testing.T) {
	t1 := skiplist.New[int, struct{}]()
	t1.Insert(1, struct{}{})
	t1.Insert(2, struct{}{})
	t1.Insert(7, struct{}{})
	t1.Insert(3, struct{}{})
	t1.Insert(5, struct{}{})
	t1.Insert(6, struct{}{})
	t1.Insert(4, struct{}{})
	t1.Insert(8, struct{}{})

	if !strings.HasPrefix(t1.String(), "Skiplist") {
		t.Errorf("String should start with container name")
	}

	fmt.Println(t1.String())
}

func TestTraversal(t *testing.T) {
	l := skiplist.New[int, struct{}]()
	for range 20 {
		l.Insert(rand.IntN(100), struct{}{})
	}
	fmt.Println(l.String())

	v1 := l.Keys()
	fmt.Println(v1)

	l1 := skiplist.New[int, string]()
	l1.Insert(1, "a")
	l1.Insert(5, "e")
	l1.Insert(6, "f")
	l1.Insert(7, "g")
	l1.Insert(3, "c")
	l1.Insert(4, "d")
	l1.Insert(1, "x") // overwrite
	l1.Insert(2, "b")
	fmt.Println(l1.String())
	keys := l1.Keys()
	fmt.Println(keys)
	values := l1.Values()
	fmt.Println(values)
}

func TestClear(t *testing.T) {
	t1 := skiplist.New[int, string]()
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

func benchmarkGet(b *testing.B, t1 *skiplist.Skiplist[int, struct{}], size int) {
	for b.Loop() {
		for n := range size {
			t1.Get(n)
		}
	}
}

func benchmarkInsert(b *testing.B, t1 *skiplist.Skiplist[int, struct{}], size int) {
	for b.Loop() {
		for n := range size {
			t1.Insert(n, struct{}{})
		}
	}
}

func benchmarkRemove(b *testing.B, t1 *skiplist.Skiplist[int, struct{}], size int) {
	for b.Loop() {
		for n := range size {
			t1.Remove(n)
		}
	}
}

func BenchmarkSkiplistSearch100(b *testing.B) {
	b.StopTimer()
	size := 100
	t1 := skiplist.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkGet(b, t1, size)
}

func BenchmarkSkiplistSearch1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	t1 := skiplist.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkGet(b, t1, size)
}

func BenchmarkSkiplistSearch10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	t1 := skiplist.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkGet(b, t1, size)
}

func BenchmarkSkiplistSearch100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	t1 := skiplist.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkGet(b, t1, size)
}

func BenchmarkSkiplistInsert100(b *testing.B) {
	b.StopTimer()
	size := 100
	t1 := skiplist.New[int, struct{}]()
	b.StartTimer()
	benchmarkInsert(b, t1, size)
}

func BenchmarkSkiplistInsert1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	t1 := skiplist.New[int, struct{}]()
	b.StartTimer()
	benchmarkInsert(b, t1, size)
}

func BenchmarkSkiplistInsert10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	t1 := skiplist.New[int, struct{}]()
	b.StartTimer()
	benchmarkInsert(b, t1, size)
}

func BenchmarkSkiplistInsert100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	t1 := skiplist.New[int, struct{}]()
	b.StartTimer()
	benchmarkInsert(b, t1, size)
}

func BenchmarkSkiplistRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	t1 := skiplist.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, t1, size)
}

func BenchmarkSkiplistRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	t1 := skiplist.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, t1, size)
}

func BenchmarkSkiplistRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	t1 := skiplist.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, t1, size)
}

func BenchmarkSkiplistRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	t1 := skiplist.New[int, struct{}]()
	for n := range size {
		t1.Insert(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, t1, size)
}
