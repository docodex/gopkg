package doublylinkedlist_test

import (
	"cmp"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/list"
	"github.com/docodex/gopkg/container/list/doublylinkedlist"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	l := doublylinkedlist.New(1, 2, 3, 4)
	assert.Equal(t, l.FrontNode().Value, 1)
	assert.Equal(t, l.BackNode().Value, 4)
	fmt.Println(l.Values())

	l.PushFront(8, 7, 6, 5)
	assert.Equal(t, l.FrontNode().Value, 8)
	v, ok := l.Front()
	assert.True(t, ok)
	assert.Equal(t, v, 8)
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

func TestInsert(t *testing.T) {
	l := doublylinkedlist.New(1, 2, 3, 4)
	assert.Equal(t, l.FrontNode().Value, 1)
	assert.Equal(t, l.BackNode().Value, 4)
	fmt.Println(l.Values())

	x := l.InsertBefore(l.BackNode(), 8, 7, 6, 5)
	assert.Equal(t, x.Value, 8)
	fmt.Println(l.Values())

	x = l.InsertBefore(l.FrontNode(), 100)
	assert.Equal(t, x.Value, 100)
	assert.Equal(t, l.FrontNode().Value, 100)
	fmt.Println(l.Values())

	x = l.InsertAfter(l.FrontNode(), 10, 11, 12, 13)
	assert.Equal(t, x.Value, 10)
	fmt.Println(l.Values())

	x = l.InsertAfter(l.BackNode(), 200)
	assert.Equal(t, x.Value, 200)
	assert.Equal(t, l.BackNode().Value, 200)
	fmt.Println(l.Values())
}

func checkListLen(t *testing.T, l *doublylinkedlist.List[any], len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkListPointers(t *testing.T, l *doublylinkedlist.List[any], xs []*doublylinkedlist.Node[any]) {
	if !checkListLen(t, l, len(xs)) {
		return
	}

	for i, x := range xs {
		var prev *doublylinkedlist.Node[any]
		Prev := (*doublylinkedlist.Node[any])(nil)
		if i > 0 {
			prev = xs[i-1]
			Prev = prev
		}
		if y := x.Prev(); y != prev {
			t.Errorf("elt[%d](%p).prev = %p, want %p", i, x, y, prev)
		}
		if y := x.Prev(); y != Prev {
			t.Errorf("elt[%d](%p).Prev() = %p, want %p", i, x, y, Prev)
		}

		var next *doublylinkedlist.Node[any]
		Next := (*doublylinkedlist.Node[any])(nil)
		if i < len(xs)-1 {
			next = xs[i+1]
			Next = next
		}
		if n := x.Next(); n != next {
			t.Errorf("elt[%d](%p).next = %p, want %p", i, x, n, next)
		}
		if n := x.Next(); n != Next {
			t.Errorf("elt[%d](%p).Next() = %p, want %p", i, x, n, Next)
		}
	}
}

func TestList(t *testing.T) {
	l := doublylinkedlist.New[any]()
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{})

	// Single element list
	l.PushFront("a")
	x := l.FrontNode()
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x})
	l.MoveToFront(x)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x})
	l.MoveToBack(x)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x})
	l.Remove(x)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{})

	// Bigger list
	l.PushFront(2)
	x2 := l.FrontNode()
	l.PushFront(1)
	x1 := l.FrontNode()
	l.PushBack(3)
	x3 := l.BackNode()
	l.PushBack("banana")
	x4 := l.BackNode()
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x2, x3, x4})

	l.Remove(x2)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x3, x4})

	l.MoveToFront(x3) // move from middle
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x3, x1, x4})

	l.MoveToFront(x1)
	l.MoveToBack(x3) // move from middle
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x4, x3})

	l.MoveToFront(x3) // move from back
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x3, x1, x4})
	l.MoveToFront(x3) // should be no-op
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x3, x1, x4})

	l.MoveToBack(x3) // move from front
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x4, x3})
	l.MoveToBack(x3) // should be no-op
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x4, x3})

	x2 = l.InsertBefore(x1, 2) // insert before front
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x2, x1, x4, x3})
	l.Remove(x2)
	x2 = l.InsertBefore(x4, 2) // insert before middle
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x2, x4, x3})
	l.Remove(x2)
	x2 = l.InsertBefore(x3, 2) // insert before back
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x4, x2, x3})
	l.Remove(x2)

	x2 = l.InsertAfter(x1, 2) // insert after front
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x2, x4, x3})
	l.Remove(x2)
	x2 = l.InsertAfter(x4, 2) // insert after middle
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x4, x2, x3})
	l.Remove(x2)
	x2 = l.InsertAfter(x3, 2) // insert after back
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x4, x3, x2})
	l.Remove(x2)

	// Check standard iteration.
	sum := 0
	for x := l.FrontNode(); x != nil; x = x.Next() {
		if i, ok := x.Value.(int); ok {
			sum += i
		}
	}
	if sum != 4 {
		t.Errorf("sum over l = %d, want 4", sum)
	}

	// Clear all elements by iterating
	var next *doublylinkedlist.Node[any]
	for x := l.FrontNode(); x != nil; x = next {
		next = x.Next()
		l.Remove(x)
	}
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{})
}

func checkList(t *testing.T, l *doublylinkedlist.List[any], xs []any) {
	if !checkListLen(t, l, len(xs)) {
		return
	}

	i := 0
	for x := l.FrontNode(); x != nil; x = x.Next() {
		lx := x.Value.(int)
		if lx != xs[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, lx, xs[i])
		}
		i++
	}
}

func TestExtending(t *testing.T) {
	l1 := doublylinkedlist.New[any]()
	l2 := doublylinkedlist.New[any]()

	l1.PushBack(1)
	l1.PushBack(2)
	l1.PushBack(3)

	l2.PushBack(4)
	l2.PushBack(5)

	l3 := doublylinkedlist.New[any]()
	l3.PushBackList(l1)
	checkList(t, l3, []any{1, 2, 3})
	l3.PushBackList(l2)
	checkList(t, l3, []any{1, 2, 3, 4, 5})

	l3 = doublylinkedlist.New[any]()
	l3.PushFrontList(l2)
	checkList(t, l3, []any{4, 5})
	l3.PushFrontList(l1)
	checkList(t, l3, []any{1, 2, 3, 4, 5})

	checkList(t, l1, []any{1, 2, 3})
	checkList(t, l2, []any{4, 5})

	l3 = doublylinkedlist.New[any]()
	l3.PushBackList(l1)
	checkList(t, l3, []any{1, 2, 3})
	l3.PushBackList(l3)
	checkList(t, l3, []any{1, 2, 3, 1, 2, 3})

	l3 = doublylinkedlist.New[any]()
	l3.PushFrontList(l1)
	checkList(t, l3, []any{1, 2, 3})
	l3.PushFrontList(l3)
	checkList(t, l3, []any{1, 2, 3, 1, 2, 3})

	l3 = doublylinkedlist.New[any]()
	l1.PushBackList(l3)
	checkList(t, l1, []any{1, 2, 3})
	l1.PushFrontList(l3)
	checkList(t, l1, []any{1, 2, 3})
}

func TestRemove(t *testing.T) {
	l := doublylinkedlist.New[any]()
	l.PushBack(1)
	x1 := l.BackNode()
	l.PushBack(2)
	x2 := l.BackNode()
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x2})
	x := l.FrontNode()
	l.Remove(x)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x2})
	l.Remove(x)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x2})
}

func TestRemove2(t *testing.T) {
	l1 := doublylinkedlist.New[any]()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := doublylinkedlist.New[any]()
	l2.PushBack(3)
	l2.PushBack(4)

	x := l1.FrontNode()
	l2.Remove(x) // l2 should not change because x is not an element of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(x, 8)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestRemove3(t *testing.T) {
	l := doublylinkedlist.New[any]()
	l.PushBack(1)
	l.PushBack(2)

	x := l.FrontNode()
	l.Remove(x)
	if x.Value != 1 {
		t.Errorf("x.value = %d, want 1", x.Value)
	}
	if x.Next() != nil {
		t.Errorf("x.Next() != nil")
	}
	if x.Prev() != nil {
		t.Errorf("x.Prev() != nil")
	}
}

func TestMove(t *testing.T) {
	l := doublylinkedlist.New[any]()
	l.PushBack(1)
	x1 := l.BackNode()
	l.PushBack(2)
	x2 := l.BackNode()
	l.PushBack(3)
	x3 := l.BackNode()
	l.PushBack(4)
	x4 := l.BackNode()

	l.MoveAfter(x3, x3)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x2, x3, x4})
	l.MoveBefore(x2, x2)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x2, x3, x4})

	l.MoveAfter(x3, x2)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x2, x3, x4})
	l.MoveBefore(x2, x3)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x2, x3, x4})

	l.MoveBefore(x2, x4)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x3, x2, x4})
	x2, x3 = x3, x2

	l.MoveBefore(x4, x1)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x4, x1, x2, x3})
	x1, x2, x3, x4 = x4, x1, x2, x3

	l.MoveAfter(x4, x1)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x4, x2, x3})
	x2, x3, x4 = x4, x2, x3

	l.MoveAfter(x2, x3)
	checkListPointers(t, l, []*doublylinkedlist.Node[any]{x1, x3, x2, x4})
}

// Test that a list l is not modified when calling InsertBefore with a mark that is not an element of l.
func TestInsertBeforeUnknownMark(t *testing.T) {
	l := doublylinkedlist.New[any]()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.InsertBefore(new(doublylinkedlist.Node[any]), 1)
	checkList(t, l, []any{1, 2, 3})
}

// Test that a list l is not modified when calling InsertAfter with a mark that is not an element of l.
func TestInsertAfterUnknownMark(t *testing.T) {
	l := doublylinkedlist.New[any]()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.InsertAfter(new(doublylinkedlist.Node[any]), 1)
	checkList(t, l, []any{1, 2, 3})
}

// Test that a list l is not modified when calling MoveAfter or MoveBefore with a mark that is not an element of l.
func TestMoveUnknownMark(t *testing.T) {
	l1 := doublylinkedlist.New[any]()
	l1.PushBack(1)
	x1 := l1.BackNode()

	l2 := doublylinkedlist.New[any]()
	l2.PushBack(2)
	x2 := l2.BackNode()

	l1.MoveAfter(x1, x2)
	checkList(t, l1, []any{1})
	checkList(t, l2, []any{2})

	l1.MoveBefore(x1, x2)
	checkList(t, l1, []any{1})
	checkList(t, l2, []any{2})
}

func TestListNew(t *testing.T) {
	list1 := doublylinkedlist.New[any]()
	if actualValue := (list1.Len() == 0); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	list2 := doublylinkedlist.New[any](1, "b")
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
	l := doublylinkedlist.New[string]()
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
	l := doublylinkedlist.New[string]()
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
	l := doublylinkedlist.New[string]()
	l.PushBack("a")
	l.PushBack("b", "c")
	l.Del(2)
	if actualValue, ok := l.Get(2); actualValue != "" || ok {
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
	l := doublylinkedlist.New[string]()
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
	if actualValue, ok := l.Get(3); actualValue != "" || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	l.Del(0)
	if actualValue, ok := l.Get(0); actualValue != "b" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "b")
	}
}

func TestListSwap(t *testing.T) {
	l := doublylinkedlist.New[string]()
	l.PushBack("a")
	l.PushBack("b", "c")
	l.Swap(0, 1)
	if actualValue, ok := l.Get(0); actualValue != "b" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "c")
	}
}

func TestListSort(t *testing.T) {
	l := doublylinkedlist.New[string]()
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
	l := doublylinkedlist.New[string]()
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
	l := doublylinkedlist.New[string]()
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
	l := doublylinkedlist.New[any]()
	l.PushBack("a")
	l.PushBack("b", "c")
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s", l.Values()...), "abc"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestListIndexOf(t *testing.T) {
	l := doublylinkedlist.New[string]()
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
	l := doublylinkedlist.New[string]()
	expectedIndex := -1
	if index := doublylinkedlist.LastIndex(l, "a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	l.PushBack("a")
	l.PushBack("b", "c", "a")
	expectedIndex = 3
	if index := doublylinkedlist.LastIndex(l, "a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 1
	if index := doublylinkedlist.LastIndex(l, "b"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 2
	if index := doublylinkedlist.LastIndex(l, "c"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}
}

func TestListAdd(t *testing.T) {
	l := doublylinkedlist.New[any]()
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
	l.Clear()
	l.Add(0, "a")
	l.Add(0, "b", "c", "d")
	fmt.Println(l.String())
}

func TestListSet(t *testing.T) {
	l := doublylinkedlist.New[any]()
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
	l.PushBack("d")
	l.PushBack("c")
	l.Set(2, "cc") // last to first traversal
	l.Set(0, "aa") // first to last traversal
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s", l.Values()...), "aadcc"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestListEach(t *testing.T) {
	l := doublylinkedlist.New[string]()
	l.PushBack("a", "b", "c")
	l.Range(func(index int, value string) bool {
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
	l.RRange(func(index int, value string) bool {
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
	l := doublylinkedlist.New[string]()
	l.PushBack("a", "b", "c")
	mappedList := doublylinkedlist.New[string]()
	list.Map(mappedList, l, func(index int, value string) string {
		return "mapped: " + value
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
	l := doublylinkedlist.New[string]()
	l.PushBack("a", "b", "c")
	selectedList := doublylinkedlist.New[string]()
	list.Filter(selectedList, l, func(index int, value string) bool {
		return value >= "a" && value <= "b"
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
	l := doublylinkedlist.New[string]()
	l.PushBack("a", "b", "c")
	a := list.Any(l, func(index int, value string) bool {
		return value == "c"
	})
	if a != true {
		t.Errorf("Got %v expected %v", a, true)
	}
	a = list.Any(l, func(index int, value string) bool {
		return value == "x"
	})
	if a != false {
		t.Errorf("Got %v expected %v", a, false)
	}
}

func TestListAll(t *testing.T) {
	l := doublylinkedlist.New[string]()
	l.PushBack("a", "b", "c")
	all := list.All(l, func(index int, value string) bool {
		return value >= "a" && value <= "c"
	})
	if all != true {
		t.Errorf("Got %v expected %v", all, true)
	}
	all = list.All(l, func(index int, value string) bool {
		return value >= "a" && value <= "b"
	})
	if all != false {
		t.Errorf("Got %v expected %v", all, false)
	}
}

func TestListFind(t *testing.T) {
	l := doublylinkedlist.New[any]()
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
	l := doublylinkedlist.New[any]()
	l.PushBack("a", "b", "c", "a")
	foundIndex, foundValue := doublylinkedlist.FindLast(l, func(index int, value any) bool {
		return value.(string) == "a"
	})
	if foundValue != "a" || foundIndex != 3 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundIndex, "a", 3)
	}
	foundIndex, foundValue = list.Find(l, func(index int, value any) bool {
		return value.(string) == "x"
	})
	if foundValue != nil || foundIndex != -1 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundIndex, nil, nil)
	}
}

func TestListChaining(t *testing.T) {
	l := doublylinkedlist.New[any]()
	l.PushBack("a", "b", "c")
	chainedList := doublylinkedlist.New[any]()
	list.Filter(chainedList, l, func(index int, value any) bool {
		return value.(string) > "a"
	})
	l = chainedList
	chainedList = doublylinkedlist.New[any]()
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
	l := doublylinkedlist.New[any]()
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
	l := doublylinkedlist.New[any]()
	l.PushBack(1)
	if !strings.HasPrefix(l.String(), "DoublyLinkedList") {
		t.Errorf("String should start with container name")
	}
}

func benchmarkGet(b *testing.B, l *doublylinkedlist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.Get(n)
		}
	}
}

func benchmarkPushBack(b *testing.B, l *doublylinkedlist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.PushBack(n)
		}
	}
}

func benchmarkPushFront(b *testing.B, l *doublylinkedlist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.PushFront(n)
		}
	}
}

func benchmarkAdd(b *testing.B, l *doublylinkedlist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.Add(n, n)
			l.Del(0)
		}
	}
}

func benchmarkDelete(b *testing.B, l *doublylinkedlist.List[any], size int) {
	for b.Loop() {
		for n := range size {
			l.Del(n)
			l.Add(0, n)
		}
	}
}

func BenchmarkDoublyLinkedListGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := doublylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkGet(b, l, size)
}

func BenchmarkDoublyLinkedListGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := doublylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkGet(b, l, size)
}

func BenchmarkDoublyLinkedListGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := doublylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkGet(b, l, size)
}

func BenchmarkDoublyLinkedListGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := doublylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkGet(b, l, size)
}

func BenchmarkDoublyLinkedListPushBack100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := doublylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushBack(b, l, size)
}

func BenchmarkDoublyLinkedListPushBack1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := doublylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushBack(b, l, size)
}

func BenchmarkDoublyLinkedListPushBack10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := doublylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushBack(b, l, size)
}

func BenchmarkDoublyLinkedListPushBack100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := doublylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushBack(b, l, size)
}

func BenchmarkDoublyLinkedListPushFront100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := doublylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushFront(b, l, size)
}

func BenchmarkDoublyLinkedListPushFront1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := doublylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushFront(b, l, size)
}

func BenchmarkDoublyLinkedListPushFront10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := doublylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushFront(b, l, size)
}

func BenchmarkDoublyLinkedListPushFront100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := doublylinkedlist.New[any]()
	b.StartTimer()
	benchmarkPushFront(b, l, size)
}

func BenchmarkDoublyLinkedListAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := doublylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkAdd(b, l, size)
}

func BenchmarkDoublyLinkedListAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := doublylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkAdd(b, l, size)
}

func BenchmarkDoublyLinkedListAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := doublylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkAdd(b, l, size)
}

func BenchmarkDoublyLinkedListAdd100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := doublylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkAdd(b, l, size)
}

func BenchmarkDoublyLinkedListDelete100(b *testing.B) {
	b.StopTimer()
	size := 100
	l := doublylinkedlist.New[any]()
	for n := range size {
		l.PushFront(n)
	}
	b.StartTimer()
	benchmarkDelete(b, l, size)
}

func BenchmarkDoublyLinkedListDelete1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	l := doublylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkDelete(b, l, size)
}

func BenchmarkDoublyLinkedListDelete10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	l := doublylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkDelete(b, l, size)
}

func BenchmarkDoublyLinkedListDelete100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	l := doublylinkedlist.New[any]()
	for n := range size {
		l.PushBack(n)
	}
	b.StartTimer()
	benchmarkDelete(b, l, size)
}
