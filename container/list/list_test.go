package list_test

import (
	"testing"

	"github.com/docodex/gopkg/container/list"
	"github.com/docodex/gopkg/container/list/arraylist"
	"github.com/stretchr/testify/assert"
)

func newList(vs ...int) list.List[int] {
	l := arraylist.New[int]()
	l.PushBack(vs...)
	return l
}

func TestIndex(t *testing.T) {
	// Nil list
	assert.Equal(t, -1, list.Index(nil, 1))
	// Empty list
	assert.Equal(t, -1, list.Index(newList(), 1))
	// Found
	assert.Equal(t, 1, list.Index(newList(10, 20, 30), 20))
	// Not found
	assert.Equal(t, -1, list.Index(newList(10, 20, 30), 99))
}

func TestFind(t *testing.T) {
	idx, v := list.Find(newList(1, 2, 3, 4), func(i, v int) bool { return v == 3 })
	assert.Equal(t, 2, idx)
	assert.Equal(t, 3, v)

	// not found
	idx, _ = list.Find(newList(1, 2, 3), func(i, v int) bool { return false })
	assert.Equal(t, -1, idx)

	// nil list / nil f / empty list
	idx, _ = list.Find(nil, func(i, v int) bool { return true })
	assert.Equal(t, -1, idx)
	idx, _ = list.Find(newList(1), nil)
	assert.Equal(t, -1, idx)
	idx, _ = list.Find(newList(), func(i, v int) bool { return true })
	assert.Equal(t, -1, idx)
}

func TestContainsHelper(t *testing.T) {
	// Nil list
	assert.False(t, list.Contains(nil, 1))
	// No query values: vacuously true
	assert.True(t, list.Contains(newList(1), []int{}...))
	// Empty list with query
	assert.False(t, list.Contains(newList(), 1))
	// All present
	assert.True(t, list.Contains(newList(1, 2, 3), 1, 2, 3))
	// Missing one
	assert.False(t, list.Contains(newList(1, 2, 3), 1, 9))
	// Dedup query
	assert.True(t, list.Contains(newList(1, 2), 1, 1, 2, 2))
}

func TestContainsAnyHelper(t *testing.T) {
	// Nil list
	assert.False(t, list.ContainsAny(nil, 1))
	// Empty query
	assert.False(t, list.ContainsAny(newList(1, 2)))
	// Empty list
	assert.False(t, list.ContainsAny(newList(), 1))
	// Found
	assert.True(t, list.ContainsAny(newList(1, 2, 3), 9, 2))
	// Not found
	assert.False(t, list.ContainsAny(newList(1, 2, 3), 7, 8))
}

func TestAllHelper(t *testing.T) {
	// Nil predicate
	assert.False(t, list.All(newList(1), nil))
	// Nil list (vacuous truth)
	assert.True(t, list.All(nil, func(i, v int) bool { return false }))
	// Empty list (vacuous truth)
	assert.True(t, list.All(newList(), func(i, v int) bool { return false }))
	// All match
	assert.True(t, list.All(newList(2, 4, 6), func(i, v int) bool { return v%2 == 0 }))
	// At least one fails
	assert.False(t, list.All(newList(2, 4, 5), func(i, v int) bool { return v%2 == 0 }))
}

func TestAnyHelper(t *testing.T) {
	// Nil list / empty list / nil predicate
	assert.False(t, list.Any(nil, func(i, v int) bool { return true }))
	assert.False(t, list.Any(newList(), func(i, v int) bool { return true }))
	assert.False(t, list.Any(newList(1), nil))
	// Some match
	assert.True(t, list.Any(newList(1, 2, 3), func(i, v int) bool { return v == 2 }))
	// None match
	assert.False(t, list.Any(newList(1, 2, 3), func(i, v int) bool { return v == 99 }))
}

func TestFilter(t *testing.T) {
	// nil inputs: no-op
	dst := arraylist.New[int]()
	list.Filter(dst, nil, func(i, v int) bool { return true })
	assert.Equal(t, 0, dst.Len())
	list.Filter(nil, newList(1), func(i, v int) bool { return true })
	list.Filter(dst, newList(1), nil)
	assert.Equal(t, 0, dst.Len())

	// Distinct dst and src
	src := newList(1, 2, 3, 4, 5)
	dst = arraylist.New[int]()
	list.Filter(list.List[int](dst), src, func(i, v int) bool { return v%2 == 0 })
	assert.Equal(t, []int{2, 4}, dst.Values())

	// Aliased dst and src: buffer branch.
	al := arraylist.New[int]()
	al.PushBack(1, 2, 3, 4, 5, 6)
	list.Filter(al, al, func(i, v int) bool { return v <= 3 })
	// Original 6 values + 3 appended matches = 9.
	assert.Equal(t, 9, al.Len())
}

func TestMapFn(t *testing.T) {
	// nil inputs
	dst := arraylist.New[int]()
	list.Map(dst, nil, func(i, v int) int { return v })
	assert.Equal(t, 0, dst.Len())
	list.Map(nil, newList(1), func(i, v int) int { return v })
	list.Map(dst, newList(1), nil)
	assert.Equal(t, 0, dst.Len())

	// Different types, distinct dst/src
	srcInt := newList(1, 2, 3)
	dstStr := arraylist.New[string]()
	list.Map(list.List[string](dstStr), srcInt, func(i, v int) string {
		return map[int]string{1: "a", 2: "b", 3: "c"}[v]
	})
	assert.Equal(t, []string{"a", "b", "c"}, dstStr.Values())

	// Aliased dst and src (same types): buffered append branch.
	al := arraylist.New[int]()
	al.PushBack(1, 2, 3)
	list.Map(al, al, func(i, v int) int { return v * 10 })
	assert.Equal(t, []int{1, 2, 3, 10, 20, 30}, al.Values())
}
