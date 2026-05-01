package set_test

import (
	"slices"
	"testing"

	"github.com/docodex/gopkg/container/set"
	"github.com/docodex/gopkg/container/set/hashset"
	"github.com/stretchr/testify/assert"
)

func sortedValues(s set.Set[int]) []int {
	out := make([]int, 0, s.Len())
	s.Range(func(v int) bool {
		out = append(out, v)
		return true
	})
	slices.Sort(out)
	return out
}

func TestIntersectionBasic(t *testing.T) {
	a := hashset.New(1, 2, 3, 4)
	b := hashset.New(2, 3, 5)
	c := hashset.New(3, 2, 6)
	dst := hashset.New[int]()

	set.Intersection(dst, a, b, c)
	assert.Equal(t, []int{2, 3}, sortedValues(dst))
}

func TestIntersectionAliasedDst(t *testing.T) {
	a := hashset.New(1, 2, 3)
	b := hashset.New(2, 3, 4)
	// dst aliases a
	set.Intersection(a, a, b)
	assert.Equal(t, []int{2, 3}, sortedValues(a))
}

func TestIntersectionEdgeCases(t *testing.T) {
	// nil dst: no-op
	set.Intersection(nil, hashset.New(1))

	// No src sets: clears dst
	dst := hashset.New(1, 2)
	set.Intersection(dst)
	assert.Equal(t, 0, dst.Len())

	// Nil src: clears dst
	dst = hashset.New(1, 2)
	var nilSet set.Set[int]
	set.Intersection(dst, nilSet)
	assert.Equal(t, 0, dst.Len())

	// Empty src: clears dst
	dst = hashset.New(1, 2)
	set.Intersection(dst, hashset.New[int]())
	assert.Equal(t, 0, dst.Len())
}

func TestUnionBasic(t *testing.T) {
	a := hashset.New(1, 2, 3)
	b := hashset.New(3, 4)
	c := hashset.New(5)
	dst := hashset.New[int]()
	set.Union(dst, a, b, c)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, sortedValues(dst))
}

func TestUnionAliasedAndNil(t *testing.T) {
	// nil dst: no-op
	set.Union(nil, hashset.New(1))

	// Aliased dst
	a := hashset.New(1, 2)
	b := hashset.New(2, 3)
	set.Union(a, a, b)
	assert.Equal(t, []int{1, 2, 3}, sortedValues(a))

	// Skip nil src members
	var nilSet set.Set[int]
	dst := hashset.New[int]()
	set.Union(dst, nilSet, hashset.New(7, 8))
	assert.Equal(t, []int{7, 8}, sortedValues(dst))

	// All nil src: dst ends up cleared and empty
	dst = hashset.New(1, 2)
	set.Union(dst, nilSet, nilSet)
	assert.Equal(t, 0, dst.Len())
}

func TestDifferenceBasic(t *testing.T) {
	a := hashset.New(1, 2, 3, 4)
	b := hashset.New(2, 4)
	dst := hashset.New[int]()
	set.Difference(dst, a, b)
	assert.Equal(t, []int{1, 3}, sortedValues(dst))
}

func TestDifferenceEdgeCases(t *testing.T) {
	// nil dst: no-op
	set.Difference(nil, hashset.New(1), hashset.New(2))

	// nil a: dst cleared
	dst := hashset.New(9)
	set.Difference(dst, nil, hashset.New(1))
	assert.Equal(t, 0, dst.Len())

	// nil b: dst equals a
	a := hashset.New(1, 2, 3)
	dst = hashset.New[int]()
	set.Difference(dst, a, nil)
	assert.Equal(t, []int{1, 2, 3}, sortedValues(dst))

	// Aliased dst == a
	a = hashset.New(1, 2, 3)
	b := hashset.New(2)
	set.Difference(a, a, b)
	assert.Equal(t, []int{1, 3}, sortedValues(a))
}
