// Package treeset implements a set backed by a red-black tree.
package treeset

import (
	"cmp"
	"encoding/json"
	"sync"

	"github.com/docodex/gopkg/container"
	"github.com/docodex/gopkg/container/tree/redblacktree"
	"github.com/docodex/gopkg/jsonx"
)

// Map represents a treeset which holds the values in a red-black tree.
type Set[T comparable] struct {
	values *redblacktree.Tree[T, struct{}] // current set values
	mu     *sync.RWMutex                   // for concurrent use
}

// New returns an initialized set with [cmp.Compare] as the cmp function for the backing red-black
// tree.
func New[T cmp.Ordered](v ...T) *Set[T] {
	s := &Set[T]{
		values: redblacktree.New[T, struct{}](),
		mu:     nil,
	}
	for i := range v {
		s.values.Insert(v[i], struct{}{})
	}
	return s
}

// NewFunc returns an initialized set with the given function cmp as the cmp function for the
// backing red-black tree.
func NewFunc[T comparable](cmp container.Compare[T]) *Set[T] {
	s := &Set[T]{
		values: redblacktree.NewFunc[T, struct{}](cmp),
		mu:     nil,
	}
	return s
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
	return s.values.Len()
}

// Values returns all values in set.
func (s *Set[T]) Values() []T {
	if s.mu != nil {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	return s.values.Keys()
}

// String returns the string representation of set.
// Ref: std fmt.Stringer.
func (s *Set[T]) String() string {
	values, _ := jsonx.MarshalToString(s.Values())
	return "TreeSet: " + values
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
	m := make(map[T]struct{}, len(v))
	for i := range v {
		m[v[i]] = struct{}{}
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if s.mu != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	return s.values.UnmarshalJSON(data)
}

// Add adds the given values v to set.
func (s *Set[T]) Add(v ...T) {
	if s.mu != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	for i := range v {
		s.values.Insert(v[i], struct{}{})
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
		s.values.Remove(v[i])
	}
}

// Contains returns true if set contains all of the given values v.
func (s *Set[T]) Contains(v ...T) bool {
	if s.mu != nil {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	for i := range v {
		if s.values.Search(v[i]) == nil {
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
		if s.values.Search(v[i]) != nil {
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
	s.values.Clear()
}

// Range calls f sequentially for each value v present in the set.
func (s *Set[T]) Range(f func(v T)) {
	if f == nil {
		return
	}
	if s.mu != nil {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	v := s.values.Keys()
	for i := range v {
		f(v[i])
	}
}
