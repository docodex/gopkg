package treebidimap_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/dict/treebidimap"
	"github.com/stretchr/testify/assert"
)

func TestMapPut(t *testing.T) {
	m := treebidimap.New[int, string]()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") //overwrite

	if actualValue := m.Len(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}

	keys := m.Keys()
	slices.Sort(keys)
	values := m.Values()
	slices.Sort(values)
	assert.Equal(t, keys, []int{1, 2, 3, 4, 5, 6, 7})
	assert.Equal(t, values, []string{"a", "b", "c", "d", "e", "f", "g"})

	// key,expectedValue,expectedFound
	tests1 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "", false},
	}

	for _, test := range tests1 {
		// retrievals
		actualValue, actualOk := m.Get(test[0].(int))
		if actualValue != test[1] || actualOk != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}
}

func TestMapRemove(t *testing.T) {
	m := treebidimap.New[int, string]()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") //overwrite

	m.Remove(5)
	m.Remove(6)
	m.Remove(7)
	m.Remove(8)
	m.Remove(5)

	keys := m.Keys()
	slices.Sort(keys)
	values := m.Values()
	slices.Sort(values)
	assert.Equal(t, keys, []int{1, 2, 3, 4})
	assert.Equal(t, values, []string{"a", "b", "c", "d"})

	if actualValue := m.Len(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 4)
	}

	tests2 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "", false},
		{6, "", false},
		{7, "", false},
		{8, "", false},
	}

	for _, test := range tests2 {
		actualValue, actualFound := m.Get(test[0].(int))
		if actualValue != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}

	m.Remove(1)
	m.Remove(4)
	m.Remove(2)
	m.Remove(3)
	m.Remove(2)
	m.Remove(2)

	assert.Equal(t, m.Keys(), []int{})
	assert.Equal(t, m.Values(), []string{})
	if actualValue := m.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestMapGetKey(t *testing.T) {
	m := treebidimap.New[int, string]()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") //overwrite

	// key,expectedValue,expectedFound
	tests1 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{0, "x", false},
	}

	for _, test := range tests1 {
		// retrievals
		actualValue, actualOk := m.GetKey(test[1].(string))
		if actualValue != test[0] || actualOk != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[0])
		}
	}
}

func TestMapSerialization(t *testing.T) {
	m := treebidimap.New[string, int]()
	m.Put("a", 1.0)
	m.Put("b", 2.0)
	m.Put("c", 3.0)

	fmt.Println(m)

	var err error
	assert := func() {
		keys := m.Keys()
		slices.Sort(keys)
		values := m.Values()
		slices.Sort(values)
		assert.Equal(t, keys, []string{"a", "b", "c"})
		assert.Equal(t, values, []int{1, 2, 3})
		if actualValue, expectedValue := m.Len(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	bytes, err := m.MarshalJSON()
	assert()

	m.Clear()
	fmt.Println(m)

	err = m.UnmarshalJSON(bytes)
	assert()
	fmt.Println(m)

	bytes, err = json.Marshal([]any{"a", "b", "c", m})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	err = json.Unmarshal([]byte(`{"a":1,"b":2}`), &m)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
}

func TestMapString(t *testing.T) {
	c := treebidimap.New[string, int]()
	c.Put("a", 1)
	if !strings.HasPrefix(c.String(), "TreeBidiMap") {
		t.Errorf("String should start with container name")
	}
}

func benchmarkGet(b *testing.B, m *treebidimap.Map[int, int], size int) {
	for b.Loop() {
		for n := range size {
			m.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, m *treebidimap.Map[int, int], size int) {
	for b.Loop() {
		for n := range size {
			m.Put(n, n)
		}
	}
}

func benchmarkRemove(b *testing.B, m *treebidimap.Map[int, int], size int) {
	for b.Loop() {
		for n := range size {
			m.Remove(n)
		}
	}
}

func BenchmarkTreeBidiMapGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := treebidimap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeBidiMapGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := treebidimap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeBidiMapGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := treebidimap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeBidiMapGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	m := treebidimap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeBidiMapPut100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := treebidimap.New[int, int]()
	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeBidiMapPut1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := treebidimap.New[int, int]()
	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeBidiMapPut10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := treebidimap.New[int, int]()
	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeBidiMapPut100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	m := treebidimap.New[int, int]()
	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeBidiMapRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := treebidimap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkTreeBidiMapRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := treebidimap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkTreeBidiMapRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := treebidimap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkTreeBidiMapRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	m := treebidimap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkRemove(b, m, size)
}
