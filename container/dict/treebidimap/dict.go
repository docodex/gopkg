// Package treebidimap implements a bidirectional map backed by two red-black trees.
//
// In computer science, a bidirectional map is an associative data structure in which
// the (key, value) pairs form a one-to-one correspondence. Thus the binary relation
// is functional in each direction: each value can also be mapped to a unique key. A
// pair (a, b) thus provides a unique coupling between 'a' and 'b' so that 'b' can be
// found when 'a' is used as a key and 'a' can be found when 'b' is used as a key.
//
// Reference: https://en.wikipedia.org/wiki/Bidirectional_map
package treebidimap

import (
	"cmp"
	"encoding/json"
	"sync"

	"github.com/docodex/gopkg/container"
	"github.com/docodex/gopkg/container/tree/redblacktree"
	"github.com/docodex/gopkg/jsonx"
)

// Map represents a bidirectional treemap which holds the entries in two red-black trees.
type Map[K comparable, V comparable] struct {
	forward *redblacktree.Tree[K, V] // current forward map entries
	inverse *redblacktree.Tree[V, K] // current inverse map entries
	mu      *sync.RWMutex            // for concurrent use
}

// New returns an initialized bidirectional map with [cmp.Compare] as the cmp function for the
// backing red-black trees.
func New[K cmp.Ordered, V cmp.Ordered]() *Map[K, V] {
	return &Map[K, V]{
		forward: redblacktree.New[K, V](),
		inverse: redblacktree.New[V, K](),
		mu:      nil,
	}
}

// NewFunc returns an initialized bidirectional map with the given functions cmp as the cmp
// function for the backing red-black trees.
func NewFunc[K comparable, V comparable](cmp1 container.Compare[K], cmp2 container.Compare[V]) *Map[K, V] {
	return &Map[K, V]{
		forward: redblacktree.NewFunc[K, V](cmp1),
		inverse: redblacktree.NewFunc[V, K](cmp2),
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
	return m.forward.Len()
}

// Values returns all values in map.
func (m *Map[K, V]) Values() []V {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	values := make([]V, 0, m.inverse.Len())
	m.inverse.Range(func(k V, v K) bool {
		values = append(values, k)
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
	keys := make([]K, 0, m.forward.Len())
	m.forward.Range(func(k K, v V) bool {
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
	m1 := make(map[K]V, m.forward.Len())
	m.forward.Range(func(k K, v V) bool {
		m1[k] = v
		return true
	})
	entries, _ := jsonx.MarshalToString(m1)
	return "TreeBidiMap: " + entries
}

// MarshalJSON marshals map into valid JSON.
// Ref: std json.Marshaler.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	return m.forward.MarshalJSON()
}

// UnmarshalJSON unmarshals a JSON description of map.
// The input can be assumed to be a valid encoding of a JSON value.
// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
// Ref: std json.Unmarshaler.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	var m1 map[K]V
	if err := json.Unmarshal(data, &m1); err != nil {
		return err
	}
	m2 := make(map[V]K, len(m1))
	for k, v := range m1 {
		m2[v] = k
	}
	inverseData, err := json.Marshal(m2)
	if err != nil {
		return err
	}
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	if err := m.forward.UnmarshalJSON(data); err != nil {
		return err
	}
	if err := m.inverse.UnmarshalJSON(inverseData); err != nil {
		return err
	}
	return nil
}

// Put adds the key-value pair (k, v) to map.
func (m *Map[K, V]) Put(k K, v V) {
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	if v1, ok := m.forward.Get(k); ok {
		m.inverse.Remove(v1)
	}
	if k1, ok := m.inverse.Get(v); ok {
		m.forward.Remove(k1)
	}
	m.forward.Insert(k, v)
	m.inverse.Insert(v, k)
}

// Get returns the corresponding value of the given key k if exists in map.
// The ok result indicates whether such value was found in map.
func (m *Map[K, V]) Get(k K) (value V, ok bool) {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	value, ok = m.forward.Get(k)
	return
}

// GetKey returns the corresponding key of the given value v if exists in map.
// The ok result indicates whether such key was found in map.
func (m *Map[K, V]) GetKey(v V) (key K, ok bool) {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	key, ok = m.inverse.Get(v)
	return
}

// Remove removes the given key k and the corresponding value if exists in map.
// If there is no such key and value found in map, do nothing.
func (m *Map[K, V]) Remove(k K) {
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	if v, ok := m.forward.Get(k); ok {
		m.inverse.Remove(v)
		m.forward.Remove(k)
	}
}

// RemoveValue removes the value v and the corresponding key if exists in map.
// If there is no such value and key found in map, do nothing.
func (m *Map[K, V]) RemoveValue(v V) {
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	if k, ok := m.inverse.Get(v); ok {
		m.forward.Remove(k)
		m.inverse.Remove(v)
	}
}

// Contains returns true if map contains all of the given keys k.
func (m *Map[K, V]) Contains(k ...K) bool {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	for i := range k {
		if m.forward.Search(k[i]) == nil {
			return false
		}
	}
	return true
}

// ContainsValues returns true if map contains all of the given values v.
func (m *Map[K, V]) ContainsValues(v ...V) bool {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	for i := range v {
		if m.inverse.Search(v[i]) == nil {
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
		if m.forward.Search(k[i]) != nil {
			return true
		}
	}
	return false
}

// ContainsAnyValues returns true if map contains any of the given values v.
func (m *Map[K, V]) ContainsAnyValues(v ...V) bool {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	for i := range v {
		if m.inverse.Search(v[i]) != nil {
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
	m.forward.Clear()
	m.inverse.Clear()
}

// Range calls f sequentially for each key-value pair present in map.
func (m *Map[K, V]) Range(f func(k K, v V)) {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	m.forward.Range(func(k K, v V) bool {
		f(k, v)
		return true
	})
}
