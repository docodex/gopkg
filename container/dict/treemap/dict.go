// Package treemap implements a map backed by a red-black tree.
package treemap

import (
	"cmp"
	"sync"

	"github.com/docodex/gopkg/container"
	"github.com/docodex/gopkg/container/tree/redblacktree"
	"github.com/docodex/gopkg/jsonx"
)

// Map represents a treemap which holds the entries in a red-black tree.
type Map[K comparable, V any] struct {
	entries *redblacktree.Tree[K, V] // current map entries
	mu      *sync.RWMutex            // for concurrent use
}

// New returns an initialized map with [cmp.Compare] as the cmp function for the backing red-black
// tree.
func New[K cmp.Ordered, V any]() *Map[K, V] {
	return &Map[K, V]{
		entries: redblacktree.New[K, V](),
		mu:      nil,
	}
}

// NewFunc returns an initialized map with the given function cmp as the cmp function for the
// backing red-black tree.
func NewFunc[K comparable, V any](cmp container.Compare[K]) *Map[K, V] {
	return &Map[K, V]{
		entries: redblacktree.NewFunc[K, V](cmp),
		mu:      nil,
	}
}

// WithLock adds sync.RWMutex to support concurrent use by multiple goroutines without additional
// locking or coordination.
func (m *Map[K, V]) WithLock() *Map[K, V] {
	m.mu = &sync.RWMutex{}
	return m
}

// Len returns the number of entries of map m.
func (m *Map[K, V]) Len() int {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	return m.entries.Len()
}

// Values returns all values in map.
func (m *Map[K, V]) Values() []V {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	values := make([]V, 0, m.entries.Len())
	m.entries.Range(func(k K, v V) bool {
		values = append(values, v)
		return true
	})
	return values
}

// Values returns all keys in map.
func (m *Map[K, V]) Keys() []K {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	keys := make([]K, 0, m.entries.Len())
	m.entries.Range(func(k K, v V) bool {
		keys = append(keys, k)
		return true
	})
	return keys
}

// String returns the string representation of map.
// Ref: std fmt.Stringer.
func (m *Map[K, V]) String() string {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	m1 := make(map[K]V, m.entries.Len())
	m.entries.Range(func(k K, v V) bool {
		m1[k] = v
		return true
	})
	entries, _ := jsonx.MarshalToString(m1)
	return "TreeMap: " + entries
}

// MarshalJSON marshals map into valid JSON.
// Ref: std json.Marshaler.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	return m.entries.MarshalJSON()
}

// UnmarshalJSON unmarshals a JSON description of map.
// The input can be assumed to be a valid encoding of a JSON value.
// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
// Ref: std json.Unmarshaler.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	return m.entries.UnmarshalJSON(data)
}

// Put adds the key-value pair (k, v) to map.
func (m *Map[K, V]) Put(k K, v V) {
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	m.entries.Insert(k, v)
}

// Get returns the corresponding value of the given key k if exists in map.
// The ok result indicates whether such value was found in map.
func (m *Map[K, V]) Get(k K) (value V, ok bool) {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	value, ok = m.entries.Get(k)
	return
}

// Remove removes the given key k and the corresponding value if exists in map.
// If there is no such key and value found in map, do nothing.
func (m *Map[K, V]) Remove(k K) {
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	m.entries.Remove(k)
}

// Contains returns true if map contains all of the given keys k.
func (m *Map[K, V]) Contains(k ...K) bool {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	for i := range k {
		if m.entries.Search(k[i]) == nil {
			return false
		}
	}
	return true
}

// Contains returns true if map contains any of the given keys k.
func (m *Map[K, V]) ContainsAny(k ...K) bool {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	for i := range k {
		if m.entries.Search(k[i]) != nil {
			return true
		}
	}
	return false
}

// Clear removes all key-value pairs in map.
func (m *Map[K, V]) Clear() {
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	m.entries.Clear()
}

// Range calls f sequentially for each key-value pair present in map.
func (m *Map[K, V]) Range(f func(k K, v V)) {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	m.entries.Range(func(k K, v V) bool {
		f(k, v)
		return true
	})
}
