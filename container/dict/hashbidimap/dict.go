// Package hashbidimap implements a bidirectional map backed by two hash tables.
//
// In computer science, a bidirectional map is an associative data structure in which
// the (key, value) pairs form a one-to-one correspondence. Thus the binary relation
// is functional in each direction: each value can also be mapped to a unique key. A
// pair (a, b) thus provides a unique coupling between 'a' and 'b' so that 'b' can be
// found when 'a' is used as a key and 'a' can be found when 'b' is used as a key.
//
// Reference: https://en.wikipedia.org/wiki/Bidirectional_map
package hashbidimap

import (
	"encoding/json"
	"sync"

	"github.com/docodex/gopkg/jsonx"
)

const defaultCapacity = 32

// Map represents a bidirectional hashmap which holds the entries in two hash tables.
type Map[K comparable, V comparable] struct {
	forward map[K]V       // current forward map entries
	inverse map[V]K       // current inverse map entries
	mu      *sync.RWMutex // for concurrent use
}

// New returns an initialized bidirectional map with the default capacity as the initial capacity
// for the backing hash tables.
func New[K comparable, V comparable]() *Map[K, V] {
	return &Map[K, V]{
		forward: make(map[K]V, defaultCapacity),
		inverse: make(map[V]K, defaultCapacity),
		mu:      nil,
	}
}

// NewWithCapacity returns an initialized bidirectional map with the given capacity as the initial
// capacity for the backing hash tables.
func NewWithCapacity[K comparable, V comparable](capacity int) *Map[K, V] {
	capacity = max(capacity, defaultCapacity)
	return &Map[K, V]{
		forward: make(map[K]V, capacity),
		inverse: make(map[V]K, capacity),
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
	return len(m.forward)
}

// Values returns all values in map.
func (m *Map[K, V]) Values() []V {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	values := make([]V, 0, len(m.inverse))
	for v := range m.inverse {
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
	keys := make([]K, 0, len(m.forward))
	for k := range m.forward {
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
	entries, _ := jsonx.MarshalToString(m.forward)
	return "HashBidiMap: " + entries
}

// MarshalJSON marshals map into valid JSON.
// Ref: std json.Marshaler.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	return json.Marshal(m.forward)
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
	capacity := max(len(m1), defaultCapacity)
	m.forward = make(map[K]V, capacity)
	m.inverse = make(map[V]K, capacity)
	for k, v := range m1 {
		m.forward[k] = v
		m.inverse[v] = k
	}
	return nil
}

// Put adds the key-value pair (k, v) to map.
func (m *Map[K, V]) Put(k K, v V) {
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	if v1, ok := m.forward[k]; ok {
		delete(m.inverse, v1)
	}
	if k1, ok := m.inverse[v]; ok {
		delete(m.forward, k1)
	}
	m.forward[k] = v
	m.inverse[v] = k
}

// Get returns the corresponding value of the given key k if exists in map.
// The ok result indicates whether such value was found in map.
func (m *Map[K, V]) Get(k K) (value V, ok bool) {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	value, ok = m.forward[k]
	return
}

// GetKey returns the corresponding key of the given value v if exists in map.
// The ok result indicates whether such key was found in map.
func (m *Map[K, V]) GetKey(v V) (key K, ok bool) {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	key, ok = m.inverse[v]
	return
}

// Remove removes the given key k and the corresponding value if exists in map.
// If there is no such key and value found in map, do nothing.
func (m *Map[K, V]) Remove(k K) {
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	if v, ok := m.forward[k]; ok {
		delete(m.inverse, v)
		delete(m.forward, k)
	}
}

// RemoveValue removes the value v and the corresponding key if exists in map.
// If there is no such value and key found in map, do nothing.
func (m *Map[K, V]) RemoveValue(v V) {
	if m.mu != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
	}
	if k, ok := m.inverse[v]; ok {
		delete(m.forward, k)
		delete(m.inverse, v)
	}
}

// Contains returns true if map contains all of the given keys k.
func (m *Map[K, V]) Contains(k ...K) bool {
	if m.mu != nil {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	for i := range k {
		if _, ok := m.forward[k[i]]; !ok {
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
		if _, ok := m.inverse[v[i]]; !ok {
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
		if _, ok := m.forward[k[i]]; ok {
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
		if _, ok := m.inverse[v[i]]; ok {
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
	m.forward = make(map[K]V, defaultCapacity)
	m.inverse = make(map[V]K, defaultCapacity)
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
	for k, v := range m.forward {
		f(k, v)
	}
}
