// Package arraylist implements an array list.
//
// To iterate over a list (where l is a *List):
//
//	for i := range l.Len() {
//		// do something with l.Get(i)
//	}
//
// or:
//
//	l.Range(func(index int, value T) bool {
//		// do something with index and value
//		return true
//	})
package arraylist

import (
	"encoding/json"
	"slices"

	"github.com/docodex/gopkg/container"
	"github.com/docodex/gopkg/jsonx"
)

// List represents an array list which holds the elements in a slice.
type List[T any] struct {
	values    []T // current list elements
	low, high int // low is the first index, high is the last index + 1, length = high - low
}

// New returns an initialized list with the values v added.
func New[T any](v ...T) *List[T] {
	l := new(List[T]).init()
	l.PushBack(v...)
	return l
}

const (
	defaultCapacity  = 128
	defaultInitIndex = defaultCapacity >> 1
)

// init initializes or clears list l.
func (l *List[T]) init() *List[T] {
	l.values = make([]T, defaultCapacity)
	l.low = defaultInitIndex
	l.high = defaultInitIndex
	return l
}

// insert inserts the values v to position i in list: prepend v to list if i equals
// to l.low; append v to list if i equals to l.high; insert v to list at position i
// if i is between l.low and l.high (exclusive); or do nothing if i is less than
// l.low or greater than l.high.
// At the same time, insert checks and expands the underlying array or moves list
// elements if necessary.
func (l *List[T]) insert(i int, v ...T) {
	if i < l.low || i > l.high {
		// invalid position i, do nothing
		return
	}
	switch i {
	case l.low:
		low := l.low - len(v)
		if low >= 0 {
			// prepend
			copy(l.values[low:l.low], v)
			l.low = low
			return
		}
	case l.high:
		high := l.high + len(v)
		if high <= cap(l.values) {
			// append
			copy(l.values[l.high:high], v)
			l.high = high
			return
		}
	default:
		tmp := cap(l.values) - l.high
		if tmp > l.low {
			// free space on high side is more than low side, try to insert to high side
			if len(v) <= tmp {
				// insert
				high := l.high + len(v)
				j := i + len(v)
				copy(l.values[j:high], l.values[i:l.high])
				copy(l.values[i:j], v)
				l.high = high
				return
			}
		} else {
			// free space on high side is not more than low side, try to insert to low side
			if len(v) <= l.low {
				// insert
				low := l.low - len(v)
				j := i - len(v)
				copy(l.values[low:j], l.values[l.low:i])
				copy(l.values[j:i], v)
				l.low = low
				return
			}
		}
	}
	// process exceeded cases
	s1 := l.Len()
	s2 := s1 + len(v)
	capacity := max(s2<<1, defaultCapacity)
	if capacity > cap(l.values) {
		// expand & migrate
		v1 := make([]T, capacity)
		low := (capacity - s2) >> 1
		high := low + s2
		switch i {
		case l.low:
			// prepend
			j := low + len(v)
			copy(v1[low:j], v)
			copy(v1[j:high], l.values[l.low:l.high])
		case l.high:
			// append
			j := low + s1
			copy(v1[low:j], l.values[l.low:l.high])
			copy(v1[j:high], v)
		default:
			// insert
			j := low + (i - l.low)
			copy(v1[low:j], l.values[l.low:i])
			k := j + len(v)
			copy(v1[j:k], v)
			copy(v1[k:high], l.values[i:l.high])
		}
		l.values = v1
		l.low = low
		l.high = high
	} else {
		// move
		low := (cap(l.values) - s2) >> 1
		high := low + s2
		switch i {
		case l.low:
			// prepend
			i := low + len(v)
			copy(l.values[i:high], l.values[l.low:l.high])
			copy(l.values[low:i], v)
		case l.high:
			// append
			i := low + s1
			copy(l.values[low:i], l.values[l.low:l.high])
			copy(l.values[i:high], v)
		default:
			// insert
			j := low + (i - l.low)
			copy(l.values[low:j], l.values[l.low:i])
			k := j + len(v)
			copy(l.values[k:high], l.values[i:l.high])
			copy(l.values[j:k], v)
		}
		l.low = low
		l.high = high
	}
}

// delete removes the value at position i if i is between l.low and l.high (exclusive), or do
// nothing.
// At the same time, delete checks and shrinks the underlying array if necessary.
func (l *List[T]) delete(i int) {
	if i < l.low || i >= l.high {
		return
	}
	if cap(l.values) > defaultCapacity {
		size := l.Len() - 1
		if size<<2 <= cap(l.values) {
			// shrink & migrate & delete
			capacity := max(size<<1, defaultCapacity)
			v := make([]T, capacity)
			low := (capacity - size) >> 1
			high := low + size
			switch i {
			case l.low:
				copy(v[low:high], l.values[i+1:l.high])
			case l.high - 1:
				copy(v[low:high], l.values[l.low:i])
			default:
				j := low + (i - l.low)
				copy(v[low:j], l.values[l.low:i])
				copy(v[j:high], l.values[i+1:l.high])
			}
			l.values = v
			l.low = low
			l.high = high
			return
		}
	}
	// delete
	switch i {
	case l.low:
		l.low++
	case l.high - 1:
		l.high--
	default:
		// delete: move the smaller part
		if l.high-i-1 > i-l.low {
			low := l.low + 1
			copy(l.values[low:i+1], l.values[l.low:i])
			l.low = low
		} else {
			high := l.high - 1
			copy(l.values[i:high], l.values[i+1:l.high])
			l.high = high
		}
	}
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List[T]) Len() int {
	return l.high - l.low
}

// Values returns a slice of all values of list.
func (l *List[T]) Values() []T {
	values := make([]T, l.Len())
	copy(values, l.values[l.low:l.high])
	return values
}

// String returns the string representation of list.
// Ref: std fmt.Stringer.
func (l *List[T]) String() string {
	values, _ := jsonx.MarshalToString(l.values[l.low:l.high])
	return "ArrayList: " + values
}

// MarshalJSON marshals list into valid JSON.
// Ref: std json.Marshaler.
func (l *List[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.values[l.low:l.high])
}

// UnmarshalJSON unmarshals a JSON description of list.
// The input can be assumed to be a valid encoding of a JSON value.
// UnmarshalJSON must copy the JSON data if it wishes to retain the data after returning.
// Ref: std json.Unmarshaler.
func (l *List[T]) UnmarshalJSON(data []byte) error {
	var v []T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	l.Clear()
	l.PushBack(v...)
	return nil
}

// Front returns the first value if exists in list.
// The ok result indicates whether such value was found in list.
func (l *List[T]) Front() (value T, ok bool) {
	if l.low < l.high {
		value = l.values[l.low]
		ok = true
	}
	return
}

// Back returns the last value if exists in list.
// The ok result indicates whether such value was found in list.
func (l *List[T]) Back() (value T, ok bool) {
	if l.low < l.high {
		value = l.values[l.high-1]
		ok = true
	}
	return
}

// PushFront inserts the given values v at the front of list.
func (l *List[T]) PushFront(v ...T) {
	l.insert(l.low, v...)
}

// PushBack inserts the given values v at the back of list.
func (l *List[T]) PushBack(v ...T) {
	l.insert(l.high, v...)
}

// PopFront removes the first value if exists in list and returns it.
// The ok result indicates whether such value was removed from list.
func (l *List[T]) PopFront() (value T, ok bool) {
	if l.low < l.high {
		value = l.values[l.low]
		l.low++
		ok = true
	}
	return
}

// PopBack removes the last value if exists in list and returns it.
// The ok result indicates whether such value was removed from list.
func (l *List[T]) PopBack() (value T, ok bool) {
	if l.low < l.high {
		l.high--
		value = l.values[l.high]
		ok = true
	}
	return
}

// Clear removes all values in list.
func (l *List[T]) Clear() {
	l.init()
}

// Get returns the value of index i if exists in list.
// The ok result indicates whether such value was found in list.
func (l *List[T]) Get(i int) (value T, ok bool) {
	i += l.low
	if i >= l.low && i < l.high {
		value = l.values[i]
		ok = true
	}
	return
}

// Set sets the value to v of index i if exists in list.
func (l *List[T]) Set(i int, v T) {
	i += l.low
	if i >= l.low && i < l.high {
		l.values[i] = v
	}
}

// Add inserts the values v to index i if exists in list, or appends the value v to the back of
// list if index i points to the next index of the last element in list.
func (l *List[T]) Add(i int, v ...T) {
	l.insert(i+l.low, v...)
}

// Del removes the value at index i if exists in list.
func (l *List[T]) Del(i int) {
	l.delete(i + l.low)
}

// Swap swaps the values with indices i and j if both indices exist in list.
func (l *List[T]) Swap(i, j int) {
	if i == j {
		return
	}
	i, j = i+l.low, j+l.low
	if i >= l.low && i < l.high && j >= l.low && j < l.high {
		l.values[i], l.values[j] = l.values[j], l.values[i]
	}
}

// Sort sorts list values (in-place) with the given cmp.
func (l *List[T]) Sort(cmp container.Compare[T]) {
	if cmp != nil && l.Len() > 1 {
		slices.SortFunc(l.values[l.low:l.high], cmp)
	}
}

// Range calls f sequentially for each index i and value v present in list.
// If f returns false, range stops the iteration.
func (l *List[T]) Range(f func(i int, v T) bool) {
	if f == nil {
		return
	}
	for i := l.low; i < l.high; i++ {
		if !f(i-l.low, l.values[i]) {
			break
		}
	}
}

// RRange calls f sequentially (in reverse order) for each index i and value v present in list.
// If f returns false, range stops the iteration.
func (l *List[T]) RRange(f func(i int, v T) bool) {
	if f == nil {
		return
	}
	for i := l.high - 1; i >= l.low; i-- {
		if !f(i-l.low, l.values[i]) {
			break
		}
	}
}

// LastIndex returns the index of the last occurrence of value v in list l,
// or -1 if not present.
func LastIndex[T comparable](l *List[T], v T) (index int) {
	index = -1
	if l == nil || l.Len() == 0 {
		return
	}
	l.RRange(func(i int, v1 T) bool {
		if v1 == v {
			index = i
			return false
		}
		return true
	})
	return
}

// Find returns the last index i and the corresponding value v in list l satisfying
// condition f(i, v), or first return parameter would be -1 if none do.
func FindLast[T any](l *List[T], f func(i int, v T) bool) (index int, value T) {
	index = -1
	if l == nil || l.Len() == 0 || f == nil {
		return
	}
	l.RRange(func(i int, v T) bool {
		if f(i, v) {
			index = i
			value = v
			return false
		}
		return true
	})
	return
}
