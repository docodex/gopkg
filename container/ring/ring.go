// Package ring provides an abstract Ring interface.
package ring

import "github.com/docodex/gopkg/container"

type Ring[T any] interface {
	container.Container[T]

	// Range calls f sequentially for each value v present in ring.
	// If f returns false, range stops the iteration.
	Range(f func(v T) bool)
}
