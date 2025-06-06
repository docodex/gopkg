// Package set provides an abstract Map interface.
package dict

import "github.com/docodex/gopkg/container"

type Map[K comparable, V any] interface {
	container.Container[V]

	// Keys returns a slice of all keys in map.
	Keys() []K

	// MarshalJSON marshals map into valid JSON.
	// Ref: std json.Marshaler.
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals a JSON description of map.
	// The input can be assumed to be a valid encoding of a JSON value.
	// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
	// Ref: std json.Unmarshaler.
	UnmarshalJSON(data []byte) error

	// Put adds the key-value pair (k, v) to map.
	Put(k K, v V)
	// Get returns the corresponding value of the given key k if exists in map.
	// The ok result indicates whether such value was found in map.
	Get(k K) (value V, ok bool)
	// Remove removes the given key k and the corresponding value if exists in map.
	// If there is no such key and value found in map, do nothing.
	Remove(k K)
	// Contains returns true if map contains all of the given keys k.
	Contains(k ...K) bool
	// Contains returns true if map contains any of the given keys k.
	ContainsAny(k ...K) bool
	// Clear removes all key-value pairs in map.
	Clear()

	// Range calls f for each key-value pair present in map.
	Range(f func(k K, v V))
}

type BidiMap[K comparable, V comparable] interface {
	Map[K, V]

	// GetKey returns the corresponding key of the given value v if exists in map.
	// The ok result indicates whether such key was found in map.
	GetKey(v V) (key K, ok bool)
	// RemoveValue removes the value v and the corresponding key if exists in map.
	// If there is no such value and key found in map, do nothing.
	RemoveValue(v V)
	// ContainsValues returns true if map contains all of the given values v.
	ContainsValues(v ...V) bool
	// ContainsAnyValues returns true if map contains any of the given values v.
	ContainsAnyValues(v ...V) bool
}
