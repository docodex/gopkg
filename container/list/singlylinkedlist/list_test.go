package singlylinkedlist_test

import (
	"cmp"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/list"
	"github.com/docodex/gopkg/container/list/singlylinkedlist"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	l := singlylinkedlist.New(1, 2, 3, 4)
	assert.Equal(t, l.FrontNode().Value, 1)
	assert.Equal(t, l.FrontNode().Next().Value, 2)
	assert.Equal(t, l.BackNode().Value, 4)
	fmt.Println(l.Values())

	l.PushFront(8, 7, 6, 5)
	assert.Equal(t, l.FrontNode().Value, 8)
	v, ok := l.Front()
	assert.True(t, ok)
	assert.Equal(t, v, 8)
	fmt.Println(l.Values())
	l.PopFront()
	assert.Equal(t, l.FrontNode().Value, 7)
	fmt.Println(l.Values())

	l.PushFront(100)
	assert.Equal(t, l.FrontNode().Value, 100)
	fmt.Println(l.Values())

	l.PushBack(10, 11, 12, 13)
	assert.Equal(t, l.BackNode().Value, 13)
	assert.Nil(t, l.BackNode().Next())
	v, ok = l.Back()
	assert.True(t, ok)
	assert.Equal(t, v, 13)
	l.PopBack()
	assert.Equal(t, l.BackNode().Value, 12)
	fmt.Println(l.Values())

	l.PushBack(200)
	assert.Equal(t, l.BackNode().Value, 200)
	fmt.Println(l.Values())
}

func TestRemove(t *testing.T) {
	l := singlylinkedlist.New(1, 2, 3, 4, 5)
	fmt.Println(l.Values())
	v, ok := l.RemoveAfter(l.BackNode())
	assert.False(t, ok)
	fmt.Println(l.Values())
	v, ok = l.RemoveAfter(l.FrontNode())
	assert.True(t, ok)
	assert.Equal(t, v, 2)
	fmt.Println(l.Values())
	l.Del(5)
	assert.Equal(t, l.Len(), 4)
	fmt.Println(l.Values())
	l.Del(2)
	assert.Equal(t, l.Len(), 3)
	fmt.Println(l.Values())
}

func TestInsert(t *testing.T) {
	l := singlylinkedlist.New(1, 2, 3, 4)
	assert.Equal(t, l.FrontNode().Value, 1)
	assert.Equal(t, l.BackNode().Value, 4)
	fmt.Println(l.Values())

	x := l.InsertAfter(l.FrontNode(), 10, 11, 12, 13)
	assert.Equal(t, x.Value, 10)
	fmt.Println(l.Values())

	x = l.InsertAfter(l.BackNode(), 200)
	assert.Equal(t, x.Value, 200)
	assert.Equal(t, l.BackNode().Value, 200)
	fmt.Println(l.Values())
}

func TestPushList(t *testing.T) {
	l1 := singlylinkedlist.New[any](1, 2, 3, 4, 5)
	l2 := singlylinkedlist.New[any](11, 12, 13, 14, 15)
	l1.PushFrontList(l2)
	assert.Equal(t, 11, l1.FrontNode().Value)
	fmt.Println(l1.Values())
	l1.PushBackList(l2)
	assert.Equal(t, 15, l1.BackNode().Value)
	fmt.Println(l1.Values())
}

func TestListNew(t *testing.T) {
	list1 := singlylinkedlist.New[any]()
	if actualValue := (list1.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	list2 := singlylinkedlist.New[any](1, "b")
	if actualValue := list2.Len(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := list2.Get(0); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
	if actualValue, ok := list2.Get(1); actualValue != "b" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "b")
	}
	if actualValue, ok := list2.Get(2); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
}

func TestListPushBack(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a")
	l.PushBack("b", "c")
	if actualValue := (l.Len() == 0); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := l.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := l.Get(2); actualValue != "c" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "c")
	}
}

func TestListPushBackAndPushFront(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("b")
	l.PushFront("a")
	l.PushBack("c")
	if actualValue := (l.Len() == 0); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := l.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := l.Get(0); actualValue != "a" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "c")
	}
	if actualValue, ok := l.Get(1); actualValue != "b" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "c")
	}
	if actualValue, ok := l.Get(2); actualValue != "c" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "c")
	}
}

func TestListDelete(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a")
	l.PushBack("b", "c")
	l.Del(2)
	if actualValue, ok := l.Get(2); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	l.Del(1)
	l.Del(0)
	l.Del(0) // no effect
	if actualValue := (l.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := l.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestListGet(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a")
	l.PushBack("b", "c")
	if actualValue, ok := l.Get(0); actualValue != "a" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "a")
	}
	if actualValue, ok := l.Get(1); actualValue != "b" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "b")
	}
	if actualValue, ok := l.Get(2); actualValue != "c" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "c")
	}
	if actualValue, ok := l.Get(3); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	l.Del(0)
	if actualValue, ok := l.Get(0); actualValue != "b" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "b")
	}
}

func TestListSwap(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a")
	l.PushBack("b", "c")
	l.Swap(0, 1)
	if actualValue, ok := l.Get(0); actualValue != "b" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "c")
	}
}

func TestListSort(t *testing.T) {
	l := singlylinkedlist.New[string]()
	l.PushBack("e", "f", "g", "a", "b", "c", "d")
	l.Sort(cmp.Compare)
	for i := 1; i < l.Len(); i++ {
		a, _ := l.Get(i - 1)
		b, _ := l.Get(i)
		if a > b {
			t.Errorf("Not sorted! %s > %s", a, b)
		}
	}
}

func TestListClear(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("e", "f", "g", "a", "b", "c", "d")
	l.Clear()
	if actualValue := (l.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := l.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestListContains(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a")
	l.PushBack("b", "c")
	if actualValue := list.Contains(l, "a"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := list.ContainsAny(l, "a", "f"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := list.Contains(l, "a", "b", "c"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := list.Contains(l, "a", "b", "c", "d"); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := list.Contains(l, "e", "f", "g"); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	l.Clear()
	if actualValue := list.Contains(l, "a"); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := list.ContainsAny(l, "a"); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := list.Contains(l, "a", "b", "c"); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
}

func TestListValues(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a")
	l.PushBack("b", "c")
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s", l.Values()...), "abc"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestListIndexOf(t *testing.T) {
	l := singlylinkedlist.New[any]()
	expectedIndex := -1
	if index := list.Index(l, "a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	l.PushBack("a")
	l.PushBack("b", "c")

	expectedIndex = 0
	if index := list.Index(l, "a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 1
	if index := list.Index(l, "b"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 2
	if index := list.Index(l, "c"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}
}

func TestListAdd(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.Add(0, "b", "c")
	l.Add(0, "a")
	l.Add(10, "x") // ignore
	if actualValue := l.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	l.Add(3, "d") // append
	if actualValue := l.Len(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 4)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s%s", l.Values()...), "abcd"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestListSet(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.Set(0, "a")
	l.Set(1, "b")
	if actualValue := l.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
	l.PushBack("c") // append
	if actualValue := l.Len(); actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
	l.Set(4, "d")  // ignore
	l.Set(0, "bb") // update
	if actualValue := l.Len(); actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s", l.Values()...), "bb"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestListEach(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a", "b", "c")
	l.Range(func(index int, value any) bool {
		switch index {
		case 0:
			if actualValue, expectedValue := value, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
		return true
	})
}

func TestListMap(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a", "b", "c")
	mappedList := singlylinkedlist.New[any]()
	list.Map(mappedList, l, func(index int, value any) any {
		return "mapped: " + value.(string)
	})
	if actualValue, _ := mappedList.Get(0); actualValue != "mapped: a" {
		t.Errorf("Got %v expected %v", actualValue, "mapped: a")
	}
	if actualValue, _ := mappedList.Get(1); actualValue != "mapped: b" {
		t.Errorf("Got %v expected %v", actualValue, "mapped: b")
	}
	if actualValue, _ := mappedList.Get(2); actualValue != "mapped: c" {
		t.Errorf("Got %v expected %v", actualValue, "mapped: c")
	}
	if mappedList.Len() != 3 {
		t.Errorf("Got %v expected %v", mappedList.Len(), 3)
	}
}

func TestListFilter(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a", "b", "c")
	selectedList := singlylinkedlist.New[any]()
	list.Filter(selectedList, l, func(index int, value any) bool {
		return value.(string) >= "a" && value.(string) <= "b"
	})
	if actualValue, _ := selectedList.Get(0); actualValue != "a" {
		t.Errorf("Got %v expected %v", actualValue, "value: a")
	}
	if actualValue, _ := selectedList.Get(1); actualValue != "b" {
		t.Errorf("Got %v expected %v", actualValue, "value: b")
	}
	if selectedList.Len() != 2 {
		t.Errorf("Got %v expected %v", selectedList.Len(), 3)
	}
}

func TestListAny(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a", "b", "c")
	a := list.Any(l, func(index int, value any) bool {
		return value.(string) == "c"
	})
	if a != true {
		t.Errorf("Got %v expected %v", a, true)
	}
	a = list.Any(l, func(index int, value any) bool {
		return value.(string) == "x"
	})
	if a != false {
		t.Errorf("Got %v expected %v", a, false)
	}
}

func TestListAll(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a", "b", "c")
	all := list.All(l, func(index int, value any) bool {
		return value.(string) >= "a" && value.(string) <= "c"
	})
	if all != true {
		t.Errorf("Got %v expected %v", all, true)
	}
	all = list.All(l, func(index int, value any) bool {
		return value.(string) >= "a" && value.(string) <= "b"
	})
	if all != false {
		t.Errorf("Got %v expected %v", all, false)
	}
}

func TestListFind(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a", "b", "c")
	foundIndex, foundValue := list.Find(l, func(index int, value any) bool {
		return value.(string) == "c"
	})
	if foundValue != "c" || foundIndex != 2 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundIndex, "c", 2)
	}
	foundIndex, foundValue = list.Find(l, func(index int, value any) bool {
		return value.(string) == "x"
	})
	if foundValue != nil || foundIndex != -1 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundIndex, nil, nil)
	}
}

func TestListChaining(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a", "b", "c")
	chainedList := singlylinkedlist.New[any]()
	list.Filter(chainedList, l, func(index int, value any) bool {
		return value.(string) > "a"
	})
	l = chainedList
	chainedList = singlylinkedlist.New[any]()
	list.Map(chainedList, l, func(index int, value any) any {
		return value.(string) + value.(string)
	})
	if chainedList.Len() != 2 {
		t.Errorf("Got %v expected %v", chainedList.Len(), 2)
	}
	if actualValue, ok := chainedList.Get(0); actualValue != "bb" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "b")
	}
	if actualValue, ok := chainedList.Get(1); actualValue != "cc" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "c")
	}
}

func TestListSerialization(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack("a", "b", "c")

	var err error
	assert := func() {
		if actualValue, expectedValue := fmt.Sprintf("%s%s%s", l.Values()...), "abc"; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if actualValue, expectedValue := l.Len(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	bytes, err := l.MarshalJSON()
	assert()

	err = l.UnmarshalJSON(bytes)
	assert()

	bytes, err = json.Marshal([]any{"a", "b", "c", l})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	err = json.Unmarshal([]byte(`[1,2,3]`), &l)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
}

func TestListString(t *testing.T) {
	l := singlylinkedlist.New[any]()
	l.PushBack(1)
	if !strings.HasPrefix(l.String(), "SinglyLinkedList") {
		t.Errorf("String should start with container name")
	}
}

func benchmarkGet(b *testing.B, l *singlylinkedlist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.Get(n)
		}
	}
}

func benchmarkPushBack(b *testing.B, l *singlylinkedlist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.PushBack(n)
		}
	}
}

func benchmarkPushFront(b *testing.B, l *singlylinkedlist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.PushFront(n)
		}
	}
}

func benchmarkAdd(b *testing.B, l *singlylinkedlist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.Add(n, n)
			l.Del(0)
		}
	}
}

func benchmarkDelete(b *testing.B, l *singlylinkedlist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.Del(n)
			l.Add(0, n)
		}
	}
}

func BenchmarkSinglyLinkedListGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := singlylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkGet(b, l, size)
}

func BenchmarkSinglyLinkedListGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := singlylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkGet(b, l, size)
}

func BenchmarkSinglyLinkedListGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := singlylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkGet(b, l, size)
}

func BenchmarkSinglyLinkedListGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := singlylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkGet(b, l, size)
}

func BenchmarkSinglyLinkedListPushBack100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := singlylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushBack(b, l, size)
}

func BenchmarkSinglyLinkedListPushBack1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := singlylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushBack(b, l, size)
}

func BenchmarkSinglyLinkedListPushBack10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := singlylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushBack(b, l, size)
}

func BenchmarkSinglyLinkedListPushBack100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := singlylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushBack(b, l, size)
}

func BenchmarkSinglyLinkedListPushFront100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := singlylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushFront(b, l, size)
}

func BenchmarkSinglyLinkedListPushFront1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := singlylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushFront(b, l, size)
}

func BenchmarkSinglyLinkedListPushFront10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := singlylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushFront(b, l, size)
}

func BenchmarkSinglyLinkedListPushFront100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := singlylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushFront(b, l, size)
}

func BenchmarkSinglyLinkedListAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := singlylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkAdd(b, l, size)
}

func BenchmarkSinglyLinkedListInsert1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := singlylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkAdd(b, l, size)
}

func BenchmarkSinglyLinkedListAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := singlylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkAdd(b, l, size)
}

func BenchmarkSinglyLinkedListAdd100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := singlylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkAdd(b, l, size)
}

func BenchmarkSinglyLinkedListDelete100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := singlylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkDelete(b, l, size)
}

func BenchmarkSinglyLinkedListDelete1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := singlylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkDelete(b, l, size)
}

func BenchmarkSinglyLinkedListDelete10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := singlylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkDelete(b, l, size)
}

func BenchmarkSinglyLinkedListDelete100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := singlylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkDelete(b, l, size)
}
