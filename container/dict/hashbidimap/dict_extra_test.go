package hashbidimap_test

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/docodex/gopkg/container/dict/hashbidimap"
	"github.com/stretchr/testify/assert"
)

func factories() []*hashbidimap.Map[int, string] {
	return []*hashbidimap.Map[int, string]{
		hashbidimap.New[int, string](),
		hashbidimap.NewWithCapacity[int, string](0),
		hashbidimap.NewWithCapacity[int, string](64),
		hashbidimap.NewWithLock[int, string](),
		hashbidimap.NewWithCapacityAndLock[int, string](0),
		hashbidimap.NewWithCapacityAndLock[int, string](64),
	}
}

func TestConstructors(t *testing.T) {
	for _, m := range factories() {
		m.Put(1, "a")
		assert.Equal(t, 1, m.Len())
		v, ok := m.Get(1)
		assert.True(t, ok)
		assert.Equal(t, "a", v)
		k, ok := m.GetKey("a")
		assert.True(t, ok)
		assert.Equal(t, 1, k)
	}
}

func TestPutReplacesExistingKey(t *testing.T) {
	m := hashbidimap.New[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(1, "c") // replace value for key 1
	assert.Equal(t, 2, m.Len())
	v, _ := m.Get(1)
	assert.Equal(t, "c", v)
	_, ok := m.GetKey("a")
	assert.False(t, ok)
}

func TestPutReplacesExistingValue(t *testing.T) {
	m := hashbidimap.New[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "a") // key 3 takes value "a"; key 1 must be evicted
	assert.Equal(t, 2, m.Len())
	_, ok := m.Get(1)
	assert.False(t, ok)
	k, _ := m.GetKey("a")
	assert.Equal(t, 3, k)
}

func TestRemoveAndRemoveValue(t *testing.T) {
	for _, m := range factories() {
		m.Put(1, "a")
		m.Put(2, "b")
		m.Put(3, "c")

		m.Remove(2)
		_, ok := m.Get(2)
		assert.False(t, ok)
		_, ok = m.GetKey("b")
		assert.False(t, ok)

		m.RemoveValue("c")
		_, ok = m.GetKey("c")
		assert.False(t, ok)
		_, ok = m.Get(3)
		assert.False(t, ok)

		// Removing non-existent key/value is a no-op.
		m.Remove(999)
		m.RemoveValue("nope")
		assert.Equal(t, 1, m.Len())
	}
}

func TestContainsFamily(t *testing.T) {
	for _, m := range factories() {
		m.Put(1, "a")
		m.Put(2, "b")
		m.Put(3, "c")

		assert.True(t, m.Contains(1, 2, 3))
		assert.False(t, m.Contains(1, 9))
		assert.True(t, m.ContainsAny(9, 2))
		assert.False(t, m.ContainsAny(9, 8))

		assert.True(t, m.ContainsValues("a", "b"))
		assert.False(t, m.ContainsValues("a", "z"))
		assert.True(t, m.ContainsAnyValues("z", "a"))
		assert.False(t, m.ContainsAnyValues("y", "z"))
	}
}

func TestClearAndRange(t *testing.T) {
	for _, m := range factories() {
		m.Put(1, "a")
		m.Put(2, "b")

		m.Range(nil) // nil f: no-op
		got := map[int]string{}
		m.Range(func(k int, v string) bool {
			got[k] = v
			return true
		})
		assert.Equal(t, map[int]string{1: "a", 2: "b"}, got)

		count := 0
		m.Range(func(k int, v string) bool { count++; return false })
		assert.Equal(t, 1, count)

		m.Clear()
		assert.Equal(t, 0, m.Len())
		assert.Empty(t, m.Keys())
		assert.Empty(t, m.Values())
	}
}

func TestLockedAll(t *testing.T) {
	m := hashbidimap.NewWithLock[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	assert.Equal(t, 3, m.Len())
	assert.Len(t, m.Keys(), 3)
	assert.Len(t, m.Values(), 3)
	assert.True(t, strings.HasPrefix(m.String(), "HashBidiMap"))

	data, err := m.MarshalJSON()
	assert.NoError(t, err)
	assert.NoError(t, m.UnmarshalJSON(data))
	assert.Equal(t, 3, m.Len())
}

func TestUnmarshalJSONErrors(t *testing.T) {
	m := hashbidimap.New[int, string]()
	// Invalid JSON
	assert.Error(t, m.UnmarshalJSON([]byte("not json")))

	// Duplicate values must leave map unchanged and return ErrDuplicateValue.
	m.Put(10, "x")
	err := m.UnmarshalJSON([]byte(`{"1":"a","2":"a"}`))
	assert.True(t, errors.Is(err, hashbidimap.ErrDuplicateValue))
	v, ok := m.Get(10)
	assert.True(t, ok)
	assert.Equal(t, "x", v)

	// Valid input (locked variant).
	ml := hashbidimap.NewWithLock[int, string]()
	assert.NoError(t, ml.UnmarshalJSON([]byte(`{"1":"a","2":"b"}`)))
	assert.Equal(t, 2, ml.Len())
}

func TestConcurrentAccess(t *testing.T) {
	m := hashbidimap.NewWithCapacityAndLock[int, int](16)
	var wg sync.WaitGroup
	for w := range 4 {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for i := range 50 {
				key := base*1000 + i
				val := key // 1-to-1
				m.Put(key, val)
				_, _ = m.Get(key)
				_, _ = m.GetKey(val)
				_ = m.Contains(key)
				_ = m.ContainsValues(val)
				_ = m.ContainsAny(key)
				_ = m.ContainsAnyValues(val)
				_ = m.Len()
			}
		}(w)
	}
	wg.Wait()

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	m2 := hashbidimap.NewWithLock[int, int]()
	assert.NoError(t, json.Unmarshal(data, m2))
	assert.Equal(t, m.Len(), m2.Len())
}
