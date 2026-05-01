package doublylinkedring_test

import (
	"fmt"

	"github.com/docodex/gopkg/container/ring/doublylinkedring"
)

func ExampleRing_Len() {
	// Create a new ring of size 4
	r := doublylinkedring.New(1, 2, 3, 4)

	// Print out its length
	fmt.Println(r.Len())

	// Output:
	// 4
}

func ExampleRing_Next() {
	// Create a new ring of size 5
	r := doublylinkedring.New(0, 1, 2, 3, 4)

	// Get the length of the ring
	n := r.Len()

	// Iterate through the ring and print its contents
	for range n {
		fmt.Println(r.Value)
		r = r.Next()
	}

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleRing_Prev() {
	// Create a new ring of size 5
	r := doublylinkedring.New(0, 1, 2, 3, 4)

	// Get the length of the ring
	n := r.Len()

	// Iterate through the ring backwards and print its contents
	for range n {
		r = r.Prev()
		fmt.Println(r.Value)
	}

	// Output:
	// 4
	// 3
	// 2
	// 1
	// 0
}

func ExampleRing_Range() {
	// Create a new ring of size 5
	r := doublylinkedring.New(0, 1, 2, 3, 4)

	// Iterate through the ring and print its contents
	r.Range(func(value int) bool {
		fmt.Println(value)
		return true
	})

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleRing_Move() {
	// Create a new ring of size 5
	r := doublylinkedring.New(0, 1, 2, 3, 4)

	// Get the length of the ring
	n := r.Len()

	// Initialize the ring with some integer values
	for i := range n {
		r.Value = i
		r = r.Next()
	}

	// Move the pointer forward by three steps
	r = r.Move(3)

	// Iterate through the ring and print its contents
	r.Range(func(value int) bool {
		fmt.Println(value)
		return true
	})

	// Output:
	// 3
	// 4
	// 0
	// 1
	// 2
}

func ExampleRing_Link() {
	// Create two rings, r and s, of size 2
	r := doublylinkedring.New(0, 0)
	s := doublylinkedring.New(1, 1)

	// Get the length of the ring
	lr := r.Len()
	ls := s.Len()

	// Initialize r with 0s
	for range lr {
		r.Value = 0
		r = r.Next()
	}

	// Initialize s with 1s
	for range ls {
		s.Value = 1
		s = s.Next()
	}

	// Link ring r and ring s
	rs := r.Link(s)

	// Iterate through the combined ring and print its contents
	rs.Range(func(value int) bool {
		fmt.Println(value)
		return true
	})

	// Output:
	// 0
	// 0
	// 1
	// 1
}

func ExampleRing_Unlink() {
	// Create a new ring of size 6
	r := doublylinkedring.New(0, 1, 2, 3, 4, 5)

	// Unlink three elements from r, starting from r.Next()
	r.Unlink(3)

	// Iterate through the remaining ring and print its contents
	r.Range(func(value int) bool {
		fmt.Println(value)
		return true
	})

	// Output:
	// 0
	// 4
	// 5
}
