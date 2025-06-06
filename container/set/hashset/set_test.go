package hashset_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/set"
	"github.com/docodex/gopkg/container/set/hashset"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	var slice []string
	fmt.Println(slice == nil)
	s := hashset.New[string]()
	s.Add(slice...)
	fmt.Println("ok")
	fmt.Println(s)

	for i := range 10 {
		switch i {
		case 1:
			fmt.Println(i)
		case 2:
			fmt.Println(i)
		case 3:
			continue
		case 5:
			fmt.Println(i)
		default:
			return
		}
	}
}

func TestSet_IntersectionAndUnion(t *testing.T) {
	s1 := hashset.New(1, 2, 3, 4, 5)
	s2 := hashset.New(3, 4, 5, 6, 7)
	s3 := hashset.New(4, 5, 6, 7, 8, 9)
	s := hashset.NewWithCapacity[int](s1.Len())
	set.Intersection(s, s1, s2, s3)
	assert.True(t, s.Len() == 2)
	assert.True(t, s.Contains(4, 5))
	fmt.Println(s.Values())
	s.Clear()
	set.Union(s, s1, s2, s3)
	assert.True(t, s.Len() == 9)
	fmt.Println(s.Values())
}

func TestSet_Unmarshal(t *testing.T) {
	txt := "[1,2,3,4,5]"
	s1 := hashset.New[int]()
	fmt.Println(s1)
	e1 := s1.UnmarshalJSON([]byte(txt))
	assert.Nil(t, e1)
	fmt.Println(s1)
	fmt.Println(s1.String())
}

func TestSetNew(t *testing.T) {
	s := hashset.New(2, 1)
	if actualValue := s.Len(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue := s.Contains(1); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := s.Contains(2); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := s.Contains(3); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSetAdd(t *testing.T) {
	s := hashset.New[int]()
	s.Add()
	s.Add(1)
	s.Add(2)
	s.Add(2, 3)
	s.Add()
	if actualValue := (s.Len() == 0); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := s.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestSetContains(t *testing.T) {
	s := hashset.New[int]()
	s.Add(3, 1, 2)
	s.Add(2, 3)
	s.Add()
	if actualValue := s.Contains(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := s.Contains(1); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := s.Contains(1, 2, 3); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := s.Contains(1, 2, 3, 4); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
}

func TestSetRemove(t *testing.T) {
	s := hashset.New[int]()
	s.Add(3, 1, 2)
	s.Remove()
	if actualValue := s.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	s.Remove(1)
	if actualValue := s.Len(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	s.Remove(3)
	s.Remove(3)
	s.Remove()
	s.Remove(2)
	if actualValue := s.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestSetSerialization(t *testing.T) {
	s := hashset.New[string]()
	s.Add("a", "b", "c")

	var err error
	assert := func() {
		if actualValue, expectedValue := s.Len(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if actualValue := s.Contains("a", "b", "c"); actualValue != true {
			t.Errorf("Got %v expected %v", actualValue, true)
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

	err = json.Unmarshal([]byte(`["a","b","c"]`), &s)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
	assert()
}

func TestSetString(t *testing.T) {
	s := hashset.New[int]()
	s.Add(1)
	if !strings.HasPrefix(s.String(), "HashSet") {
		t.Errorf("String should start with container name")
	}
}

func TestSetIntersection(t *testing.T) {
	s := hashset.New[string]()
	another := hashset.New[string]()

	intersection := hashset.New[string]()
	set.Intersection(intersection, s, another)
	if actualValue, expectedValue := intersection.Len(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	s.Add("a", "b", "c", "d")
	another.Add("c", "d", "e", "f")

	intersection = hashset.New[string]()
	set.Intersection(intersection, s, another)

	if actualValue, expectedValue := intersection.Len(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := intersection.Contains("c", "d"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSetUnion(t *testing.T) {
	s := hashset.New[string]()
	another := hashset.New[string]()

	union := hashset.New[string]()
	set.Union(union, s, another)
	if actualValue, expectedValue := union.Len(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	s.Add("a", "b", "c", "d")
	another.Add("c", "d", "e", "f")

	union = hashset.New[string]()
	set.Union(union, s, another)

	if actualValue, expectedValue := union.Len(), 6; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := union.Contains("a", "b", "c", "d", "e", "f"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSetDifference(t *testing.T) {
	a := hashset.New[string]()
	b := hashset.New[string]()

	difference := hashset.New[string]()
	set.Difference(difference, a, b)
	if actualValue, expectedValue := difference.Len(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	a.Add("a", "b", "c", "d")
	b.Add("c", "d", "e", "f")

	difference = hashset.New[string]()
	set.Difference(difference, a, b)

	if actualValue, expectedValue := difference.Len(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := difference.Contains("a", "b"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSet_WithLock(t *testing.T) {
	s := hashset.New(1, 2, 3, 4, 5).WithLock()
	if actualValue := s.Len(); actualValue != 5 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue := s.Contains(1); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := s.Contains(2); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := s.Contains(6); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSet_WithLock_Add(t *testing.T) {
	s := hashset.New[int]().WithLock()
	s.Add()
	s.Add(1)
	s.Add(2)
	s.Add(2, 3)
	s.Add()
	if actualValue := (s.Len() == 0); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := s.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestSet_WithLock_Contains(t *testing.T) {
	s := hashset.New[int]().WithLock()
	s.Add(3, 1, 2)
	s.Add(2, 3)
	s.Add()
	if actualValue := s.Contains(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := s.Contains(1); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := s.Contains(1, 2, 3); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := s.Contains(1, 2, 3, 4); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
}

func TestSet_WithLock_Remove(t *testing.T) {
	s := hashset.New[int]().WithLock()
	s.Add(3, 1, 2)
	s.Remove()
	if actualValue := s.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	s.Remove(1)
	if actualValue := s.Len(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	s.Remove(3)
	s.Remove(3)
	s.Remove()
	s.Remove(2)
	if actualValue := s.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestSet_WithLock_String(t *testing.T) {
	s := hashset.New[int]().WithLock()
	s.Add(1)
	if !strings.HasPrefix(s.String(), "HashSet") {
		t.Errorf("String should start with container name")
	}
}

func benchmarkContains(b *testing.B, set *hashset.Set[int], size int) {
	for b.Loop() {
		for n := range size {
			set.Contains(n)
		}
	}
}

func benchmarkAdd(b *testing.B, set *hashset.Set[int], size int) {
	for b.Loop() {
		for n := range size {
			set.Add(n)
		}
	}
}

func benchmarkRemove(b *testing.B, set *hashset.Set[int], size int) {
	for b.Loop() {
		for n := range size {
			set.Remove(n)
		}
	}
}

func BenchmarkHashSetContains100(b *testing.B) {
	b.StopTimer()
	size := 100
	s := hashset.New[int]()
	for n := range size {
		s.Add(n)
	}
	b.StartTimer()
	benchmarkContains(b, s, size)
}

func BenchmarkHashSetContains1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	s := hashset.New[int]()
	for n := range size {
		s.Add(n)
	}
	b.StartTimer()
	benchmarkContains(b, s, size)
}

func BenchmarkHashSetContains10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	s := hashset.New[int]()
	for n := range size {
		s.Add(n)
	}
	b.StartTimer()
	benchmarkContains(b, s, size)
}

func BenchmarkHashSetContains100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	s := hashset.New[int]()
	for n := range size {
		s.Add(n)
	}
	b.StartTimer()
	benchmarkContains(b, s, size)
}

func BenchmarkHashSetAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	s := hashset.New[int]()
	b.StartTimer()
	benchmarkAdd(b, s, size)
}

func BenchmarkHashSetAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	s := hashset.New[int]()
	b.StartTimer()
	benchmarkAdd(b, s, size)
}

func BenchmarkHashSetAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	s := hashset.New[int]()
	b.StartTimer()
	benchmarkAdd(b, s, size)
}

func BenchmarkHashSetAdd100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	s := hashset.New[int]()
	b.StartTimer()
	benchmarkAdd(b, s, size)
}

func BenchmarkHashSetRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	s := hashset.New[int]()
	for n := range size {
		s.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, s, size)
}

func BenchmarkHashSetRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	s := hashset.New[int]()
	for n := range size {
		s.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, s, size)
}

func BenchmarkHashSetRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	s := hashset.New[int]()
	for n := range size {
		s.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, s, size)
}

func BenchmarkHashSetRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	s := hashset.New[int]()
	for n := range size {
		s.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, s, size)
}
