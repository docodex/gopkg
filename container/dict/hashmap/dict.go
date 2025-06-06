// Package hashmap implements a map backed by a hash table.
package hashmap

import (
	"encoding/json"
	"maps"
	"sync"

	"github.com/docodex/gopkg/jsonx"
)

const defaultCapacity = 32

// Map represents a hashmap which holds the entries in a hash table.
type Map[K comparable, V any] struct {
	entries map[K]V       // current map entries
	mu      *sync.RWMutex // for concurrent use
}

// New returns an initialized map with the default capacity as the initial capacity for the
// backing hash table.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		entries: make(map[K]V, defaultCapacity),
		mu:      nil,
	}
}

// NewWithCapacity returns an initialized map with the given capacity as the initial capacity for
// the backing hash table.
func NewWithCapacity[K comparable, V any](capacity int) *Map[K, V] {
	return &Map[K, V]{
		entries: make(map[K]V, max(capacity, defaultCapacity)),
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
	return len(m.entries)
}

// Values returns all values in map.
func (m *Map[K, V]) Values() []V {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	values := make([]V, 0, len(m.entries))
	for _, v := range m.entries {
		values = append(values, v)
	}
	return values
}

// Values returns all keys in map.
func (m *Map[K, V]) Keys() []K {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	keys := make([]K, 0, len(m.entries))
	for k := range m.entries {
		keys = append(keys, k)
	}
	return keys
}

// String returns the string representation of map.
// Ref: std fmt.Stringer.
func (m *Map[K, V]) String() string {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	entries, _ := jsonx.MarshalToString(m.entries)
	return "HashMap: " + entries
}

// MarshalJSON marshals map into valid JSON.
// Ref: std json.Marshaler.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	return json.Marshal(m.entries)
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
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	m.entries = make(map[K]V, max(len(m1), defaultCapacity))
	maps.Copy(m.entries, m1)
	return nil
}

// Put adds the key-value pair (k, v) to map.
func (m *Map[K, V]) Put(k K, v V) {
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	m.entries[k] = v
}

// Get returns the corresponding value of the given key k if exists in map.
// The ok result indicates whether such value was found in map.
func (m *Map[K, V]) Get(k K) (value V, ok bool) {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	value, ok = m.entries[k]
	return
}

// Remove removes the given key k and the corresponding value if exists in map.
// If there is no such key and value found in map, do nothing.
func (m *Map[K, V]) Remove(k K) {
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	delete(m.entries, k)
}

// Contains returns true if map contains all of the given keys k.
func (m *Map[K, V]) Contains(k ...K) bool {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	for i := range k {
		if _, ok := m.entries[k[i]]; !ok {
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
		if _, ok := m.entries[k[i]]; ok {
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
	m.entries = make(map[K]V, defaultCapacity)
}

// Range calls f for each key-value pair present in map.
func (m *Map[K, V]) Range(f func(k K, v V)) {
	if f == nil {
		return
	}
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	for k, v := range m.entries {
		f(k, v)
	}
}
