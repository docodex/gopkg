package doublylinkedring_test

import (
	"fmt"
	"testing"

	"github.com/docodex/gopkg/container/ring/doublylinkedring"
)

func TestNewAndAdd(t *testing.T) {
	r := doublylinkedring.New(1, 2, 3, 4)
	fmt.Println(r.Values())
	r.Add(8, 7, 6, 5)
	fmt.Println(r.Values())
	r.Add(100)
	fmt.Println(r.Values())
}

// For debugging - keep around.
func dump[T any](r *doublylinkedring.Ring[T]) {
	if r == nil {
		fmt.Println("empty")
		return
	}
	i, n := 0, r.Len()
	for x := r; i < n; x = x.Next() {
		fmt.Printf("%4d: %p = {<- %p | %p ->}\n", i, x, x.Prev(), x.Next())
		i++
	}
	fmt.Println()
}

func verify(t *testing.T, r *doublylinkedring.Ring[int], N int, sum int) {
	if r == nil {
		return
	}

	dump(r)

	// Len
	n := r.Len()
	if n != N {
		t.Errorf("r.Len() == %d; expected %d", n, N)
	}

	// iteration
	n = 0
	s := 0
	r.Range(func(value int) bool {
		n++
		s += value
		return true
	})
	if n != N {
		t.Errorf("number of forward iterations == %d; expected %d", n, N)
	}
	if sum >= 0 && s != sum {
		t.Errorf("forward ring sum = %d; expected %d", s, sum)
	}

	// connections
	if r.Next() != nil {
		var x *doublylinkedring.Ring[int] // previous element
		for y := r; x == nil || y != r; y = y.Next() {
			if x != nil && x != y.Prev() {
				t.Errorf("prev = %p, expected q.prev = %p\n", x, y.Prev())
			}
			x = y
		}
		if x != r.Prev() {
			t.Errorf("prev = %p, expected r.prev = %p\n", x, r.Prev())
		}
	}

	// Move
	if r.Move(0) != r {
		t.Errorf("r.Move(0) != r")
	}
	if r.Move(N) != r {
		t.Errorf("r.Move(%d) != r", N)
	}
	if r.Move(-N) != r {
		t.Errorf("r.Move(%d) != r", -N)
	}
	for i := range 10 {
		ni := N + i
		mi := ni % N
		if r.Move(ni) != r.Move(mi) {
			t.Errorf("r.Move(%d) != r.Move(%d)", ni, mi)
		}
		if r.Move(-ni) != r.Move(-mi) {
			t.Errorf("r.Move(%d) != r.Move(%d)", -ni, -mi)
		}
	}
}

func TestCornerCases(t *testing.T) {
	var (
		r0 *doublylinkedring.Ring[int]
		r1 doublylinkedring.Ring[int]
	)
	// Basics
	verify(t, r0, 0, 0)
	verify(t, &r1, 1, 0)
	// Insert
	r1.Link(r0)
	verify(t, r0, 0, 0)
	verify(t, &r1, 1, 0)
	// Insert
	r1.Link(r0)
	verify(t, r0, 0, 0)
	verify(t, &r1, 1, 0)
	// Unlink
	r1.Unlink(0)
	verify(t, &r1, 1, 0)
}

func newDefault(t *testing.T, n int) *doublylinkedring.Ring[int] {
	if n <= 0 {
		return nil
	}
	r := doublylinkedring.New(0)
	for i := 1; i < n; i++ {
		r.Add(0)
	}
	size := r.Len()
	if size != n {
		t.Errorf("newDefault failed, r.Len=%d, expected r.Len=%d\n", size, n)
	}
	return r
}

func makeN(t *testing.T, n int) *doublylinkedring.Ring[int] {
	r := newDefault(t, n)
	for i := 1; i <= n; i++ {
		r.Value = i
		r = r.Next()
	}
	return r
}

func sumN(n int) int { return (n*n + n) / 2 }

func TestNew(t *testing.T) {
	for i := range 10 {
		r := newDefault(t, i)
		verify(t, r, i, -1)
	}
	for i := range 10 {
		r := makeN(t, i)
		verify(t, r, i, sumN(i))
	}
}

func TestLink1(t *testing.T) {
	r1a := makeN(t, 1)
	var r1b doublylinkedring.Ring[int]
	r2a := r1a.Link(&r1b)
	verify(t, r2a, 2, 1)
	if r2a != r1a {
		t.Errorf("a) 2-element link failed")
	}

	r2b := r2a.Link(r2a.Next())
	verify(t, r2b, 2, 1)
	if r2b != r2a.Next() {
		t.Errorf("b) 2-element link failed")
	}

	r1c := r2b.Link(r2b)
	verify(t, r1c, 1, 1)
	verify(t, r2b, 1, 0)
}

func TestLink2(t *testing.T) {
	var r0 *doublylinkedring.Ring[int]
	r1a := &doublylinkedring.Ring[int]{Value: 42}
	r1b := &doublylinkedring.Ring[int]{Value: 77}
	r10 := makeN(t, 10)

	r1a.Link(r0)
	verify(t, r1a, 1, 42)

	r1a.Link(r1b)
	verify(t, r1a, 2, 42+77)

	r10.Link(r0)
	verify(t, r10, 10, sumN(10))

	r10.Link(r1a)
	verify(t, r10, 12, sumN(10)+42+77)
}

func TestLink3(t *testing.T) {
	var r doublylinkedring.Ring[int]
	n := 1
	for i := 1; i < 10; i++ {
		n += i
		verify(t, r.Link(newDefault(t, i)), n, -1)
	}
}

func TestUnlink(t *testing.T) {
	r10 := makeN(t, 10)
	s10 := r10.Move(6)

	sum10 := sumN(10)

	verify(t, r10, 10, sum10)
	verify(t, s10, 10, sum10)

	r0 := r10.Unlink(0)
	verify(t, r0, 0, 0)

	r1 := r10.Unlink(1)
	verify(t, r1, 1, 2)
	verify(t, r10, 9, sum10-2)

	r9 := r10.Unlink(9)
	verify(t, r9, 9, sum10-2)
	verify(t, r10, 9, sum10-2)
}

func TestLinkUnlink(t *testing.T) {
	for i := 1; i < 4; i++ {
		ri := newDefault(t, i)
		for j := range i {
			rj := ri.Unlink(j)
			verify(t, rj, j, -1)
			verify(t, ri, i-j, -1)
			ri.Link(rj)
			verify(t, ri, i, -1)
		}
	}
}

// Test that calling Move() on an empty Ring initializes it.
func TestMoveEmptyRing(t *testing.T) {
	var r doublylinkedring.Ring[int]
	r.Move(1)
	verify(t, &r, 1, 0)
}

func TestDeleteRing(t *testing.T) {
	for i := range 5 {
		r := newDefault(t, i)
		if r == nil {
			continue
		}
		buf := make([]*doublylinkedring.Ring[int], 0, i)
		buf = append(buf, r)
		for n := r.Next(); n != r; n = n.Next() {
			buf = append(buf, n)
		}
		for _, x := range buf {
			fmt.Printf("%4d: %p = {<- %p | %p ->}\n", i, x, x.Prev(), x.Next())
			i++
		}
		fmt.Println("-------")
		doublylinkedring.Delete(r)
		for _, x := range buf {
			fmt.Printf("%4d: %p = {<- %p | %p ->}\n", i, x, x.Prev(), x.Next())
			i++
		}
		fmt.Println("=======")
	}
}
