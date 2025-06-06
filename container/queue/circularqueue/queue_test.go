package circularqueue_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/queue/circularqueue"
)

func TestQueueEnqueue(t *testing.T) {
	q := circularqueue.New[int](3)
	if actualValue := q.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	if actualValue := q.Values(); actualValue[0] != 1 || actualValue[1] != 2 || actualValue[2] != 3 {
		t.Errorf("Got %v expected %v", actualValue, "[1,2,3]")
	}
	if actualValue := q.Empty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := q.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := q.Peek(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestQueuePeek(t *testing.T) {
	q := circularqueue.New[int](3)
	if actualValue, ok := q.Peek(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	if actualValue, ok := q.Peek(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestQueueDequeue(t *testing.T) {
	assert := func(actualValue any, expectedValue any) {
		if actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	q := circularqueue.New[int](3)
	assert(q.Empty(), true)
	assert(q.Empty(), true)
	assert(q.Full(), false)
	assert(q.Len(), 0)
	q.Enqueue(1)
	assert(q.Len(), 1)
	q.Enqueue(2)
	assert(q.Len(), 2)

	q.Enqueue(3)
	assert(q.Len(), 3)
	assert(q.Empty(), false)
	assert(q.Full(), true)

	q.Dequeue()
	assert(q.Len(), 2)

	if actualValue, ok := q.Peek(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	assert(q.Len(), 2)

	if actualValue, ok := q.Dequeue(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	assert(q.Len(), 1)

	if actualValue, ok := q.Dequeue(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	assert(q.Len(), 0)
	assert(q.Empty(), true)
	assert(q.Full(), false)

	if actualValue, ok := q.Dequeue(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	assert(q.Len(), 0)

	assert(q.Empty(), true)
	assert(q.Full(), false)
	assert(len(q.Values()), 0)
}

func TestQueueDequeueFull(t *testing.T) {
	assert := func(actualValue any, expectedValue any) {
		if actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}

	q := circularqueue.New[int](2)
	assert(q.Empty(), true)
	assert(q.Full(), false)
	assert(q.Len(), 0)

	ok := q.Enqueue(1)
	assert(ok, true)
	assert(q.Len(), 1)

	ok = q.Enqueue(2)
	assert(ok, true)
	assert(q.Len(), 2)
	assert(q.Full(), true)
	if actualValue, ok := q.Peek(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	ok = q.Enqueue(3) // overwrites 1
	assert(ok, false)
	assert(q.Len(), 2)

	if actualValue, ok := q.Dequeue(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
	if actualValue, expectedValue := q.Len(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, ok := q.Peek(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, expectedValue := q.Len(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, ok := q.Dequeue(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	assert(q.Len(), 0)

	if actualValue, ok := q.Dequeue(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	assert(q.Empty(), true)
	assert(q.Full(), false)
	assert(len(q.Values()), 0)
}

func SameElements[T comparable](t *testing.T, actual, expected []T) {
	if len(actual) != len(expected) {
		t.Errorf("Got %d expected %d", len(actual), len(expected))
	}
outer:
	for _, e := range expected {
		for _, a := range actual {
			if e == a {
				continue outer
			}
		}
		t.Errorf("Did not find expected element %v in %v", e, actual)
	}
}

func TestQueueSerialization(t *testing.T) {
	q := circularqueue.New[string](3)
	q.Enqueue("a")
	q.Enqueue("b")
	q.Enqueue("c")

	var err error
	assert := func() {
		SameElements(t, q.Values(), []string{"a", "b", "c"})
		if actualValue, expectedValue := q.Len(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	bytes, err := q.MarshalJSON()
	assert()

	err = q.UnmarshalJSON(bytes)
	assert()

	bytes, err = json.Marshal([]any{"a", "b", "c", q})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	err = json.Unmarshal([]byte(`["a","b","c"]`), &q)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
	assert()
}

func TestQueueString(t *testing.T) {
	c := circularqueue.New[int](3)
	c.Enqueue(1)
	if !strings.HasPrefix(c.String(), "CircularQueue") {
		t.Errorf("String should start with container name")
	}
}

func benchmarkEnqueue(b *testing.B, q *circularqueue.Queue[int], size int) {
	for b.Loop() {
		for n := range size {
			q.Enqueue(n)
		}
	}
}

func benchmarkDequeue(b *testing.B, q *circularqueue.Queue[int], size int) {
	for b.Loop() {
		for range size {
			q.Dequeue()
		}
	}
}

func BenchmarkCircularQueueDequeue100(b *testing.B) {
	b.StopTimer()
	size := 100
	q := circularqueue.New[int](3)
	for n := range size {
		q.Enqueue(n)
	}
	b.StartTimer()
	benchmarkDequeue(b, q, size)
}

func BenchmarkCircularQueueDequeue1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	q := circularqueue.New[int](3)
	for n := range size {
		q.Enqueue(n)
	}
	b.StartTimer()
	benchmarkDequeue(b, q, size)
}

func BenchmarkCircularQueueDequeue10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	q := circularqueue.New[int](3)
	for n := range size {
		q.Enqueue(n)
	}
	b.StartTimer()
	benchmarkDequeue(b, q, size)
}

func BenchmarkCircularQueueDequeue100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	q := circularqueue.New[int](3)
	for n := range size {
		q.Enqueue(n)
	}
	b.StartTimer()
	benchmarkDequeue(b, q, size)
}

func BenchmarkCircularQueueEnqueue100(b *testing.B) {
	b.StopTimer()
	size := 100
	queue := circularqueue.New[int](3)
	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkCircularQueueEnqueue1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	q := circularqueue.New[int](3)
	b.StartTimer()
	benchmarkEnqueue(b, q, size)
}

func BenchmarkCircularQueueEnqueue10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	q := circularqueue.New[int](3)
	b.StartTimer()
	benchmarkEnqueue(b, q, size)
}

func BenchmarkCircularQueueEnqueue100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	q := circularqueue.New[int](3)
	b.StartTimer()
	benchmarkEnqueue(b, q, size)
}
