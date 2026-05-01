package deque_test

import (
	"testing"

	"github.com/docodex/gopkg/container/queue/deque"
	"github.com/stretchr/testify/assert"
)

func TestAliasesAndClear(t *testing.T) {
	q := deque.New[int]()
	// Dequeue/Peek on empty queue
	_, ok := q.Dequeue()
	assert.False(t, ok)
	_, ok = q.Peek()
	assert.False(t, ok)

	q.EnqueueBack(1)
	q.EnqueueBack(2)
	q.EnqueueBack(3)

	// Peek is alias for PeekFront.
	v, ok := q.Peek()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	// Dequeue is alias for DequeueFront.
	v, ok = q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	assert.Equal(t, 2, q.Len())

	// Clear resets the queue.
	q.Clear()
	assert.Equal(t, 0, q.Len())
	_, ok = q.Peek()
	assert.False(t, ok)
}

func TestInsertInvalidPosition(t *testing.T) {
	// Public API only exposes EnqueueFront/EnqueueBack (valid positions), so
	// the "invalid position, do nothing" branch of insert is currently
	// unreachable through the public API. Verified here via valid API only.
	q := deque.New[int]()
	q.EnqueueFront(1)
	q.EnqueueBack(2)
	q.EnqueueFront(0)
	assert.Equal(t, []int{0, 1, 2}, q.Values())
}

func TestInsertExpandPrepend(t *testing.T) {
	q := deque.New[int]()
	// EnqueueFront 200 items; uses one-at-a-time path, causing multiple expansions.
	for i := range 200 {
		q.EnqueueFront(i)
	}
	assert.Equal(t, 200, q.Len())
	v, ok := q.PeekFront()
	assert.True(t, ok)
	assert.Equal(t, 199, v)
	v, ok = q.PeekBack()
	assert.True(t, ok)
	assert.Equal(t, 0, v)
}

func TestInsertExpandAppend(t *testing.T) {
	q := deque.New[int]()
	for i := range 200 {
		q.EnqueueBack(i)
	}
	assert.Equal(t, 200, q.Len())
	v, _ := q.PeekFront()
	assert.Equal(t, 0, v)
	v, _ = q.PeekBack()
	assert.Equal(t, 199, v)
}

func TestDequeueShrinkPaths(t *testing.T) {
	// Grow beyond default cap, then dequeue from front until shrink triggers.
	q := deque.New[int]()
	for i := range 400 {
		q.EnqueueBack(i)
	}
	for i := range 350 {
		v, ok := q.DequeueFront()
		assert.True(t, ok)
		assert.Equal(t, i, v)
	}
	assert.Equal(t, 50, q.Len())
	// And from back.
	q = deque.New[int]()
	for i := range 400 {
		q.EnqueueBack(i)
	}
	for i := range 350 {
		v, ok := q.DequeueBack()
		assert.True(t, ok)
		assert.Equal(t, 399-i, v)
	}
	assert.Equal(t, 50, q.Len())
}

func TestCheckAndShrinkNoShrinkBranches(t *testing.T) {
	// cap <= defaultCapacity early return
	q := deque.New[int]()
	q.EnqueueBack(1)
	v, ok := q.DequeueFront()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	// Grow cap, then dequeue little so size<<2 > cap: no-shrink branch.
	q = deque.New[int]()
	for i := range 300 {
		q.EnqueueBack(i)
	}
	// Dequeue just a few from front; size stays above cap/4.
	for range 5 {
		q.DequeueFront()
	}
	assert.Equal(t, 295, q.Len())
}

func TestUnmarshalJSONErrors(t *testing.T) {
	q := deque.New[int]()
	assert.Error(t, q.UnmarshalJSON([]byte("not json")))

	// Valid roundtrip: ensures MarshalJSON/UnmarshalJSON happy path.
	for i := range 3 {
		q.EnqueueBack(i)
	}
	data, err := q.MarshalJSON()
	assert.NoError(t, err)
	q2 := deque.New[int]()
	assert.NoError(t, q2.UnmarshalJSON(data))
	assert.Equal(t, q.Values(), q2.Values())

	// Empty input.
	q3 := deque.New[int]()
	assert.NoError(t, q3.UnmarshalJSON([]byte(`[]`)))
	assert.Equal(t, 0, q3.Len())
}

func TestDequeBackWithEmpty(t *testing.T) {
	q := deque.New[int]()
	_, ok := q.DequeueBack()
	assert.False(t, ok)
	_, ok = q.PeekBack()
	assert.False(t, ok)
}
