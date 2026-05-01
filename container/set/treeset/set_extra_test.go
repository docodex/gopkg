package treeset_test

import (
	"cmp"
	"strings"
	"sync"
	"testing"

	"github.com/docodex/gopkg/container/set/treeset"
	"github.com/stretchr/testify/assert"
)

func factories() []*treeset.Set[int] {
	rev := func(a, b int) int { return cmp.Compare(b, a) }
	return []*treeset.Set[int]{
		treeset.New[int](),
		treeset.NewFunc(rev),
		treeset.NewWithLock[int](),
		treeset.NewFuncWithLock(rev),
	}
}

func TestConstructorsSeed(t *testing.T) {
	s := treeset.New(1, 2, 3)
	assert.Equal(t, 3, s.Len())
	s2 := treeset.NewWithLock(1, 2, 3)
	assert.Equal(t, 3, s2.Len())
}

func TestContainsAnyAndClear(t *testing.T) {
	for _, s := range factories() {
		s.Add(1, 2, 3)
		assert.True(t, s.Contains(1, 2))
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
	s := treeset.NewWithLock(1, 2, 3)
	s.Range(nil) // no-op

	got := []int{}
	s.Range(func(v int) bool {
		got = append(got, v)
		return true
	})
	assert.Equal(t, []int{1, 2, 3}, got)

	count := 0
	s.Range(func(v int) bool { count++; return false })
	assert.Equal(t, 1, count)

	assert.True(t, strings.HasPrefix(s.String(), "TreeSet"))
}

func TestMarshalUnmarshal(t *testing.T) {
	s := treeset.NewWithLock(1, 2, 3)
	data, err := s.MarshalJSON()
	assert.NoError(t, err)

	s2 := treeset.NewWithLock[int]()
	assert.NoError(t, s2.UnmarshalJSON(data))
	assert.Equal(t, 3, s2.Len())

	// Invalid JSON.
	assert.Error(t, s2.UnmarshalJSON([]byte("not json")))
}

func TestConcurrentAccess(t *testing.T) {
	s := treeset.NewWithLock[int]()
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
