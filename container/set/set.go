// Package set provides an abstract Set interface.
//
// In computer science, a set is an abstract data type that can store certain values
// and no repeated values. It is a computer implementation of the mathematical concept
// of a finite set. Unlike most other collection types, rather than retrieving a specific
// element from a set, one typically tests a value for membership in a set.
//
// Reference: https://en.wikipedia.org/wiki/Set_%28abstract_data_type%29
package set

import "github.com/docodex/gopkg/container"

type Set[T any] interface {
	container.Container[T]

	// MarshalJSON marshals set into valid JSON.
	// Ref: std json.Marshaler.
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals a JSON description of set.
	// The input can be assumed to be a valid encoding of a JSON value.
	// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
	// Ref: std json.Unmarshaler.
	UnmarshalJSON(data []byte) error

	// Add adds the given values v to set.
	Add(v ...T)
	// Remove removes the given values v if exists in set.
	// If there is no such values found in set, do nothing.
	Remove(v ...T)
	// Contains returns true if set contains all of the given values v.
	Contains(v ...T) bool
	// Contains returns true if set contains any of the given values v.
	ContainsAny(v ...T) bool
	// Clear removes all values in set.
	Clear()

	// Range calls f for each value v present in the set.
	Range(f func(v T))
}

// Intersection returns the intersection among sets.
// The dst set consists of all elements that are in all "src" sets.
// Ref: https://en.wikipedia.org/wiki/Intersection_(set_theory)
func Intersection[T comparable](dst Set[T], src ...Set[T]) {
	if dst == nil {
		return
	}
	var (
		tmp = -1 // shortest set length
		pos = -1 // shortest set index
	)
	for i := range src {
		if src[i] == nil || src[i].Len() == 0 {
			// nil or empty set exists
			return
		}
		if tmp == -1 || src[i].Len() < tmp {
			tmp = src[i].Len()
			pos = i
		}
	}
	if pos == -1 {
		return
	}
	src[pos].Range(func(v T) {
		for i := range src {
			if i != pos && !src[i].Contains(v) {
				return
			}
		}
		dst.Add(v)
	})
}

// Union returns the union of sets.
// The dst set consists of all elements that are in any "src" sets.
// Ref: https://en.wikipedia.org/wiki/Union_(set_theory)
func Union[T comparable](dst Set[T], src ...Set[T]) {
	if dst == nil {
		return
	}
	for i := range src {
		if src[i] != nil {
			src[i].Range(func(v T) {
				dst.Add(v)
			})
		}
	}
}

// Difference returns the difference between two sets.
// The dst set consists of all elements that are in "a" but not in "b".
// Ref: https://proofwiki.org/wiki/Definition:Set_Difference
func Difference[T comparable](dst, a, b Set[T]) {
	if dst == nil {
		return
	}
	a.Range(func(v T) {
		if !b.Contains(v) {
			dst.Add(v)
		}
	})
}
