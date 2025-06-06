package deque_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/queue/deque"
)

type QueryStructure[T any] struct {
	queryType      string
	parameter      T
	expectedResult any
	expectedStatus bool
}

type TestCaseData[T any] struct {
	description string
	queries     []QueryStructure[T]
}

func TestDeque(t *testing.T) {
	// Test cases with ints as params
	testCasesInt := []TestCaseData[int]{
		{
			description: "Test empty deque",
			queries: []QueryStructure[int]{
				{
					queryType:      "Empty",
					expectedResult: true,
					expectedStatus: true,
				},
				{
					queryType:      "PeekFront",
					expectedStatus: false,
				},
				{
					queryType:      "PeekBack",
					expectedStatus: false,
				},
				{
					queryType:      "DequeueFront",
					expectedStatus: false,
				},
				{
					queryType:      "DeQueueRear",
					expectedStatus: false,
				},
				{
					queryType:      "Len",
					expectedResult: 0,
					expectedStatus: true,
				},
			},
		},
		{
			description: "Test deque with one element",
			queries: []QueryStructure[int]{
				{
					queryType:      "EnqueueFront",
					parameter:      1,
					expectedStatus: true,
				},
				{
					queryType:      "Empty",
					expectedResult: false,
					expectedStatus: true,
				},
				{
					queryType:      "PeekFront",
					expectedResult: 1,
					expectedStatus: true,
				},
				{
					queryType:      "PeekBack",
					expectedResult: 1,
					expectedStatus: true,
				},
				{
					queryType:      "Len",
					expectedResult: 1,
					expectedStatus: true,
				},
				{
					queryType:      "DequeueFront",
					expectedResult: 1,
					expectedStatus: true,
				},
				{
					queryType:      "Empty",
					expectedResult: true,
					expectedStatus: true,
				},
				{
					queryType:      "Len",
					expectedResult: 0,
					expectedStatus: true,
				},
			},
		},
		{
			description: "Test deque with multiple elements",
			queries: []QueryStructure[int]{
				{
					queryType:      "EnqueueFront",
					parameter:      1,
					expectedStatus: true,
				},
				{
					queryType:      "EnqueueFront",
					parameter:      2,
					expectedStatus: true,
				},
				{
					queryType:      "EnqueueBack",
					parameter:      3,
					expectedStatus: true,
				},
				{
					queryType:      "EnqueueBack",
					parameter:      4,
					expectedStatus: true,
				},
				{
					queryType:      "Empty",
					expectedResult: false,
					expectedStatus: true,
				},
				{
					queryType:      "PeekFront",
					expectedResult: 2,
					expectedStatus: true,
				},
				{
					queryType:      "PeekBack",
					expectedResult: 4,
					expectedStatus: true,
				},
				{
					queryType:      "Len",
					expectedResult: 4,
					expectedStatus: true,
				},
				{
					queryType:      "DequeueFront",
					expectedResult: 2,
					expectedStatus: true,
				},
				{
					queryType:      "DequeueBack",
					expectedResult: 4,
					expectedStatus: true,
				},
				{
					queryType:      "Empty",
					expectedResult: false,
					expectedStatus: true,
				},
				{
					queryType:      "Len",
					expectedResult: 2,
					expectedStatus: true,
				},
			},
		},
	}

	// Test cases with strings as params
	testCasesString := []TestCaseData[string]{
		{
			description: "Test one element deque",
			queries: []QueryStructure[string]{
				{
					queryType:      "EnqueueFront",
					parameter:      "a",
					expectedStatus: true,
				},
				{
					queryType:      "Empty",
					expectedResult: false,
					expectedStatus: true,
				},
				{
					queryType:      "PeekFront",
					expectedResult: "a",
					expectedStatus: true,
				},
				{
					queryType:      "PeekBack",
					expectedResult: "a",
					expectedStatus: true,
				},
				{
					queryType:      "Len",
					expectedResult: 1,
					expectedStatus: true,
				},
				{
					queryType:      "DequeueFront",
					expectedResult: "a",
					expectedStatus: true,
				},
				{
					queryType:      "Empty",
					expectedResult: true,
					expectedStatus: true,
				},
				{
					queryType:      "Len",
					expectedResult: 0,
					expectedStatus: true,
				},
			},
		},
		{
			description: "Test multiple elements deque",
			queries: []QueryStructure[string]{
				{
					queryType:      "EnqueueFront",
					parameter:      "a",
					expectedStatus: true,
				},
				{
					queryType:      "EnqueueFront",
					parameter:      "b",
					expectedStatus: true,
				},
				{
					queryType:      "EnqueueBack",
					parameter:      "c",
					expectedStatus: true,
				},
				{
					queryType:      "EnqueueBack",
					parameter:      "d",
					expectedStatus: true,
				},
				{
					queryType:      "Empty",
					expectedResult: false,
					expectedStatus: true,
				},
				{
					queryType:      "PeekFront",
					expectedResult: "b",
					expectedStatus: true,
				},
				{
					queryType:      "PeekBack",
					expectedResult: "d",
					expectedStatus: true,
				},
				{
					queryType:      "Len",
					expectedResult: 4,
					expectedStatus: true,
				},
				{
					queryType:      "DequeueFront",
					expectedResult: "b",
					expectedStatus: true,
				},
				{
					queryType:      "DequeueBack",
					expectedResult: "d",
					expectedStatus: true,
				},
				{
					queryType:      "Empty",
					expectedResult: false,
					expectedStatus: true,
				},
				{
					queryType:      "Len",
					expectedResult: 2,
					expectedStatus: true,
				},
			},
		},
	}

	// Run tests with ints
	for _, testCase := range testCasesInt {
		t.Run(testCase.description, func(t *testing.T) {
			q := deque.New[int]()
			for _, query := range testCase.queries {
				switch query.queryType {
				case "EnqueueFront":
					q.EnqueueFront(query.parameter)
				case "EnqueueBack":
					q.EnqueueBack(query.parameter)
				case "DequeueFront":
					result, ok := q.DequeueFront()
					if ok != query.expectedStatus {
						t.Errorf("Expected status: %v, got : %v", query.expectedStatus, ok)
					}
					if ok && result != query.expectedResult {
						t.Errorf("Expected %v, got %v", query.expectedResult, result)
					}
				case "DequeueBack":
					result, ok := q.DequeueBack()
					if ok != query.expectedStatus {
						t.Errorf("Expected status: %v, got : %v", query.expectedStatus, ok)
					}
					if ok && result != query.expectedResult {
						t.Errorf("Expected %v, got %v", query.expectedResult, result)
					}
				case "PeekFront":
					result, ok := q.PeekFront()
					if ok != query.expectedStatus {
						t.Errorf("Expected status: %v, got : %v", query.expectedStatus, ok)
					}
					if ok && result != query.expectedResult {
						t.Errorf("Expected %v, got %v, %v", query.expectedResult, result, testCase.description)
					}
				case "PeekBack":
					result, ok := q.PeekBack()
					if ok != query.expectedStatus {
						t.Errorf("Expected status: %v, got : %v", query.expectedStatus, ok)
					}
					if ok && result != query.expectedResult {
						t.Errorf("Expected %v, got %v", query.expectedResult, result)
					}
				case "Empty":
					result := (q.Len() == 0)
					if result != query.expectedResult {
						t.Errorf("Expected status: %v, got : %v", query.expectedResult, result)
					}
				case "Len":
					result := q.Len()
					if result != query.expectedResult {
						t.Errorf("Expected %v got %v", query.expectedResult, result)
					}
				}
			}
		})
	}

	// Run tests with strings
	for _, testCase := range testCasesString {
		t.Run(testCase.description, func(t *testing.T) {
			q := deque.New[string]()
			for _, query := range testCase.queries {
				switch query.queryType {
				case "EnqueueFront":
					q.EnqueueFront(query.parameter)
				case "EnqueueBack":
					q.EnqueueBack(query.parameter)
				case "DequeueFront":
					result, ok := q.DequeueFront()
					if ok != query.expectedStatus {
						t.Errorf("Expected status: %v, got : %v", query.expectedStatus, ok)
					}
					if ok && result != query.expectedResult {
						t.Errorf("Expected %v, got %v", query.expectedResult, result)
					}
				case "DequeueBack":
					result, ok := q.DequeueBack()
					if ok != query.expectedStatus {
						t.Errorf("Expected status: %v, got : %v", query.expectedStatus, ok)
					}
					if ok && result != query.expectedResult {
						t.Errorf("Expected %v, got %v", query.expectedResult, result)
					}
				case "PeekFront":
					result, ok := q.PeekFront()
					if ok != query.expectedStatus {
						t.Errorf("Expected status: %v, got : %v", query.expectedStatus, ok)
					}
					if ok && result != query.expectedResult {
						t.Errorf("Expected %v, got %v, %v", query.expectedResult, result, testCase.description)
					}
				case "PeekBack":
					result, ok := q.PeekBack()
					if ok != query.expectedStatus {
						t.Errorf("Expected status: %v, got : %v", query.expectedStatus, ok)
					}
					if ok && result != query.expectedResult {
						t.Errorf("Expected %v, got %v", query.expectedResult, result)
					}
				case "Empty":
					result := (q.Len() == 0)
					if result != query.expectedResult {
						t.Errorf("Expected %v, got %v", query.expectedResult, result)
					}
				case "Len":
					result := q.Len()
					if result != query.expectedResult {
						t.Errorf("Expected %v got %v", query.expectedResult, result)
					}
				}
			}
		})
	}
}

func TestQueueEnqueue(t *testing.T) {
	q := deque.New[any]()
	if actualValue := (q.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	q.EnqueueBack(1)
	q.EnqueueBack(2)
	q.EnqueueBack(3)

	if actualValue := q.Values(); actualValue[0].(int) != 1 || actualValue[1].(int) != 2 || actualValue[2].(int) != 3 {
		t.Errorf("Got %v expected %v", actualValue, "[1,2,3]")
	}
	if actualValue := (q.Len() == 0); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := q.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := q.PeekFront(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}

	q = deque.New[any]()
	if actualValue := (q.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	q.EnqueueFront(1)
	q.EnqueueFront(2)
	q.EnqueueFront(3)

	if actualValue := q.Values(); actualValue[0].(int) != 3 || actualValue[1].(int) != 2 || actualValue[2].(int) != 1 {
		t.Errorf("Got %v expected %v", actualValue, "[3,2,1]")
	}
	if actualValue := (q.Len() == 0); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := q.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := q.PeekBack(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestQueuePeek(t *testing.T) {
	q := deque.New[any]()
	if actualValue, ok := q.PeekFront(); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	q.EnqueueBack(1)
	q.EnqueueBack(2)
	q.EnqueueBack(3)
	if actualValue, ok := q.PeekFront(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}

	q = deque.New[any]()
	if actualValue, ok := q.PeekBack(); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	q.EnqueueFront(1)
	q.EnqueueFront(2)
	q.EnqueueFront(3)
	if actualValue, ok := q.PeekBack(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestQueueDequeue(t *testing.T) {
	q := deque.New[any]()
	q.EnqueueBack(1)
	q.EnqueueBack(2)
	q.EnqueueBack(3)
	q.DequeueFront()
	if actualValue, ok := q.PeekFront(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := q.DequeueFront(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := q.DequeueFront(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := q.DequeueFront(); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := (q.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := q.Values(); len(actualValue) != 0 {
		t.Errorf("Got %v expected %v", actualValue, "[]")
	}

	q = deque.New[any]()
	q.EnqueueFront(1)
	q.EnqueueFront(2)
	q.EnqueueFront(3)
	q.DequeueBack()
	if actualValue, ok := q.PeekBack(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := q.DequeueBack(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := q.DequeueBack(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := q.DequeueBack(); actualValue != nil || ok {
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
	q := deque.New[any]()
	q.EnqueueBack("a")
	q.EnqueueBack("b")
	q.EnqueueBack("c")

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
	q := deque.New[any]()
	q.EnqueueBack(1)
	if !strings.HasPrefix(q.String(), "DoubleEndedQueue") {
		t.Errorf("String should start with container name")
	}
}

func benchmarkEnqueueBack(b *testing.B, q *deque.Queue[any], size int) {
	for b.Loop() {
		for n := range size {
			q.EnqueueBack(n)
		}
	}
}

func benchmarkEnqueueFront(b *testing.B, q *deque.Queue[any], size int) {
	for b.Loop() {
		for n := range size {
			q.EnqueueFront(n)
		}
	}
}

func benchmarkDequeueFront(b *testing.B, q *deque.Queue[any], size int) {
	for b.Loop() {
		for range size {
			q.DequeueFront()
		}
	}
}

func benchmarkDequeueBack(b *testing.B, q *deque.Queue[any], size int) {
	for b.Loop() {
		for range size {
			q.DequeueBack()
		}
	}
}

func BenchmarkDoubleEndedQueueDequeueFront100(b *testing.B) {
	b.StopTimer()
	size := 100
	q := deque.New[any]()
	for n := range size {
		q.EnqueueBack(n)
	}
	b.StartTimer()
	benchmarkDequeueFront(b, q, size)
}

func BenchmarkDoubleEndedQueueDequeueFront1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	q := deque.New[any]()
	for n := range size {
		q.EnqueueBack(n)
	}
	b.StartTimer()
	benchmarkDequeueFront(b, q, size)
}

func BenchmarkDoubleEndedQueueDequeueFront10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	q := deque.New[any]()
	for n := range size {
		q.EnqueueBack(n)
	}
	b.StartTimer()
	benchmarkDequeueFront(b, q, size)
}

func BenchmarkDoubleEndedQueueDequeueFront100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	q := deque.New[any]()
	for n := range size {
		q.EnqueueBack(n)
	}
	b.StartTimer()
	benchmarkDequeueFront(b, q, size)
}

func BenchmarkDoubleEndedQueueDequeueBack100(b *testing.B) {
	b.StopTimer()
	size := 100
	q := deque.New[any]()
	for n := range size {
		q.EnqueueFront(n)
	}
	b.StartTimer()
	benchmarkDequeueBack(b, q, size)
}

func BenchmarkDoubleEndedQueueDequeueBack1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	q := deque.New[any]()
	for n := range size {
		q.EnqueueFront(n)
	}
	b.StartTimer()
	benchmarkDequeueBack(b, q, size)
}

func BenchmarkDoubleEndedQueueDequeueBack10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	q := deque.New[any]()
	for n := range size {
		q.EnqueueFront(n)
	}
	b.StartTimer()
	benchmarkDequeueBack(b, q, size)
}

func BenchmarkDoubleEndedQueueDequeueBack100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	q := deque.New[any]()
	for n := range size {
		q.EnqueueFront(n)
	}
	b.StartTimer()
	benchmarkDequeueBack(b, q, size)
}

func BenchmarkDoubleEndedQueueEnqueueBack100(b *testing.B) {
	b.StopTimer()
	size := 100
	q := deque.New[any]()
	b.StartTimer()
	benchmarkEnqueueBack(b, q, size)
}

func BenchmarkDoubleEndedQueueEnqueueBack1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	q := deque.New[any]()
	b.StartTimer()
	benchmarkEnqueueBack(b, q, size)
}

func BenchmarkDoubleEndedQueueEnqueueBack10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	q := deque.New[any]()
	b.StartTimer()
	benchmarkEnqueueBack(b, q, size)
}

func BenchmarkDoubleEndedQueueEnqueueBack100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	q := deque.New[any]()
	b.StartTimer()
	benchmarkEnqueueBack(b, q, size)
}

func BenchmarkDoubleEndedQueueEnqueueFront100(b *testing.B) {
	b.StopTimer()
	size := 100
	q := deque.New[any]()
	b.StartTimer()
	benchmarkEnqueueFront(b, q, size)
}

func BenchmarkDoubleEndedQueueEnqueueFront1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	q := deque.New[any]()
	b.StartTimer()
	benchmarkEnqueueFront(b, q, size)
}

func BenchmarkDoubleEndedQueueEnqueueFront10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	q := deque.New[any]()
	b.StartTimer()
	benchmarkEnqueueFront(b, q, size)
}

func BenchmarkDoubleEndedQueueEnqueueFront100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	q := deque.New[any]()
	b.StartTimer()
	benchmarkEnqueueFront(b, q, size)
}
