// Package list provides an abstract List interface.
//
// In computer science, a list or sequence is an abstract data type that represents
// an ordered sequence of values, where the same value may occur more than once. An
// instance of a list is a computer representation of the mathematical concept of a
// finite sequence; the (potentially) infinite analog of a list is a stream.  Lists
// are a basic example of containers, as they contain other values. If the same
// value occurs multiple times, each occurrence is considered a distinct item.
//
// Reference: https://en.wikipedia.org/wiki/List_%28abstract_data_type%29
//
// Concurrency: implementations are NOT safe for concurrent use. Callers must
// provide external synchronization if a list is shared across goroutines.
// See package container for the project-wide concurrency policy.
package list

import "github.com/docodex/gopkg/container"

type List[T any] interface {
	container.Container[T]

	// MarshalJSON marshals list into valid JSON.
	// Ref: std json.Marshaler.
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals a JSON description of list.
	// The input can be assumed to be a valid encoding of a JSON value.
	// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
	// Ref: std json.Unmarshaler.
	UnmarshalJSON(data []byte) error

	// Front returns the first element if exists in list.
	// The ok result indicates whether such element was found in list.
	Front() (value T, ok bool)
	// Back returns the last element if exists in list.
	// The ok result indicates whether such element was found in list.
	Back() (value T, ok bool)
	// PushFront inserts new elements with the given values v at the front of list.
	PushFront(v ...T)
	// PushBack inserts new elements with the given values v at the back of list.
	PushBack(v ...T)
	// PopFront removes the first element if exists in list and returns it.
	// The ok result indicates whether such element was removed from list.
	PopFront() (value T, ok bool)
	// PopBack removes the last element if exists in list and returns it.
	// The ok result indicates whether such element was removed from list.
	PopBack() (value T, ok bool)
	// Clear removes all values in list.
	Clear()

	// Get returns the value of index i if exists in list.
	// The ok result indicates whether such value was found in list.
	Get(i int) (value T, ok bool)
	// Set sets the value to v of index i if exists in list.
	Set(i int, v T)
	// Add inserts the values v to index i if exists in list, or appends the value v to the back
	// of list if index i points to the next index of the last element in list.
	Add(i int, v ...T)
	// Del removes the value at index i if exists in list.
	Del(i int)
	// Swap swaps the values with indices i and j if both indices exist in list.
	Swap(i, j int)

	// Sort sorts list values (in-place) with the given cmp.
	Sort(cmp container.Compare[T])

	// Range calls f sequentially for each index i and value v present in list.
	// If f returns false, range stops the iteration.
	Range(f func(i int, v T) bool)
}

// Index returns the index of the first occurrence of value v in list l, or -1 if not present.
func Index[T comparable](l List[T], v T) (index int) {
	index = -1
	if l == nil || l.Len() == 0 {
		return
	}
	l.Range(func(i int, v1 T) bool {
		if v1 == v {
			index = i
			return false
		}
		return true
	})
	return
}

// Find returns the first index i and the corresponding value v in list l satisfying condition
// f(i, v), or first return parameter would be -1 if none do.
func Find[T any](l List[T], f func(i int, v T) bool) (index int, value T) {
	index = -1
	if l == nil || l.Len() == 0 || f == nil {
		return
	}
	l.Range(func(i int, v T) bool {
		if f(i, v) {
			index = i
			value = v
			return false
		}
		return true
	})
	return
}

// Contains returns true if list l contains all the given values v.
func Contains[T comparable](l List[T], v ...T) bool {
	if l == nil {
		return false
	}
	if len(v) == 0 {
		return true
	}
	if l.Len() == 0 {
		return false
	}
	// Deduplicate the query values so we can scan the list a single time.
	remaining := make(map[T]struct{}, len(v))
	for i := range v {
		remaining[v[i]] = struct{}{}
	}
	l.Range(func(_ int, v1 T) bool {
		delete(remaining, v1)
		return len(remaining) != 0
	})
	return len(remaining) == 0
}

// ContainsAny returns true if list l contains any of the given values v.
func ContainsAny[T comparable](l List[T], v ...T) bool {
	if l == nil || len(v) == 0 || l.Len() == 0 {
		return false
	}
	targets := make(map[T]struct{}, len(v))
	for i := range v {
		targets[v[i]] = struct{}{}
	}
	found := false
	l.Range(func(_ int, v1 T) bool {
		if _, ok := targets[v1]; ok {
			found = true
			return false
		}
		return true
	})
	return found
}

// All returns true if all elements in list l satisfying condition f(i, v) which i is the
// element index and v is the corresponding value of i.
//
// All returns true for an empty list and for a nil list (vacuous truth: a universal
// quantifier over an empty domain is true). A nil predicate f is treated as "unknown"
// and returns false.
func All[T any](l List[T], f func(i int, v T) bool) bool {
	if f == nil {
		return false
	}
	if l == nil || l.Len() == 0 {
		return true
	}
	except := false
	l.Range(func(i int, v T) bool {
		if !f(i, v) {
			except = true
			return false
		}
		return true
	})
	return !except
}

// Any returns true if any of elements in list l satisfying condition f(i, v) which i is the
// element index and v is the corresponding value of i, or false if none do.
func Any[T any](l List[T], f func(i int, v T) bool) bool {
	if l == nil || l.Len() == 0 || f == nil {
		return false
	}
	some := false
	l.Range(func(i int, v T) bool {
		if f(i, v) {
			some = true
			return false
		}
		return true
	})
	return some
}

// Filter filters elements in list src to list dst by condition f(i, v) which i is the element
// index and v is the corresponding value of i.
//
// Filter is safe when dst and src reference the same list: matching values are buffered before
// being appended, so concurrent mutation of src during iteration cannot occur.
func Filter[T any](dst, src List[T], f func(i int, v T) bool) {
	if src == nil || f == nil || dst == nil {
		return
	}
	if any(dst) == any(src) {
		// Collect first, then append, to avoid mutating src during iteration.
		var buf []T
		src.Range(func(i int, v T) bool {
			if f(i, v) {
				buf = append(buf, v)
			}
			return true
		})
		dst.PushBack(buf...)
		return
	}
	src.Range(func(i int, v T) bool {
		if f(i, v) {
			dst.PushBack(v)
		}
		return true
	})
}

// Map maps elements in list src to list dst by condition f(i, v) which i is the element index and
// v is the corresponding value of i.
//
// When T1 and T2 are the same type and dst and src reference the same list, Map buffers the
// mapped values before appending to avoid mutating src during iteration.
func Map[T1, T2 any](dst List[T2], src List[T1], f func(i int, v T1) T2) {
	if src == nil || f == nil || dst == nil {
		return
	}
	if any(dst) == any(src) {
		var buf []T2
		src.Range(func(i int, v T1) bool {
			buf = append(buf, f(i, v))
			return true
		})
		dst.PushBack(buf...)
		return
	}
	src.Range(func(i int, v T1) bool {
		dst.PushBack(f(i, v))
		return true
	})
}
