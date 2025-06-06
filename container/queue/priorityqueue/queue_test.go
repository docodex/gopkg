package priorityqueue_test

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/queue/priorityqueue"
	"github.com/stretchr/testify/assert"
)

func verify(t *testing.T, q *priorityqueue.Queue[int], i int, less func(a, b int) bool) {
	t.Helper()
	elements := q.Elements()
	j1 := i<<1 + 1 //left child
	j2 := i<<1 + 2 //right child
	if j1 < len(elements) {
		if less(elements[j1].Value, elements[i].Value) {
			t.Errorf("queue invariant invalidated [%d] = %d > [%d] = %d",
				i, elements[i].Value, j1, elements[j1].Value)
			return
		}
		verify(t, q, j1, less)
	}
	if j2 < len(elements) {
		if less(elements[j2].Value, elements[i].Value) {
			t.Errorf("queue invariant invalidated [%d] = %d > [%d] = %d",
				i, elements[i].Value, j1, elements[j2].Value)
			return
		}
		verify(t, q, j2, less)
	}
}

func TestUn(t *testing.T) {
	var values []int
	for range 20 {
		values = append(values, 0) // all elements are different
	}
	v, err := json.Marshal(values)
	assert.Nil(t, err)
	less := func(a, b int) bool {
		return a > b
	}
	q := priorityqueue.NewFunc(less)
	err = q.UnmarshalJSON(v)
	assert.Nil(t, err)

	verify(t, q, 0, less)

	for i := 1; q.Len() > 0; i++ {
		x, _ := q.Dequeue()
		verify(t, q, 0, less)
		if x != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, x, 0)
		}
	}
}

func TestInit1(t *testing.T) {
	var values []int
	for i := range 20 {
		values = append(values, i) // all elements are different
	}
	v, err := json.Marshal(values)
	assert.Nil(t, err)
	less := func(a, b int) bool {
		return a > b
	}
	q := priorityqueue.NewFunc(less)
	err = q.UnmarshalJSON(v)
	assert.Nil(t, err)

	verify(t, q, 0, less)

	for i := 19; q.Len() > 0; i-- {
		x, _ := q.Dequeue()
		verify(t, q, 0, less)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestQueue(t *testing.T) {
	less := func(a, b int) bool {
		return a < b
	}
	q := priorityqueue.NewFunc(less)
	verify(t, q, 0, less)

	var values []int
	for i := 20; i > 10; i-- {
		values = append(values, i)
	}
	v, err := json.Marshal(values)
	assert.Nil(t, err)
	err = q.UnmarshalJSON(v)
	assert.Nil(t, err)
	verify(t, q, 0, less)

	for i := 10; i > 0; i-- {
		q.Enqueue(i)
		verify(t, q, 0, less)
	}

	for i := 1; q.Len() > 0; i++ {
		x, _ := q.Dequeue()
		if i < 20 {
			q.Enqueue(20 + i)
		}
		verify(t, q, 0, less)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestRemove0(t *testing.T) {
	var values []int
	for i := range 10 {
		values = append(values, i)
	}
	less := func(a, b int) bool {
		return a < b
	}
	q := priorityqueue.NewFunc(less, values...)
	verify(t, q, 0, less)

	for q.Len() > 0 {
		i := q.Len() - 1
		x, _ := q.Remove(i)
		if x != i {
			t.Errorf("Remove(%d) got %d; want %d", i, x, i)
		}
		verify(t, q, 0, less)
	}
}

func TestRemove1(t *testing.T) {
	var values []int
	for i := range 10 {
		values = append(values, i)
	}
	less := func(a, b int) bool {
		return a < b
	}
	q := priorityqueue.NewFunc(less, values...)
	verify(t, q, 0, less)

	for i := 0; q.Len() > 0; i++ {
		x, _ := q.Remove(0)
		if x != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}
		verify(t, q, 0, less)
	}
}

func TestRemove2(t *testing.T) {
	N := 10

	var values []int
	for i := range N {
		values = append(values, i)
	}
	less := func(a, b int) bool {
		return a < b
	}
	q := priorityqueue.NewFunc(less, values...)
	verify(t, q, 0, less)

	m := make(map[int]bool)
	for q.Len() > 0 {
		v, _ := q.Remove((q.Len() - 1) / 2)
		m[v] = true
		verify(t, q, 0, less)
	}

	if len(m) != N {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}
	for i := range len(m) {
		if !m[i] {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

func TestFix(t *testing.T) {
	less := func(a, b int) bool {
		return a < b
	}
	q := priorityqueue.NewFunc(less)
	verify(t, q, 0, less)

	for i := 200; i > 0; i -= 10 {
		q.Enqueue(i)
	}
	verify(t, q, 0, less)

	elements := q.Elements()
	if elements[0].Value != 10 {
		t.Fatalf("Expected head to be 10, was %d", elements[0].Value)
	}
	elements[0].Value = 210
	q.Fix(0)
	verify(t, q, 0, less)

	for i := 100; i > 0; i-- {
		elements := q.Elements()
		elem := rand.IntN(q.Len())
		if i&1 == 0 {
			elements[elem].Value *= 2
		} else {
			elements[elem].Value /= 2
		}
		q.Fix(elem)
		verify(t, q, 0, less)
	}
}

func TestUpdate(t *testing.T) {
	less := func(a, b int) bool {
		return a < b
	}
	q := priorityqueue.NewFunc(less)
	verify(t, q, 0, less)

	for i := 200; i > 0; i -= 10 {
		q.Enqueue(i)
	}
	verify(t, q, 0, less)

	elements := q.Elements()
	if elements[0].Value != 10 {
		t.Fatalf("Expected head to be 10, was %d", elements[0].Value)
	}
	q.Update(0, 210)
	verify(t, q, 0, less)

	for i := 100; i > 0; i-- {
		elements := q.Elements()
		elem := rand.IntN(q.Len())
		if i&1 == 0 {
			q.Update(elem, elements[elem].Value*2)
		} else {
			q.Update(elem, elements[elem].Value/2)
		}
		verify(t, q, 0, less)
	}
}

func TestPriorityQueue_Enqueue(t *testing.T) {
	t.Parallel()

	q := priorityqueue.New[int]()
	assert.Equal(t, true, q.Len() == 0)
	q.Enqueue(3)
	q.Enqueue(1)
	q.Enqueue(2)
	v := q.Values()
	assert.Equal(t, []int{1, 2, 3}, v)
	assert.True(t, q.Len() != 0)
}

func TestPriorityQueue_Dequeue(t *testing.T) {
	t.Parallel()

	q := priorityqueue.New[int]()
	_, ok := q.Dequeue()
	assert.Equal(t, false, ok)
	q.Enqueue(3)
	q.Enqueue(1)
	q.Enqueue(2)
	assert.Equal(t, 3, q.Len())
	val, ok := q.Dequeue()
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, val)
}

type Item struct {
	priority int
	name     string
}

func (element Item) String() string {
	return fmt.Sprintf("{%v %v}", element.priority, element.name)
}

// Comparator function (sort by priority value in descending order)
func byPriority(a, b Item) bool {
	return a.priority > b.priority
}

func TestBinaryQueueEnqueue(t *testing.T) {
	q := priorityqueue.NewFunc(byPriority)

	if actualValue := q.Len() == 0; actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	a := Item{name: "a", priority: 1}
	c := Item{name: "c", priority: 3}
	b := Item{name: "b", priority: 2}

	q.Enqueue(a)
	q.Enqueue(c)
	q.Enqueue(b)

	if actualValue := q.Values(); actualValue[0].name != "c" || actualValue[1].name != "b" || actualValue[2].name != "a" {
		t.Errorf("Got %v expected %v", actualValue, `[{3 c} {2 b} {1 a}]`)
	}

	count := 0
	for q.Len() != 0 {
		value, _ := q.Dequeue()
		switch count {
		case 0:
			if actualValue, expectedValue := value.name, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value.name, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value.name, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
		count++
	}
}

func TestBinaryQueueEnqueueBulk(t *testing.T) {
	q := priorityqueue.New[int]()

	q.Enqueue(15)
	q.Enqueue(20)
	q.Enqueue(3)
	q.Enqueue(1)
	q.Enqueue(2)

	if actualValue, ok := q.Dequeue(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
	if actualValue, ok := q.Dequeue(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := q.Dequeue(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := q.Dequeue(); actualValue != 15 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 15)
	}
	if actualValue, ok := q.Dequeue(); actualValue != 20 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 20)
	}

	q.Clear()
	if actualValue := q.Len() == 0; !actualValue {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestBinaryQueueDequeue(t *testing.T) {
	q := priorityqueue.New[int]()

	if actualValue := q.Len() == 0; actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	q.Enqueue(3)
	q.Enqueue(2)
	q.Enqueue(1)
	q.Dequeue() // removes 1

	if actualValue, ok := q.Dequeue(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := q.Dequeue(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := q.Dequeue(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := q.Len() == 0; actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := q.Values(); len(actualValue) != 0 {
		t.Errorf("Got %v expected %v", actualValue, "[]")
	}
}

func TestBinaryQueueRandom(t *testing.T) {
	q := priorityqueue.New[int]()

	for range 10000 {
		r := int(rand.Int32N(30))
		q.Enqueue(r)
	}

	prev, _ := q.Dequeue()
	for q.Len() != 0 {
		curr, _ := q.Dequeue()
		if prev > curr {
			t.Errorf("Queue property invalidated. prev: %v current: %v", prev, curr)
		}
		prev = curr
	}
}

func TestBinaryQueueSerialization(t *testing.T) {
	q := priorityqueue.New[string]()

	q.Enqueue("c")
	q.Enqueue("b")
	q.Enqueue("a")

	var err error
	assert := func() {
		if actualValue := q.Values(); actualValue[0] != "a" || actualValue[1] != "b" || actualValue[2] != "c" {
			t.Errorf("Got %v expected %v", actualValue, "[1,3,2]")
		}
		if actualValue := q.Len(); actualValue != 3 {
			t.Errorf("Got %v expected %v", actualValue, 3)
		}
		if actualValue, ok := q.Peek(); actualValue != "a" || !ok {
			t.Errorf("Got %v expected %v", actualValue, "a")
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

func TestBTreeString(t *testing.T) {
	q := priorityqueue.New[int]()
	q.Enqueue(1)
	if !strings.HasPrefix(q.String(), "PriorityQueue") {
		t.Errorf("String should start with container name")
	}
}

func benchmarkEnqueue(b *testing.B, q *priorityqueue.Queue[Item], size int) {
	for b.Loop() {
		for range size {
			q.Enqueue(Item{})
		}
	}
}

func benchmarkDequeue(b *testing.B, q *priorityqueue.Queue[Item], size int) {
	for b.Loop() {
		for range size {
			q.Dequeue()
		}
	}
}

func BenchmarkBinaryQueueDequeue100(b *testing.B) {
	b.StopTimer()
	size := 100
	q := priorityqueue.NewFunc(byPriority)
	for range size {
		q.Enqueue(Item{})
	}
	b.StartTimer()
	benchmarkDequeue(b, q, size)
}

func BenchmarkBinaryQueueDequeue1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	q := priorityqueue.NewFunc(byPriority)
	for range size {
		q.Enqueue(Item{})
	}
	b.StartTimer()
	benchmarkDequeue(b, q, size)
}

func BenchmarkBinaryQueueDequeue10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	q := priorityqueue.NewFunc(byPriority)
	for range size {
		q.Enqueue(Item{})
	}
	b.StartTimer()
	benchmarkDequeue(b, q, size)
}

func BenchmarkBinaryQueueDequeue100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	q := priorityqueue.NewFunc(byPriority)
	for range size {
		q.Enqueue(Item{})
	}
	b.StartTimer()
	benchmarkDequeue(b, q, size)
}

func BenchmarkBinaryQueueEnqueue100(b *testing.B) {
	b.StopTimer()
	size := 100
	q := priorityqueue.NewFunc(byPriority)
	b.StartTimer()
	benchmarkEnqueue(b, q, size)
}

func BenchmarkBinaryQueueEnqueue1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	q := priorityqueue.NewFunc(byPriority)
	b.StartTimer()
	benchmarkEnqueue(b, q, size)
}

func BenchmarkBinaryQueueEnqueue10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	q := priorityqueue.NewFunc(byPriority)
	b.StartTimer()
	benchmarkEnqueue(b, q, size)
}

func BenchmarkBinaryQueueEnqueue100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	q := priorityqueue.NewFunc(byPriority)
	b.StartTimer()
	benchmarkEnqueue(b, q, size)
}
