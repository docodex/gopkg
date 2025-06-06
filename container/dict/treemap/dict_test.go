package treemap_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/docodex/gopkg/container/dict/treemap"
	"github.com/stretchr/testify/assert"
)

func TestMapInit(t *testing.T) {
	m := map[int]string{
		1: "a",
		2: "b",
		3: "c",
	}
	v, err := json.Marshal(m)
	assert.Nil(t, err)
	m1 := treemap.New[int, string]()
	fmt.Println(m1)
	err = m1.UnmarshalJSON(v)
	assert.Nil(t, err)
	fmt.Println(m1)
	buf, _ := m1.MarshalJSON()
	fmt.Println(string(buf))
}

func TestMapPut(t *testing.T) {
	m := treemap.New[int, string]()
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

	// key,expectedValue,expectedOk
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
	m := treemap.New[int, string]()
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

func TestMapSerialization(t *testing.T) {
	m := treemap.New[string, float64]()
	m.Put("a", 1.0)
	m.Put("b", 2.0)
	m.Put("c", 3.0)

	var err error
	assert := func() {
		keys := m.Keys()
		slices.Sort(keys)
		values := m.Values()
		slices.Sort(values)
		assert.Equal(t, keys, []string{"a", "b", "c"})
		assert.Equal(t, values, []float64{1.0, 2.0, 3.0})
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

	err = m.UnmarshalJSON(bytes)
	assert()

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
	c := treemap.New[string, int]()
	c.Put("a", 1)
	if !strings.HasPrefix(c.String(), "TreeMap") {
		t.Errorf("String should start with container name")
	}
}

func benchmarkGet(b *testing.B, m *treemap.Map[int, int], size int) {
	for b.Loop() {
		for n := range size {
			m.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, m *treemap.Map[int, int], size int) {
	for b.Loop() {
		for n := range size {
			m.Put(n, n)
		}
	}
}

func benchmarkRemove(b *testing.B, m *treemap.Map[int, int], size int) {
	for b.Loop() {
		for n := range size {
			m.Remove(n)
		}
	}
}

func BenchmarkTreeMapGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := treemap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeMapGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := treemap.New[int, int]()
	for n := 0; n < size; n++ {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeMapGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := treemap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeMapGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	m := treemap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkGet(b, m, size)
}

func BenchmarkTreeMapPut100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := treemap.New[int, int]()
	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeMapPut1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := treemap.New[int, int]()
	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeMapPut10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := treemap.New[int, int]()
	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeMapPut100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	m := treemap.New[int, int]()
	b.StartTimer()
	benchmarkPut(b, m, size)
}

func BenchmarkTreeMapRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := treemap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkTreeMapRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := treemap.New[int, int]()
	for n := 0; n < size; n++ {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkTreeMapRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := treemap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkRemove(b, m, size)
}

func BenchmarkTreeMapRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	m := treemap.New[int, int]()
	for n := range size {
		m.Put(n, n)
	}
	b.StartTimer()
	benchmarkRemove(b, m, size)
}
