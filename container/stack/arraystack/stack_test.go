package arraystack_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/stack/arraystack"
)

func TestStackPush(t *testing.T) {
	s := arraystack.New[any]()
	if actualValue := (s.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	s.Push(1)
	s.Push(2)
	s.Push(3)

	if actualValue := s.Values(); actualValue[0].(int) != 3 || actualValue[1].(int) != 2 || actualValue[2].(int) != 1 {
		t.Errorf("Got %v expected %v", actualValue, "[3,2,1]")
	}
	if actualValue := (s.Len() == 0); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := s.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := s.Peek(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestStackPeek(t *testing.T) {
	s := arraystack.New[any]()
	if actualValue, ok := s.Peek(); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if actualValue, ok := s.Peek(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestStackPop(t *testing.T) {
	s := arraystack.New[any]()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Pop()
	if actualValue, ok := s.Peek(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := s.Pop(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := s.Pop(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
	if actualValue, ok := s.Pop(); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := (s.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := s.Values(); len(actualValue) != 0 {
		t.Errorf("Got %v expected %v", actualValue, "[]")
	}
}

func TestStackSerialization(t *testing.T) {
	s := arraystack.New[any]()
	s.Push("a")
	s.Push("b")
	s.Push("c")

	var err error
	assert := func() {
		if actualValue, expectedValue := fmt.Sprintf("%s%s%s", s.Values()...), "cba"; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if actualValue, expectedValue := s.Len(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	bytes, err := s.MarshalJSON()
	assert()

	err = s.UnmarshalJSON(bytes)
	assert()

	bytes, err = json.Marshal([]any{"a", "b", "c", s})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	err = json.Unmarshal([]byte(`[1,2,3]`), &s)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
}

func TestStackString(t *testing.T) {
	s := arraystack.New[any]()
	s.Push(1)
	if !strings.HasPrefix(s.String(), "ArrayStack") {
		t.Errorf("String should start with container name")
	}
}

func benchmarkPush(b *testing.B, s *arraystack.Stack[any], size int) {
	for b.Loop() {
		for n := range size {
			s.Push(n)
		}
	}
}

func benchmarkPop(b *testing.B, s *arraystack.Stack[any], size int) {
	for b.Loop() {
		for range size {
			s.Pop()
		}
	}
}

func BenchmarkArrayStackPop100(b *testing.B) {
	b.StopTimer()
	size := 100
	s := arraystack.New[any]()
	for n := range size {
		s.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, s, size)
}

func BenchmarkArrayStackPop1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	s := arraystack.New[any]()
	for n := range size {
		s.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, s, size)
}

func BenchmarkArrayStackPop10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	s := arraystack.New[any]()
	for n := range size {
		s.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, s, size)
}

func BenchmarkArrayStackPop100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	s := arraystack.New[any]()
	for n := range size {
		s.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, s, size)
}

func BenchmarkArrayStackPush100(b *testing.B) {
	b.StopTimer()
	size := 100
	s := arraystack.New[any]()
	b.StartTimer()
	benchmarkPush(b, s, size)
}

func BenchmarkArrayStackPush1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	s := arraystack.New[any]()
	b.StartTimer()
	benchmarkPush(b, s, size)
}

func BenchmarkArrayStackPush10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	s := arraystack.New[any]()
	b.StartTimer()
	benchmarkPush(b, s, size)
}

func BenchmarkArrayStackPush100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	s := arraystack.New[any]()
	b.StartTimer()
	benchmarkPush(b, s, size)
}
