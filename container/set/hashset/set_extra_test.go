package hashset_test

import (
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/docodex/gopkg/container/set/hashset"
	"github.com/stretchr/testify/assert"
)

func factories() []*hashset.Set[int] {
	return []*hashset.Set[int]{
		hashset.New[int](),
		hashset.NewWithCapacity[int](0),
		hashset.NewWithCapacity[int](64),
		hashset.NewWithLock[int](),
		hashset.NewWithCapacityAndLock[int](0),
		hashset.NewWithCapacityAndLock[int](64),
	}
}

func TestConstructorsSeed(t *testing.T) {
	s := hashset.New(1, 2, 3)
	assert.Equal(t, 3, s.Len())
	s2 := hashset.NewWithLock(1, 2, 3)
	assert.Equal(t, 3, s2.Len())
}

func TestContainsAnyAndClear(t *testing.T) {
	for _, s := range factories() {
		s.Add(1, 2, 3)
		assert.True(t, s.Contains(1, 2, 3))
		assert.False(t, s.Contains(1, 9))
		assert.True(t, s.ContainsAny(9, 2))
		assert.False(t, s.ContainsAny(8, 9))
		assert.False(t, s.ContainsAny())

		s.Remove(2)
		assert.False(t, s.Contains(2))

		s.Clear()
		assert.Equal(t, 0, s.Len())
		s.Add(5)
		assert.True(t, s.Contains(5))
	}
}

func TestRangeAndString(t *testing.T) {
	s := hashset.NewWithLock(1, 2, 3)
	s.Range(nil) // no-op

	got := []int{}
	s.Range(func(v int) bool {
		got = append(got, v)
		return true
	})
	slices.Sort(got)
	assert.Equal(t, []int{1, 2, 3}, got)

	count := 0
	s.Range(func(v int) bool { count++; return false })
	assert.Equal(t, 1, count)

	assert.True(t, strings.HasPrefix(s.String(), "HashSet"))
}

func TestMarshalUnmarshal(t *testing.T) {
	s := hashset.NewWithLock(1, 2, 3)
	data, err := s.MarshalJSON()
	assert.NoError(t, err)

	s2 := hashset.NewWithLock[int]()
	assert.NoError(t, s2.UnmarshalJSON(data))
	assert.Equal(t, 3, s2.Len())

	// Invalid JSON.
	assert.Error(t, s2.UnmarshalJSON([]byte("not json")))
}

func TestConcurrentAccess(t *testing.T) {
	s := hashset.NewWithCapacityAndLock[int](16)
	var wg sync.WaitGroup
	for w := range 4 {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for i := range 50 {
				s.Add(base*100 + i)
				_ = s.Contains(base * 100)
				_ = s.ContainsAny(base * 100)
				_ = s.Len()
			}
		}(w)
	}
	wg.Wait()
	assert.GreaterOrEqual(t, s.Len(), 0)
}
