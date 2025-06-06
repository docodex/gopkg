package arrayqueue_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/queue/arrayqueue"
)

func TestQueueEnqueue(t *testing.T) {
	q := arrayqueue.New[any]()
	if actualValue := (q.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	if actualValue := q.Values(); actualValue[0].(int) != 1 || actualValue[1].(int) != 2 || actualValue[2].(int) != 3 {
		t.Errorf("Got %v expected %v", actualValue, "[1,2,3]")
	}
	if actualValue := (q.Len() == 0); actualValue != false {
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
	q := arrayqueue.New[any]()
	if actualValue, ok := q.Peek(); actualValue != nil || ok {
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
	q := arrayqueue.New[any]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Dequeue()
	if actualValue, ok := q.Peek(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := q.Dequeue(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := q.Dequeue(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := q.Dequeue(); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := (q.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := q.Values(); len(actualValue) != 0 {
		t.Errorf("Got %v expected %v", actualValue, "[]")
	}
}

func TestQueueSerialization(t *testing.T) {
	q := arrayqueue.New[any]()
	q.Enqueue("a")
	q.Enqueue("b")
	q.Enqueue("c")

	var err error
	assert := func() {
		if actualValue, expectedValue := fmt.Sprintf("%s%s%s", q.Values()...), "abc"; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
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

	err = json.Unmarshal([]byte(`[1,2,3]`), &q)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
}

func TestQueueString(t *testing.T) {
	q := arrayqueue.New[any]()
	q.Enqueue(1)
	if !strings.HasPrefix(q.String(), "ArrayQueue") {
		t.Errorf("String should start with container name")
	}
}

func benchmarkEnqueue(b *testing.B, q *arrayqueue.Queue[any], size int) {
	for b.Loop() {
		for n := range size {
			q.Enqueue(n)
		}
	}
}

func benchmarkDequeue(b *testing.B, q *arrayqueue.Queue[any], size int) {
	for b.Loop() {
		for range size {
			q.Dequeue()
		}
	}
}

func BenchmarkArrayQueueDequeue100(b *testing.B) {
	b.StopTimer()
	size := 100
	q := arrayqueue.New[any]()
	for n := range size {
		q.Enqueue(n)
	}
	b.StartTimer()
	benchmarkDequeue(b, q, size)
}

func BenchmarkArrayQueueDequeue1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	q := arrayqueue.New[any]()
	for n := range size {
		q.Enqueue(n)
	}
	b.StartTimer()
	benchmarkDequeue(b, q, size)
}

func BenchmarkArrayQueueDequeue10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	q := arrayqueue.New[any]()
	for n := range size {
		q.Enqueue(n)
	}
	b.StartTimer()
	benchmarkDequeue(b, q, size)
}

func BenchmarkArrayQueueDequeue100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	q := arrayqueue.New[any]()
	for n := range size {
		q.Enqueue(n)
	}
	b.StartTimer()
	benchmarkDequeue(b, q, size)
}

func BenchmarkArrayQueueEnqueue100(b *testing.B) {
	b.StopTimer()
	size := 100
	q := arrayqueue.New[any]()
	b.StartTimer()
	benchmarkEnqueue(b, q, size)
}

func BenchmarkArrayQueueEnqueue1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	q := arrayqueue.New[any]()
	b.StartTimer()
	benchmarkEnqueue(b, q, size)
}

func BenchmarkArrayQueueEnqueue10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	q := arrayqueue.New[any]()
	b.StartTimer()
	benchmarkEnqueue(b, q, size)
}

func BenchmarkArrayQueueEnqueue100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	q := arrayqueue.New[any]()
	b.StartTimer()
	benchmarkEnqueue(b, q, size)
}
