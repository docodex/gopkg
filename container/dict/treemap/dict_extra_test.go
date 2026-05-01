package treemap_test

import (
	"cmp"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/docodex/gopkg/container/dict/treemap"
	"github.com/stretchr/testify/assert"
)

func factories() []*treemap.Map[int, string] {
	reverse := func(a, b int) int { return cmp.Compare(b, a) }
	return []*treemap.Map[int, string]{
		treemap.New[int, string](),
		treemap.NewFunc[int, string](reverse),
		treemap.NewWithLock[int, string](),
		treemap.NewFuncWithLock[int, string](reverse),
	}
}

func TestConstructors(t *testing.T) {
	for _, m := range factories() {
		assert.Equal(t, 0, m.Len())
		m.Put(1, "a")
		m.Put(2, "b")
		assert.Equal(t, 2, m.Len())
		v, ok := m.Get(1)
		assert.True(t, ok)
		assert.Equal(t, "a", v)
	}
}

func TestContainsAndContainsAny(t *testing.T) {
	for _, m := range factories() {
		m.Put(1, "a")
		m.Put(2, "b")
		m.Put(3, "c")
		assert.True(t, m.Contains(1, 2, 3))
		assert.False(t, m.Contains(1, 9))
		assert.True(t, m.Contains())
		assert.True(t, m.ContainsAny(9, 2))
		assert.False(t, m.ContainsAny(8, 9))
		assert.False(t, m.ContainsAny())
	}
}

func TestClearAndRange(t *testing.T) {
	for _, m := range factories() {
		m.Put(1, "a")
		m.Put(2, "b")
		m.Put(3, "c")

		m.Range(nil) // nil f: no-op
		collected := map[int]string{}
		m.Range(func(k int, v string) bool {
			collected[k] = v
			return true
		})
		assert.Equal(t, map[int]string{1: "a", 2: "b", 3: "c"}, collected)

		count := 0
		m.Range(func(k int, v string) bool {
			count++
			return false
		})
		assert.Equal(t, 1, count)

		m.Clear()
		assert.Equal(t, 0, m.Len())
		m.Put(4, "d")
		assert.Equal(t, 1, m.Len())
	}
}

func TestLockedAll(t *testing.T) {
	m := treemap.NewWithLock[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	assert.Equal(t, 3, m.Len())
	assert.Len(t, m.Keys(), 3)
	assert.Len(t, m.Values(), 3)
	assert.True(t, strings.HasPrefix(m.String(), "TreeMap"))
	assert.True(t, m.Contains("a"))
	assert.False(t, m.Contains("z"))
	assert.True(t, m.ContainsAny("z", "a"))
	assert.False(t, m.ContainsAny("z"))

	data, err := m.MarshalJSON()
	assert.NoError(t, err)
	assert.NoError(t, m.UnmarshalJSON(data))
	assert.Equal(t, 3, m.Len())

	v, ok := m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	seen := 0
	m.Range(func(k string, v int) bool { seen++; return true })
	assert.Equal(t, 3, seen)

	m.Remove("a")
	assert.False(t, m.Contains("a"))
	m.Clear()
	assert.Equal(t, 0, m.Len())
}

func TestUnmarshalInvalid(t *testing.T) {
	m := treemap.New[int, string]()
	assert.Error(t, m.UnmarshalJSON([]byte("{not json")))

	ml := treemap.NewWithLock[int, string]()
	assert.NoError(t, ml.UnmarshalJSON([]byte(`{"1":"a","2":"b"}`)))
	assert.Equal(t, 2, ml.Len())
}

func TestConcurrentAccess(t *testing.T) {
	m := treemap.NewWithLock[int, int]()
	var wg sync.WaitGroup
	for w := range 4 {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for i := range 50 {
				m.Put(base*100+i, i)
				_, _ = m.Get(base*100 + i)
				_ = m.Contains(base * 100)
				_ = m.ContainsAny(base * 100)
				_ = m.Len()
			}
		}(w)
	}
	wg.Wait()

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	m2 := treemap.NewWithLock[int, int]()
	assert.NoError(t, json.Unmarshal(data, m2))
	assert.Equal(t, m.Len(), m2.Len())
}
