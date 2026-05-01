package priorityqueue_test

import (
	"fmt"

	"github.com/docodex/gopkg/container/queue/priorityqueue"
)

func ExampleQueue_Enqueue() {
	h := priorityqueue.NewFunc(func(a, b int) bool {
		return a > b
	})
	h.Enqueue(3)
	h.Enqueue(5)
	h.Enqueue(1)
	h.Enqueue(2)
	v, _ := h.Peek()
	fmt.Printf("maximum: %d\n", v)
	for h.Len() > 0 {
		v, _ := h.Dequeue()
		fmt.Printf("%d ", v)
	}

	// Output:
	// maximum: 5
	// 5 3 2 1
}

func ExampleQueue_Update() {
	// Some items and their priorities.
	items := map[string]int{
		"apple": 2, "pear": 4,
	}

	// An Item is something we manage in a priority queue.
	type Item struct {
		value    string // The value of the item; arbitrary.
		priority int    // The priority of the item in queue.
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	q := priorityqueue.NewFunc(func(a, b *Item) bool {
		return a.priority > b.priority
	})
	for value, priority := range items {
		q.Enqueue(&Item{
			value:    value,
			priority: priority,
		})
	}

	// Insert a new item and then modify its priority.
	orange := q.Enqueue(&Item{
		value:    "orange",
		priority: 1,
	})
	orange.Value.priority = 5
	q.Fix(orange.Index())

	// Insert a new item and then update its priority.
	banana := q.Enqueue(&Item{
		value:    "banana",
		priority: 3,
	})
	q.Update(banana.Index(), &Item{
		value:    "banana",
		priority: 1,
	})

	// Take the items out; they arrive in decreasing priority order.
	for q.Len() > 0 {
		item, _ := q.Dequeue()
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
	// Output:
	// 05:orange 04:pear 02:apple 01:banana
}
