// Package hashset implements a set backed by a hash table.
package hashset

import (
	"encoding/json"
	"sync"

	"github.com/docodex/gopkg/jsonx"
)

const defaultCapacity = 32

// Set represents a hashset which holds the values in a hash table.
type Set[T comparable] struct {
	values map[T]struct{} // current set values
	mu     *sync.RWMutex  // for concurrent use
}

// New returns an initialized set with the default capacity as the initial capacity for the
// backing hash table.
func New[T comparable](v ...T) *Set[T] {
	s := &Set[T]{
		values: make(map[T]struct{}, max(len(v), defaultCapacity)),
		mu:     nil,
	}
	for i := range v {
		s.values[v[i]] = struct{}{}
	}
	return s
}

// NewWithCapacity returns an initialized set with the given capacity as the initial capacity for
// the backing hash table.
func NewWithCapacity[T comparable](capacity int) *Set[T] {
	return &Set[T]{
		values: make(map[T]struct{}, max(capacity, defaultCapacity)),
		mu:     nil,
	}
}

// WithLock adds sync.RWMutex to support concurrent use by multiple goroutines without additional
// locking or coordination.
func (s *Set[T]) WithLock() *Set[T] {
	s.mu = &sync.RWMutex{}
	return s
}

// Len returns the number of values of set s.
func (s *Set[T]) Len() int {
	if s.mu != nil {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	return len(s.values)
}

// Values returns all values in set.
func (s *Set[T]) Values() []T {
	if s.mu != nil {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	values := make([]T, 0, len(s.values))
	for value := range s.values {
		values = append(values, value)
	}
	return values
}

// String returns the string representation of set.
// Ref: std fmt.Stringer.
func (s *Set[T]) String() string {
	values, _ := jsonx.MarshalToString(s.Values())
	return "HashSet: " + values
}

// MarshalJSON marshals set into valid JSON.
// Ref: std json.Marshaler.
func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON unmarshals a JSON description of set.
// The input can be assumed to be a valid encoding of a JSON value.
// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
// Ref: std json.Unmarshaler.
func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var v []T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if s.mu != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	s.values = make(map[T]struct{}, max(len(v), defaultCapacity))
	for i := range v {
		s.values[v[i]] = struct{}{}
	}
	return nil
}

// Add adds the given values v to set.
func (s *Set[T]) Add(v ...T) {
	if s.mu != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	for i := range v {
		s.values[v[i]] = struct{}{}
	}
}

// Remove removes the given values v if exists in set.
// If there is no such values found in set, do nothing.
func (s *Set[T]) Remove(v ...T) {
	if s.mu != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	for i := range v {
		delete(s.values, v[i])
	}
}

// Contains returns true if set contains all of the given values v.
func (s *Set[T]) Contains(v ...T) bool {
	if s.mu != nil {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	for i := range v {
		if _, ok := s.values[v[i]]; !ok {
			return false
		}
	}
	return true
}

// Contains returns true if set contains any of the given values v.
func (s *Set[T]) ContainsAny(v ...T) bool {
	if s.mu != nil {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	for i := range v {
		if _, ok := s.values[v[i]]; ok {
			return true
		}
	}
	return false
}

// Clear removes all values in set.
func (s *Set[T]) Clear() {
	if s.mu != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	s.values = make(map[T]struct{}, defaultCapacity)
}

// Range calls f for each value v present in the set.
func (s *Set[T]) Range(f func(v T)) {
	if f == nil {
		return
	}
	if s.mu != nil {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	for v := range s.values {
		f(v)
	}
}
