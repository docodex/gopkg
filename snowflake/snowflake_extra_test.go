package snowflake_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/docodex/gopkg/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestNewWithNilOption(t *testing.T) {
	// Nil Option is skipped without error.
	s, err := snowflake.New(nil, snowflake.WithNode(7))
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, int64(7), s.Node())
}

func TestWithCheckNode(t *testing.T) {
	// Nil checkNode: option is a no-op, no error.
	s, err := snowflake.New(snowflake.WithNode(3), snowflake.WithCheckNode(nil))
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// Accepting checkNode.
	s, err = snowflake.New(
		snowflake.WithNode(3),
		snowflake.WithCheckNode(func(node int64) bool { return true }),
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), s.Node())

	// Rejecting checkNode returns ErrCheckNodeFailed.
	_, err = snowflake.New(
		snowflake.WithNode(3),
		snowflake.WithCheckNode(func(node int64) bool { return false }),
	)
	assert.True(t, errors.Is(err, snowflake.ErrCheckNodeFailed))
}

func TestNodeAccessor(t *testing.T) {
	s, err := snowflake.New(snowflake.WithNode(42))
	assert.NoError(t, err)
	assert.Equal(t, int64(42), s.Node())
}

func TestGenerateSequenceWrap(t *testing.T) {
	// Rapid generation forces the sequence to increment within the same
	// millisecond and eventually wrap, exercising the wait-for-next-time-unit
	// branch in Generate.
	s, err := snowflake.New(snowflake.WithNode(1))
	assert.NoError(t, err)

	// 2^11 = 2048 is maxSequence; generate > 2048 IDs as fast as possible.
	const n = 5000
	ids := make([]int64, 0, n)
	for range n {
		id, err := s.Generate()
		assert.NoError(t, err)
		ids = append(ids, id)
	}
	// All IDs must be unique and strictly monotonically increasing.
	seen := make(map[int64]struct{}, n)
	for i, id := range ids {
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id at %d: %d", i, id)
		}
		seen[id] = struct{}{}
		if i > 0 && id <= ids[i-1] {
			t.Fatalf("not monotonic at %d: %d <= %d", i, id, ids[i-1])
		}
	}
}

func TestGenerateConcurrent(t *testing.T) {
	s, err := snowflake.New(snowflake.WithNode(9))
	assert.NoError(t, err)

	const workers = 8
	const per = 1000
	results := make(chan int64, workers*per)
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range per {
				id, err := s.Generate()
				if err != nil {
					t.Errorf("generate: %v", err)
					return
				}
				results <- id
			}
		}()
	}
	wg.Wait()
	close(results)

	seen := make(map[int64]struct{}, workers*per)
	for id := range results {
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id across goroutines: %d", id)
		}
		seen[id] = struct{}{}
	}
}

func TestOptionErrorPropagates(t *testing.T) {
	// An Option returning error must cause New to fail with that error.
	_, err := snowflake.New(snowflake.WithNode(1 << 30))
	assert.True(t, errors.Is(err, snowflake.ErrOverNodeLimit))
}
