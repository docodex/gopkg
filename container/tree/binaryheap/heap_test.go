package binaryheap

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func verify(t *testing.T, h *Heap[int], i int) {
	t.Helper()
	n := h.Len()
	j1 := h.left(i)
	j2 := h.right(i)
	if j1 < n {
		if h.less(h.values[j1], h.values[i]) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d",
				i, h.values[i], j1, h.values[j1])
			return
		}
		verify(t, h, j1)
	}
	if j2 < n {
		if h.less(h.values[j2], h.values[i]) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d",
				i, h.values[i], j1, h.values[j2])
			return
		}
		verify(t, h, j2)
	}
}

func TestInit0(t *testing.T) {
	var values []int
	for i := 20; i > 0; i-- {
		values = append(values, 0) // all elements are the same
	}
	h := New(values...)
	verify(t, h, 0)

	for i := 1; h.Len() > 0; i++ {
		x, _ := h.Pop()
		verify(t, h, 0)
		if x != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, x, 0)
		}
	}
}

func TestInit1(t *testing.T) {
	var values []int
	for i := 20; i > 0; i-- {
		values = append(values, i) // all elements are different
	}
	h := New(values...)
	verify(t, h, 0)

	for i := 1; h.Len() > 0; i++ {
		x, _ := h.Pop()
		verify(t, h, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestInit2(t *testing.T) {
	var values []int
	for i := 20; i > 0; i-- {
		values = append(values, i) // all elements are different
	}
	v, err := json.Marshal(values)
	assert.Nil(t, err)
	h := New[int]()
	err = h.UnmarshalJSON(v)
	assert.Nil(t, err)

	verify(t, h, 0)

	for i := 1; h.Len() > 0; i++ {
		x, _ := h.Pop()
		verify(t, h, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestHeap(t *testing.T) {
	h := New[int]()
	verify(t, h, 0)

	for i := 20; i > 10; i-- {
		h.values = append(h.values, i)
	}
	h.init()
	verify(t, h, 0)

	for i := 10; i > 0; i-- {
		h.Push(i)
		verify(t, h, 0)
	}

	for i := 1; h.Len() > 0; i++ {
		x, _ := h.Pop()
		if i < 20 {
			h.Push(20 + i)
		}
		verify(t, h, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestRemove0(t *testing.T) {
	h := New[int]()
	for i := range 10 {
		h.values = append(h.values, i)
	}
	h.init()
	verify(t, h, 0)

	for h.Len() > 0 {
		i := h.Len() - 1
		x, _ := h.Remove(i)
		if x != i {
			t.Errorf("Remove(%d) got %d; want %d", i, x, i)
		}
		verify(t, h, 0)
	}
}

func TestRemove1(t *testing.T) {
	h := New[int]()
	for i := range 10 {
		h.values = append(h.values, i)
	}
	h.init()
	verify(t, h, 0)

	for i := 0; h.Len() > 0; i++ {
		x, _ := h.Remove(0)
		if x != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}
		verify(t, h, 0)
	}
}

func TestRemove2(t *testing.T) {
	N := 10

	h := New[int]()
	for i := range N {
		h.values = append(h.values, i)
	}
	h.init()
	verify(t, h, 0)

	m := make(map[int]bool)
	for h.Len() > 0 {
		x, _ := h.Remove((h.Len() - 1) / 2)
		m[x] = true
		verify(t, h, 0)
	}

	if len(m) != N {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}
	for i := 0; i < len(m); i++ {
		if !m[i] {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

type element struct {
	value    string
	priority int
}

func (e *element) String() string {
	return fmt.Sprintf("{%s:%d}", e.value, e.priority)
}

func TestElements(t *testing.T) {
	h := NewFunc(
		func(a, b *element) bool {
			return a.priority > b.priority
		},
		&element{value: "b", priority: 2},
		&element{value: "d", priority: 4},
		&element{value: "c", priority: 3},
		&element{value: "a", priority: 1},
	)
	fmt.Println(h.Values())
	fmt.Println(h.Elements())
	es := h.Elements()
	for i, e := range es {
		e.priority = rand.IntN(10)
		h.Update(i, e)
	}
	fmt.Println(h.Values())
	fmt.Println(h.Elements())
	h.Clear()
	assert.True(t, h.Len() == 0)
	fmt.Println(h.Values())
	fmt.Println(h.Elements())
}

func TestFix(t *testing.T) {
	h := New[int]()
	verify(t, h, 0)

	for i := 200; i > 0; i -= 10 {
		h.Push(i)
	}
	verify(t, h, 0)

	if x, _ := h.Peek(); x != 10 {
		t.Fatalf("Expected head to be 10, was %d", x)
	}
	h.values[0] = 210
	h.Fix(0)
	verify(t, h, 0)

	for i := 100; i > 0; i-- {
		elem := rand.IntN(h.Len())
		if i&1 == 0 {
			h.values[elem] *= 2
		} else {
			h.values[elem] /= 2
		}
		h.Fix(elem)
		verify(t, h, 0)
	}
}

type testInt int

func (u testInt) Less(o testInt) bool {
	return u < o
}

type testStudent struct {
	Name  string
	Score int64
}

func (u testStudent) Less(o testStudent) bool {
	if u.Score == o.Score {
		return u.Name < o.Name
	}
	return u.Score > o.Score
}

func TestHeap_Empty(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "empty", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New[testInt]()
			if got := (h.Len() == 0); got != tt.want {
				t.Errorf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testOpType int

const (
	testPush  = 1
	testPop   = 2
	testTop   = 3
	testEmpty = 4
)

type testOp[T any] struct {
	typ     testOpType
	x       T
	isEmpty bool
}

type testStruct[T any] struct {
	name string
	ops  []testOp[T]
}

func TestHeapExample1(t *testing.T) {
	tests1 := []testStruct[testInt]{
		{
			name: "example 1",
			ops: []testOp[testInt]{
				{typ: testEmpty, isEmpty: true},
				{typ: testPush, x: 10},
				{typ: testEmpty, isEmpty: false},
				{typ: testTop, x: 10},
				{typ: testPop},
				{typ: testEmpty, isEmpty: true},
				{typ: testPush, x: 9},
				{typ: testPush, x: 8},
				{typ: testPop},
				{typ: testPush, x: 3},
				{typ: testTop, x: 3},
				{typ: testPush, x: 2},
				{typ: testTop, x: 2},
				{typ: testPush, x: 4},
				{typ: testPush, x: 6},
				{typ: testPush, x: 5},
				{typ: testTop, x: 2},
				{typ: testPop},
				{typ: testTop, x: 3},
				{typ: testPop},
				{typ: testPop},
				{typ: testTop, x: 5},
				{typ: testEmpty, isEmpty: false},
			},
		},
	}
	testFunc(t, tests1, testInt.Less)
}

func TestHeapExample2(t *testing.T) {
	tests1 := []testStruct[testStudent]{
		{
			name: "example 2",
			ops: []testOp[testStudent]{
				{typ: testPush, x: testStudent{Name: "Alan", Score: 87}},
				{typ: testPush, x: testStudent{Name: "Bob", Score: 98}},
				{typ: testTop, x: testStudent{Name: "Bob", Score: 98}},
				{typ: testPop},
				{typ: testPush, x: testStudent{Name: "Carl", Score: 70}},
				{typ: testTop, x: testStudent{Name: "Alan", Score: 87}},
			},
		},
	}
	testFunc(t, tests1, testStudent.Less)
}

func testFunc[T any](t *testing.T, tests []testStruct[T], less func(a, b T) bool) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewFunc(less)
			for i, op := range tt.ops {
				switch op.typ {
				case testPush:
					oldSize := h.Len()
					h.Push(op.x)
					newSize := h.Len()
					if oldSize+1 != newSize {
						t.Errorf("op %d testPush %v failed", i, op.x)
					}
				case testPop:
					oldSize := h.Len()
					h.Pop()
					newSize := h.Len()
					if oldSize-1 != newSize {
						t.Errorf("op %d testPop %v failed", i, op.x)
					}
				case testTop:
					if got, _ := h.Peek(); !reflect.DeepEqual(got, op.x) {
						t.Errorf("op %d testTop %v, want %v", i, got, op.x)
					}
				case testEmpty:
					if got := (h.Len() == 0); got != op.isEmpty {
						t.Errorf("op %d Empty() = %v, want %v", i, got, op.isEmpty)
					}
				}
			}
		})
	}
}

func TestMaxHeap_BuildMaxHeap(t *testing.T) {
	t.Parallel()

	values := []int{6, 5, 2, 4, 7, 10, 12, 1, 3, 8, 9, 11}
	h := NewFunc(func(a, b int) bool {
		return a > b
	}, values...)

	expected := []int{12, 9, 11, 4, 8, 10, 2, 1, 3, 5, 7, 6}
	assert.Equal(t, expected, h.values)
	assert.Equal(t, 12, h.Len())
}

func TestMaxHeap_Push(t *testing.T) {
	t.Parallel()

	h := NewFunc(func(a, b int) bool {
		return a > b
	})

	values := []int{6, 5, 2, 4, 7, 10, 12, 1, 3, 8, 9, 11}
	for _, v := range values {
		h.Push(v)
	}

	expected := []int{12, 9, 11, 4, 8, 10, 7, 1, 3, 5, 6, 2}
	assert.Equal(t, expected, h.values)
	assert.Equal(t, 12, h.Len())
}

func TestMaxHeap_Pop(t *testing.T) {
	t.Parallel()

	h := NewFunc(func(a, b int) bool {
		return a > b
	})

	_, ok := h.Pop()
	assert.Equal(t, false, ok)

	values := []int{6, 5, 2, 4, 7, 10, 12, 1, 3, 8, 9, 11}
	for _, v := range values {
		h.Push(v)
	}

	val, ok := h.Pop()
	assert.Equal(t, 12, val)
	assert.Equal(t, true, ok)

	assert.Equal(t, 11, h.Len())
}

func TestMaxHeap_Peek(t *testing.T) {
	t.Parallel()

	h := NewFunc(func(a, b int) bool {
		return a > b
	})

	_, ok := h.Peek()
	assert.Equal(t, false, ok)

	values := []int{6, 5, 2, 4, 7, 10, 12, 1, 3, 8, 9, 11}
	for _, v := range values {
		h.Push(v)
	}

	val, ok := h.Peek()
	assert.Equal(t, 12, val)
	assert.Equal(t, true, ok)

	assert.Equal(t, 12, h.Len())
}

func TestBinaryHeapPush(t *testing.T) {
	h := New[int]()

	if actualValue := (h.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	h.Push(3)
	h.Push(2)
	h.Push(1)
	fmt.Println(h.values)
	fmt.Println(h.Values())
	if actualValue, expectedValue := h.Values(), []int{1, 2, 3}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := (h.Len() == 0); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := h.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := h.Peek(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestBinaryHeapValues(t *testing.T) {
	h := New(15, 20, 3, 1, 2)

	if actualValue, expectedValue := h.Values(), []int{1, 2, 3, 15, 20}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, ok := h.Pop(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestBinaryHeapPop(t *testing.T) {
	h := New[int]()

	if actualValue := (h.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	h.Push(3)
	h.Push(2)
	h.Push(1)
	h.Pop()

	if actualValue, ok := h.Peek(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := h.Pop(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := h.Pop(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := h.Pop(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := (h.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := h.Values(); len(actualValue) != 0 {
		t.Errorf("Got %v expected %v", actualValue, "[]")
	}
}

func TestBinaryHeapRandom(t *testing.T) {
	h := New[int]()

	for range 10000 {
		r := int(rand.Int32N(30))
		h.Push(r)
	}

	prev, _ := h.Pop()
	for !(h.Len() == 0) {
		curr, _ := h.Pop()
		if prev > curr {
			t.Errorf("Heap property invalidated. prev: %v current: %v", prev, curr)
		}
		prev = curr
	}
}

func TestBinaryHeapString(t *testing.T) {
	h := New[int]()
	h.Push(1)
	h.Push(5)
	h.Push(3)
	h.Push(1)
	h.Push(4)
	h.Push(2)
	fmt.Println(h)
	if !strings.HasPrefix(h.String(), "BinaryHeap") {
		t.Errorf("String should start with container name")
	}
}

func BenchmarkDup(b *testing.B) {
	const n = 10000
	h := New[int]()
	for b.Loop() {
		for range n {
			h.Push(0) // all elements are the same
		}
		for h.Len() > 0 {
			h.Pop()
		}
	}
}

func benchmarkPush(b *testing.B, h *Heap[int], size int) {
	for b.Loop() {
		for n := range size {
			h.Push(n)
		}
	}
}

func benchmarkPop(b *testing.B, h *Heap[int], size int) {
	for b.Loop() {
		for range size {
			h.Pop()
		}
	}
}

func BenchmarkBinaryHeapPop100(b *testing.B) {
	b.StopTimer()
	size := 100
	heap := New[int]()
	for n := range size {
		heap.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, heap, size)
}

func BenchmarkBinaryHeapPop1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	heap := New[int]()
	for n := range size {
		heap.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, heap, size)
}

func BenchmarkBinaryHeapPop10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	heap := New[int]()
	for n := range size {
		heap.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, heap, size)
}

func BenchmarkBinaryHeapPop100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	heap := New[int]()
	for n := range size {
		heap.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, heap, size)
}

func BenchmarkBinaryHeapPush100(b *testing.B) {
	b.StopTimer()
	size := 100
	heap := New[int]()
	b.StartTimer()
	benchmarkPush(b, heap, size)
}

func BenchmarkBinaryHeapPush1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	heap := New[int]()
	b.StartTimer()
	benchmarkPush(b, heap, size)
}

func BenchmarkBinaryHeapPush10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	heap := New[int]()
	b.StartTimer()
	benchmarkPush(b, heap, size)
}

func BenchmarkBinaryHeapPush100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	heap := New[int]()
	b.StartTimer()
	benchmarkPush(b, heap, size)
}
