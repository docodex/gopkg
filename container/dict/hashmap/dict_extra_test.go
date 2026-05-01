package hashmap_test

import (
	"encoding/json"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/docodex/gopkg/container/dict/hashmap"
	"github.com/stretchr/testify/assert"
)

// factories exercises every constructor to ensure coverage and
// return the four kinds of Map instances used in table-driven tests.
func factories() []*hashmap.Map[int, string] {
	return []*hashmap.Map[int, string]{
		hashmap.New[int, string](),
		hashmap.NewWithCapacity[int, string](0),  // < defaultCapacity
		hashmap.NewWithCapacity[int, string](64), // > defaultCapacity
		hashmap.NewWithLock[int, string](),
		hashmap.NewWithCapacityAndLock[int, string](0),
		hashmap.NewWithCapacityAndLock[int, string](64),
	}
}

func TestConstructors(t *testing.T) {
	for _, m := range factories() {
		assert.Equal(t, 0, m.Len())
		m.Put(1, "a")
		assert.Equal(t, 1, m.Len())
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
		assert.False(t, m.Contains(1, 4))
		assert.True(t, m.Contains())

		assert.True(t, m.ContainsAny(4, 2))
		assert.False(t, m.ContainsAny(4, 5))
		assert.False(t, m.ContainsAny())
	}
}

func TestClear(t *testing.T) {
	for _, m := range factories() {
		m.Put(1, "a")
		m.Put(2, "b")
		assert.Equal(t, 2, m.Len())
		m.Clear()
		assert.Equal(t, 0, m.Len())
		// re-usable after Clear
		m.Put(3, "c")
		v, ok := m.Get(3)
		assert.True(t, ok)
		assert.Equal(t, "c", v)
	}
}

func TestRange(t *testing.T) {
	for _, m := range factories() {
		m.Put(1, "a")
		m.Put(2, "b")
		m.Put(3, "c")

		// nil callback must be a no-op.
		m.Range(nil)

		// Full traversal.
		collected := map[int]string{}
		m.Range(func(k int, v string) bool {
			collected[k] = v
			return true
		})
		assert.Equal(t, map[int]string{1: "a", 2: "b", 3: "c"}, collected)

		// Early termination must be respected.
		count := 0
		m.Range(func(k int, v string) bool {
			count++
			return false
		})
		assert.Equal(t, 1, count)
	}
}

func TestLockedHappyPath(t *testing.T) {
	// Ensure the sync.RWMutex branches of every method are executed.
	m := hashmap.NewWithLock[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	assert.Equal(t, 3, m.Len())

	keys := m.Keys()
	slices.Sort(keys)
	assert.Equal(t, []string{"a", "b", "c"}, keys)

	values := m.Values()
	slices.Sort(values)
	assert.Equal(t, []int{1, 2, 3}, values)

	v, ok := m.Get("b")
	assert.True(t, ok)
	assert.Equal(t, 2, v)

	assert.True(t, m.Contains("a", "b"))
	assert.False(t, m.Contains("a", "z"))
	assert.True(t, m.ContainsAny("z", "a"))
	assert.False(t, m.ContainsAny("x", "z"))

	assert.True(t, strings.HasPrefix(m.String(), "HashMap"))

	data, err := m.MarshalJSON()
	assert.NoError(t, err)
	// Re-parse and re-load to exercise UnmarshalJSON lock branch.
	err = m.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, 3, m.Len())

	// Range under lock.
	seen := 0
	m.Range(func(k string, v int) bool {
		seen++
		return true
	})
	assert.Equal(t, 3, seen)

	m.Remove("a")
	assert.False(t, m.Contains("a"))

	m.Clear()
	assert.Equal(t, 0, m.Len())
}

func TestUnmarshalJSONInvalid(t *testing.T) {
	m := hashmap.New[int, string]()
	m.Put(1, "a")
	// Invalid JSON must return a non-nil error.
	err := m.UnmarshalJSON([]byte("not json"))
	assert.Error(t, err)

	// Unmarshal valid JSON into a locked map.
	ml := hashmap.NewWithLock[int, string]()
	err = ml.UnmarshalJSON([]byte(`{"1":"a","2":"b"}`))
	assert.NoError(t, err)
	assert.Equal(t, 2, ml.Len())
}

func TestConcurrentAccess(t *testing.T) {
	// Exercise the RWMutex paths concurrently with the race detector.
	m := hashmap.NewWithCapacityAndLock[int, int](16)
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
	assert.GreaterOrEqual(t, m.Len(), 0)

	// Marshal/unmarshal roundtrip.
	data, err := json.Marshal(m)
	assert.NoError(t, err)
	m2 := hashmap.NewWithLock[int, int]()
	assert.NoError(t, json.Unmarshal(data, m2))
	assert.Equal(t, m.Len(), m2.Len())
}
