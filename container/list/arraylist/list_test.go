package arraylist_test

import (
	"cmp"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/list"
	"github.com/docodex/gopkg/container/list/arraylist"
	"github.com/stretchr/testify/assert"
)

func TestListNew(t *testing.T) {
	list1 := arraylist.New[any]()
	if actualValue := (list1.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	list2 := arraylist.New[any](1, "b")
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
	l := arraylist.New[any]()
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
	fmt.Println(l.Values())
}

func TestListPushFront(t *testing.T) {
	l := arraylist.New[any]()
	l.PushFront("a")
	l.PushFront("b", "c")
	if actualValue := (l.Len() == 0); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := l.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := l.Get(2); actualValue != "a" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "a")
	}
	fmt.Println(l.Values())
}

func TestListIndexOf(t *testing.T) {
	l := arraylist.New[any]()
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

func TestListLastIndexOf(t *testing.T) {
	l := arraylist.New[any]()
	expectedIndex := -1
	if index := arraylist.LastIndex(l, "a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	l.PushBack("a")
	l.PushBack("b", "c", "a")

	expectedIndex = 3
	if index := arraylist.LastIndex(l, "a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 1
	if index := arraylist.LastIndex(l, "b"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 2
	if index := arraylist.LastIndex(l, "c"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}
}

func TestListDelete(t *testing.T) {
	l := arraylist.New[any]()
	l.PushBack("a")
	l.PushBack("b", "c", "d")
	assert.Equal(t, l.Values(), []any{"a", "b", "c", "d"})
	l.Del(2)
	if actualValue, ok := l.Get(3); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	assert.Equal(t, l.Values(), []any{"a", "b", "d"})
	l.Del(0)
	assert.Equal(t, l.Values(), []any{"b", "d"})
	l.Del(1)
	assert.Equal(t, l.Values(), []any{"b"})
	l.Del(0)
	assert.Equal(t, l.Values(), []any{})
	l.Del(0) // no effect
	assert.Equal(t, l.Values(), []any{})
	if actualValue := (l.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := l.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	l.PushBack("a", "b", "c")
	l.Del(1)
	assert.Equal(t, l.Values(), []any{"a", "c"})
}

func TestListGet(t *testing.T) {
	l := arraylist.New[any]()
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
	l := arraylist.New[any]()
	l.PushBack("a")
	l.PushBack("b", "c")
	l.Swap(0, 1)
	if actualValue, ok := l.Get(0); actualValue != "b" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "b")
	}
}

func TestListSort(t *testing.T) {
	l := arraylist.New[string]()
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
	l := arraylist.New[any]()
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
	l := arraylist.New[any]()
	l.PushBack("a")
	l.PushBack("b", "c")
	if actualValue := list.Contains(l, "a"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := list.Contains(l, nil); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := list.Contains(l); actualValue != true {
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
	l := arraylist.New[any]()
	l.PushBack("a")
	l.PushBack("b", "c")
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s", l.Values()...), "abc"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestListAdd(t *testing.T) {
	l := arraylist.New[any]()
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
	l := arraylist.New[any]()
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
	l := arraylist.New[any]()
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
	l.RRange(func(index int, value any) bool {
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
	l := arraylist.New[any]()
	l.PushBack("a", "b", "c")
	mappedList := arraylist.New[any]()
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
	l := arraylist.New[any]()
	l.PushBack("a", "b", "c")
	selectedList := arraylist.New[any]()
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
	l := arraylist.New[any]()
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
	l := arraylist.New[any]()
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
	l := arraylist.New[any]()
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

func TestListFindLast(t *testing.T) {
	l := arraylist.New[any]()
	l.PushBack("a", "b", "c", "a")
	foundIndex, foundValue := arraylist.FindLast(l, func(index int, value any) bool {
		return value.(string) == "a"
	})
	if foundValue != "a" || foundIndex != 3 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundIndex, "a", 3)
	}
	foundIndex, foundValue = arraylist.FindLast(l, func(index int, value any) bool {
		return value.(string) == "x"
	})
	if foundValue != nil || foundIndex != -1 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundIndex, nil, nil)
	}
}

func TestListChaining(t *testing.T) {
	l := arraylist.New[any]()
	l.PushBack("a", "b", "c")
	chainedList := arraylist.New[any]()
	list.Filter(chainedList, l, func(index int, value any) bool {
		return value.(string) > "a"
	})
	l = chainedList
	chainedList = arraylist.New[any]()
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
	l := arraylist.New[any]()
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
	l := arraylist.New[any]()
	l.PushBack(1)
	if !strings.HasPrefix(l.String(), "ArrayList") {
		t.Errorf("String should start with container name")
	}
}

type SimpleList[T any] struct {
	values []T // current list values
	len    int // current list length
}

func NewSimpleList[T any](values ...T) *SimpleList[T] {
	l := new(SimpleList[T]).init()
	l.PushBack(values...)
	return l
}

func (l *SimpleList[T]) init() *SimpleList[T] {
	l.values = nil
	l.len = 0
	return l
}

func (l *SimpleList[T]) Len() int {
	return l.len
}

const defaultCapacity = 128

// checkAndExpand checks and expands the underlying array if necessary.
func (l *SimpleList[T]) checkAndExpand(delta int) {
	size := l.len + delta
	if size <= cap(l.values) {
		return
	}
	// expand & migrate
	capacity := max(size<<1, defaultCapacity)
	v := make([]T, capacity)
	copy(v, l.values[:l.len])
	l.values = v
}

// checkAndShrink checks and shrinks the underlying array if necessary.
func (l *SimpleList[T]) checkAndShrink() {
	if cap(l.values) <= defaultCapacity {
		return
	}
	if l.len<<2 > cap(l.values) {
		return
	}
	// shrink & migrate
	capacity := max(l.len<<1, defaultCapacity)
	v := make([]T, capacity)
	copy(v, l.values[:l.len])
	l.values = v
}

func (l *SimpleList[T]) PushBack(v ...T) {
	l.checkAndExpand(len(v))
	size := l.len + len(v)
	copy(l.values[l.len:size], v)
	l.len = size
}

func (l *SimpleList[T]) PushFront(v ...T) {
	l.checkAndExpand(len(v))
	size := l.len + len(v)
	copy(l.values[len(v):size], l.values[:l.len])
	copy(l.values[:len(v)], v)
	l.len = size
}

func (l *SimpleList[T]) Get(i int) (value T, ok bool) {
	if i >= 0 || i < l.len {
		value = l.values[i]
		ok = true
	}
	return
}

func (l *SimpleList[T]) Set(i int, v T) {
	if i >= 0 || i < l.len {
		l.values[i] = v
	}
}

func (l *SimpleList[T]) Add(i int, v ...T) {
	if i == l.len {
		l.PushBack(v...)
		return
	}
	if i < 0 || i >= l.len {
		return
	}
	l.checkAndExpand(len(v))
	size := l.len + len(v)
	j := i + len(v)
	copy(l.values[j:size], l.values[i:l.len])
	copy(l.values[i:j], v)
	l.len = size
}

func (l *SimpleList[T]) Delete(i int) {
	if i < 0 || i >= l.len {
		return
	}
	if i != l.len-1 {
		copy(l.values[i:l.len-1], l.values[i+1:l.len])
	}
	l.len--
	l.checkAndShrink()
}

func benchmarkSimpleGet(b *testing.B, l *SimpleList[any], size int) {
	for b.Loop() {
		for n := range size {
			l.Get(n)
		}
	}
}

func benchmarkSimplePushBack(b *testing.B, l *SimpleList[any], size int) {
	for b.Loop() {
		for n := range size {
			l.PushBack(n)
		}
	}
}

func benchmarkSimplePushFront(b *testing.B, l *SimpleList[any], size int) {
	for b.Loop() {
		for n := range size {
			l.PushFront(n)
		}
	}
}

func benchmarkSimpleAdd(b *testing.B, l *SimpleList[any], size int) {
	for b.Loop() {
		for n := range size {
			l.Add(n, n)
			l.Delete(l.Len() - 1)
		}
	}
}

func benchmarkSimpleDelete(b *testing.B, l *SimpleList[any], size int) {
	for b.Loop() {
		for n := range size {
			l.Delete(n)
			l.Add(l.Len(), n)
		}
	}
}

func benchmarkGet(b *testing.B, l *arraylist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.Get(n)
		}
	}
}

func benchmarkPushBack(b *testing.B, l *arraylist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.PushBack(n)
		}
	}
}

func benchmarkPushFront(b *testing.B, l *arraylist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.PushFront(n)
		}
	}
}

func benchmarkAdd(b *testing.B, l *arraylist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.Add(n, n)
			l.Del(l.Len() - 1)
		}
	}
}

func benchmarkDelete(b *testing.B, l *arraylist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.Del(n)
			l.Add(l.Len(), n)
		}
	}
}

func BenchmarkArrayListGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := arraylist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkGet(b, l, size)
}

func BenchmarkSimpleListGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := NewSimpleList[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkSimpleGet(b, l, size)
}

func BenchmarkArrayListGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := arraylist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkGet(b, l, size)
}

func BenchmarkSimpleListGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := NewSimpleList[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkSimpleGet(b, l, size)
}

func BenchmarkArrayListGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := arraylist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkGet(b, l, size)
}

func BenchmarkSimpleListGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := NewSimpleList[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkSimpleGet(b, l, size)
}

func BenchmarkArrayListGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := arraylist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkGet(b, l, size)
}

func BenchmarkSimpleListGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := NewSimpleList[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkSimpleGet(b, l, size)
}

func BenchmarkArrayListPushBack100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := arraylist.New[any]()
	b.StartTimer()
	benchmarkPushBack(b, l, size)
}

func BenchmarkSimpleListPushBack100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := NewSimpleList[any]()
	b.StartTimer()
	benchmarkSimplePushBack(b, l, size)
}

func BenchmarkArrayListPushBack1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := arraylist.New[any]()
	b.StartTimer()
	benchmarkPushBack(b, l, size)
}

func BenchmarkSimpleListPushBack1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := NewSimpleList[any]()
	b.StartTimer()
	benchmarkSimplePushBack(b, l, size)
}

func BenchmarkArrayListPushBack10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := arraylist.New[any]()
	b.StartTimer()
	benchmarkPushBack(b, l, size)
}

func BenchmarkSimpleListPushBack10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := NewSimpleList[any]()
	b.StartTimer()
	benchmarkSimplePushBack(b, l, size)
}

func BenchmarkArrayListPushBack100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := arraylist.New[any]()
	b.StartTimer()
	benchmarkPushBack(b, l, size)
}

func BenchmarkSimpleListPushBack100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := NewSimpleList[any]()
	b.StartTimer()
	benchmarkSimplePushBack(b, l, size)
}

func BenchmarkArrayListPushFront100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := arraylist.New[any]()
	b.StartTimer()
	benchmarkPushFront(b, l, size)
}

func BenchmarkSimpleListPushFront100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := NewSimpleList[any]()
	b.StartTimer()
	benchmarkSimplePushFront(b, l, size)
}

func BenchmarkArrayListPushFront1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := arraylist.New[any]()
	b.StartTimer()
	benchmarkPushFront(b, l, size)
}

func BenchmarkSimpleListPushFront1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := NewSimpleList[any]()
	b.StartTimer()
	benchmarkSimplePushFront(b, l, size)
}

func BenchmarkArrayListPushFront10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := arraylist.New[any]()
	b.StartTimer()
	benchmarkPushFront(b, l, size)
}

func BenchmarkSimpleListPushFront10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := NewSimpleList[any]()
	b.StartTimer()
	benchmarkSimplePushFront(b, l, size)
}

func BenchmarkArrayListPushFront100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := arraylist.New[any]()
	b.StartTimer()
	benchmarkPushFront(b, l, size)
}

func BenchmarkSimpleListPushFront100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := NewSimpleList[any]()
	b.StartTimer()
	benchmarkSimplePushFront(b, l, size)
}

func BenchmarkArrayListAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := arraylist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkAdd(b, l, size)
}

func BenchmarkSimpleListAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := NewSimpleList[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkSimpleAdd(b, l, size)
}

func BenchmarkArrayListAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := arraylist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkAdd(b, l, size)
}

func BenchmarkSimpleListAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := NewSimpleList[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkSimpleAdd(b, l, size)
}

func BenchmarkArrayListAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := arraylist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkAdd(b, l, size)
}

func BenchmarkSimpleListAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := NewSimpleList[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkSimpleAdd(b, l, size)
}

func BenchmarkArrayListAdd100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := arraylist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkAdd(b, l, size)
}

func BenchmarkSimpleListAdd100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := NewSimpleList[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkSimpleAdd(b, l, size)
}

func BenchmarkArrayListDelete100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := arraylist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkDelete(b, l, size)
}

func BenchmarkSimpleListDelete100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := NewSimpleList[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkSimpleDelete(b, l, size)
}

func BenchmarkArrayListDelete1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := arraylist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkDelete(b, l, size)
}

func BenchmarkSimpleListDelete1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := NewSimpleList[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkSimpleDelete(b, l, size)
}

func BenchmarkArrayListDelete10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := arraylist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkDelete(b, l, size)
}

func BenchmarkSimpleListDelete10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := NewSimpleList[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkSimpleDelete(b, l, size)
}

func BenchmarkArrayListDelete100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := arraylist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkDelete(b, l, size)
}

func BenchmarkSimpleListDelete100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := NewSimpleList[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkSimpleDelete(b, l, size)
}
