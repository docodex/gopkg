package binaryheap_test

import (
	"fmt"

	"github.com/docodex/gopkg/container/tree/binaryheap"
)

func ExampleHeap_Push() {
	h := binaryheap.New(2, 1, 5)
	h.Push(3)
	v, _ := h.Peek()
	fmt.Printf("minimum: %d\n", v)
	for h.Len() > 0 {
		v, _ := h.Pop()
		fmt.Printf("%d ", v)
	}

	// Output:
	// minimum: 1
	// 1 2 3 5
}

func ExampleHeap_Update() {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// An Item is something we manage in a priority queue.
	type Item struct {
		value    string // The value of the item.
		priority int    // The priority of the item in queue.
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make([]*Item, 0, len(items))
	for value, priority := range items {
		pq = append(pq, &Item{
			value:    value,
			priority: priority,
		})
	}

	h := binaryheap.NewFunc(func(a, b *Item) bool {
		// We want Pop to give us the highest, not lowest, priority so we use greater than here.
		return a.priority > b.priority
	}, pq...)

	// Insert a new item and then modify its priority.
	h.Push(&Item{
		value:    "orange",
		priority: 5,
	})

	item, _ := h.Peek()
	h.Update(0, &Item{
		value:    item.value,
		priority: 1,
	})

	// Take the items out; they arrive in decreasing priority order.
	for h.Len() > 0 {
		item, _ := h.Pop()
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}

	// Output:
	// 04:pear 03:banana 02:apple 01:orange
}
