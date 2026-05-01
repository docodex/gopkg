package doublylinkedlist_test

import (
	"fmt"

	"github.com/docodex/gopkg/container/list/doublylinkedlist"
)

func ExampleList() {
	// Create a new list and put some numbers in it.
	l := doublylinkedlist.New[int]()
	l.PushBack(4)
	e4 := l.BackNode()
	l.PushFront(1)
	e1 := l.FrontNode()
	l.InsertBefore(e4, 3)
	l.InsertAfter(e1, 2)

	// Iterate through list and print its contents.
	for x := l.FrontNode(); x != nil; x = x.Next() {
		fmt.Println(x.Value)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
}

func ExampleList_PushFront() {
	l := doublylinkedlist.New[int]()
	l.PushFront(4, 3, 2, 1)

	// Iterate through list and print its contents.
	for x := l.FrontNode(); x != nil; x = x.Next() {
		fmt.Println(x.Value)
	}

	// Output:
	// 4
	// 3
	// 2
	// 1
}

func ExampleList_PushBack() {
	l := doublylinkedlist.New[int]()
	l.PushFront(4, 3, 2, 1)
	l.PushBack(5, 6, 7, 8)

	// Iterate through list and print its contents.
	for x := l.FrontNode(); x != nil; x = x.Next() {
		fmt.Println(x.Value)
	}

	// Output:
	// 4
	// 3
	// 2
	// 1
	// 5
	// 6
	// 7
	// 8
}
