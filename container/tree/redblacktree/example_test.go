package redblacktree_test

import (
	"fmt"

	"github.com/docodex/gopkg/container/tree/redblacktree"
)

func ExampleTree() {
	strs := []string{"Hello", "World", "Golang", "Python", "Rust", "C", "JavaScript", "Haskell", "Pascal", "ZZ"}
	ints := []int{3, 5, 4, 1, 8, 6, 5, 7, 9, 0}
	t1 := redblacktree.New[string, int]()
	t2 := redblacktree.New[int, string]()
	for i := range len(strs) {
		t1.Insert(strs[i], ints[i])
		t2.Insert(ints[i], strs[i])
	}
	k1 := t1.Keys()
	v1 := t1.Values()
	for i := range k1 {
		fmt.Printf("%s:%d\n", k1[i], v1[i])
	}
	k2 := t2.Keys()
	v2 := t2.Values()
	for i := range k2 {
		fmt.Printf("%d:%s\n", k2[i], v2[i])
	}

	// Output:
	// C:6
	// Golang:4
	// Haskell:7
	// Hello:3
	// JavaScript:5
	// Pascal:9
	// Python:1
	// Rust:8
	// World:5
	// ZZ:0
	// 0:ZZ
	// 1:Python
	// 3:Hello
	// 4:Golang
	// 5:JavaScript
	// 6:C
	// 7:Haskell
	// 8:Rust
	// 9:Pascal
}
