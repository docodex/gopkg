// Package tree provides an abstract Tree interface.
//
// In computer science, a tree is a widely used abstract data type (ADT) or data structure
// implementing this ADT that simulates a hierarchical tree structure, with a root value
// and subtrees of children with a parent node, represented as a set of linked nodes.
//
// Reference: https://en.wikipedia.org/wiki/Tree_%28data_structure%29
package tree

import "github.com/docodex/gopkg/container"

type Tree[T any] interface {
	container.Container[T]
}
