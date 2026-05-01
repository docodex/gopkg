// Package ring provides an abstract Ring interface.
//
// Concurrency: implementations are NOT safe for concurrent use. Callers must
// provide external synchronization if a ring is shared across goroutines.
// Unlike other containers, ring types (doublylinkedring, singlylinkedring)
// support zero-value use and lazily self-initialize.
// See package container for the project-wide concurrency policy.
package ring

import "github.com/docodex/gopkg/container"

type Ring[T any] interface {
	container.Container[T]

	// Range calls f sequentially for each value v present in ring.
	// If f returns false, range stops the iteration.
	Range(f func(v T) bool)
}
