package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/docodex/gopkg/retry"
	"github.com/stretchr/testify/assert"
)

func TestDo_SuccessNoRetry(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), func(ctx context.Context) error {
		calls++
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, calls)
}

func TestDo_RetriesUntilMaxAttempts(t *testing.T) {
	sentinel := errors.New("boom")
	calls := 0
	err := retry.Do(context.Background(), func(ctx context.Context) error {
		calls++
		return sentinel
	},
		retry.WithMaxAttempts(3),
		retry.WithInitInterval(time.Millisecond),
	)

	assert.ErrorIs(t, err, sentinel)
	assert.Equal(t, 3, calls)
}

func TestDo_StopsOnFirstSuccessAfterFailure(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), func(ctx context.Context) error {
		calls++
		if calls < 2 {
			return errors.New("transient")
		}
		return nil
	},
		retry.WithMaxAttempts(5),
		retry.WithInitInterval(time.Microsecond),
	)

	assert.NoError(t, err)
	assert.Equal(t, 2, calls)
}

func TestDo_NilFnPanics(t *testing.T) {
	assert.PanicsWithValue(t, "retry: RetryableFunc must not be nil", func() {
		_ = retry.Do(context.Background(), nil)
	})
}

func TestDo_ContextAlreadyCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := retry.Do(ctx, func(ctx context.Context) error {
		calls++
		return nil
	})

	assert.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 0, calls, "fn must not run when ctx is already canceled")
}

func TestDo_SuccessWinsOverCtxCancel(t *testing.T) {
	// Documented contract: when fn returns nil, Do returns nil even
	// if ctx was canceled during fn. Callers that need cancellation
	// to take precedence must check ctx.Err() after Do returns.
	ctx, cancel := context.WithCancel(context.Background())

	calls := 0
	err := retry.Do(ctx, func(ctx context.Context) error {
		calls++
		cancel()
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, calls)
	assert.ErrorIs(t, ctx.Err(), context.Canceled, "ctx is still observably canceled")
}

func TestDo_ContextCanceledDuringSleep(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sentinel := errors.New("boom")
	calls := 0

	start := time.Now()
	err := retry.Do(ctx, func(ctx context.Context) error {
		calls++
		if calls == 1 {
			go func() {
				time.Sleep(10 * time.Millisecond)
				cancel()
			}()
		}
		return sentinel
	},
		retry.WithMaxAttempts(10),
		retry.WithInitInterval(time.Hour),
	)
	elapsed := time.Since(start)

	assert.ErrorIs(t, err, context.Canceled)
	assert.ErrorIs(t, err, sentinel)
	assert.Less(t, elapsed, time.Second, "ctx cancel must wake the sleep")
	assert.Equal(t, 1, calls)
}

func TestDo_PanicWithoutHandlerPropagates(t *testing.T) {
	assert.PanicsWithValue(t, "fn boom", func() {
		_ = retry.Do(context.Background(), func(ctx context.Context) error {
			panic("fn boom")
		})
	})
}

func TestDo_PanicHandlerConvertsToError(t *testing.T) {
	sentinel := errors.New("converted")
	calls := 0

	err := retry.Do(context.Background(), func(ctx context.Context) error {
		calls++
		panic("kaboom")
	},
		retry.WithPanicHandler(func(ctx context.Context, recovered any) error {
			return sentinel
		}),
		retry.WithInitInterval(time.Millisecond),
		retry.WithMaxAttempts(3),
	)

	assert.ErrorIs(t, err, sentinel)
	assert.Equal(t, 3, calls)
}

func TestDo_RetryablePredicateStopsImmediately(t *testing.T) {
	fatal := errors.New("fatal")
	calls := 0

	err := retry.Do(context.Background(), func(ctx context.Context) error {
		calls++
		return fatal
	},
		retry.WithMaxAttempts(5),
		retry.WithInitInterval(time.Millisecond),
		retry.WithRetryable(func(err error) bool {
			return !errors.Is(err, fatal)
		}),
	)

	assert.ErrorIs(t, err, fatal)
	assert.Equal(t, 1, calls, "non-retryable error must end the loop after one call")
}

func TestDo_MaxElapsedStops(t *testing.T) {
	sentinel := errors.New("boom")
	calls := 0
	start := time.Now()

	err := retry.Do(context.Background(), func(ctx context.Context) error {
		calls++
		return sentinel
	},
		retry.WithMaxAttempts(100),
		retry.WithInitInterval(20*time.Millisecond),
		retry.WithMaxElapsed(50*time.Millisecond),
	)
	elapsed := time.Since(start)

	assert.ErrorIs(t, err, sentinel)
	assert.Less(t, calls, 100, "must stop before maxAttempts when maxElapsed is hit")
	assert.Less(t, elapsed, 500*time.Millisecond)
}

func TestDo_ExponentialBackoffOverflowSafe(t *testing.T) {
	d := retry.ExponentialBackoff(500*time.Millisecond, 100)
	assert.Greater(t, d, time.Duration(0))
	assert.LessOrEqual(t, d, 24*time.Hour)
}

func TestConstantBackoff(t *testing.T) {
	for _, attempt := range []int{0, 1, 2, 10, 100} {
		assert.Equal(t, 500*time.Millisecond,
			retry.ConstantBackoff(500*time.Millisecond, attempt),
			"attempt=%d should still return initInterval", attempt)
	}
}

func TestLinearBackoff(t *testing.T) {
	initInterval := 100 * time.Millisecond
	cases := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, initInterval}, // attempt <= 1 returns initInterval
		{1, initInterval}, // attempt <= 1 returns initInterval
		{2, 200 * time.Millisecond},
		{5, 500 * time.Millisecond},
		{10, time.Second},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, retry.LinearBackoff(initInterval, c.attempt),
			"attempt=%d", c.attempt)
	}
}

func TestLinearBackoffOverflowSafe(t *testing.T) {
	// attempt=1_000_000 with init=1h would overflow; must clamp at 24h.
	d := retry.LinearBackoff(time.Hour, 1_000_000)
	assert.Greater(t, d, time.Duration(0))
	assert.LessOrEqual(t, d, 24*time.Hour)
}

func TestLinearBackoffNonPositiveInit(t *testing.T) {
	assert.Equal(t, 24*time.Hour, retry.LinearBackoff(-time.Second, 5))
	assert.Equal(t, 24*time.Hour, retry.LinearBackoff(0, 5))
}

func TestDo_WithLinearBackoff(t *testing.T) {
	// init=5ms, 4 retries: sleep totals 5+10+15+20 = 50ms.
	sentinel := errors.New("boom")
	calls := 0
	start := time.Now()

	err := retry.Do(context.Background(), func(ctx context.Context) error {
		calls++
		return sentinel
	},
		retry.WithMaxAttempts(5),
		retry.WithInitInterval(5*time.Millisecond),
		retry.WithLinearBackoff(),
	)
	elapsed := time.Since(start)

	assert.ErrorIs(t, err, sentinel)
	assert.Equal(t, 5, calls)
	assert.GreaterOrEqual(t, elapsed, 45*time.Millisecond)
	assert.Less(t, elapsed, 500*time.Millisecond)
}

func TestDo_WithConstantBackoff(t *testing.T) {
	sentinel := errors.New("boom")
	calls := 0
	start := time.Now()

	err := retry.Do(context.Background(), func(ctx context.Context) error {
		calls++
		return sentinel
	},
		retry.WithMaxAttempts(5),
		retry.WithInitInterval(5*time.Millisecond),
		retry.WithConstantBackoff(),
	)
	elapsed := time.Since(start)

	assert.ErrorIs(t, err, sentinel)
	assert.Equal(t, 5, calls)
	assert.GreaterOrEqual(t, elapsed, 18*time.Millisecond)
	assert.Less(t, elapsed, 500*time.Millisecond)
}

func TestDo_BackoffNonPositiveFallsBackToInit(t *testing.T) {
	sentinel := errors.New("boom")
	calls := 0
	start := time.Now()

	err := retry.Do(context.Background(), func(ctx context.Context) error {
		calls++
		return sentinel
	},
		retry.WithMaxAttempts(3),
		retry.WithInitInterval(5*time.Millisecond),
		retry.WithBackoff(func(initInterval time.Duration, attempt int) time.Duration {
			return -1
		}),
	)
	elapsed := time.Since(start)

	assert.ErrorIs(t, err, sentinel)
	assert.Equal(t, 3, calls)
	assert.GreaterOrEqual(t, elapsed, 8*time.Millisecond)
}

func TestDo_MaxIntervalCapsBackoff(t *testing.T) {
	// Without cap: 5ms, 10ms, 20ms, 40ms => total ~75ms.
	// With cap 12ms: 5ms, 10ms, 12ms, 12ms => total ~39ms.
	sentinel := errors.New("boom")
	calls := 0
	start := time.Now()

	err := retry.Do(context.Background(), func(ctx context.Context) error {
		calls++
		return sentinel
	},
		retry.WithMaxAttempts(5),
		retry.WithInitInterval(5*time.Millisecond),
		retry.WithMaxInterval(12*time.Millisecond),
		retry.WithExponentialBackoff(),
	)
	elapsed := time.Since(start)

	assert.ErrorIs(t, err, sentinel)
	assert.Less(t, elapsed, 70*time.Millisecond, "MaxInterval must cap backoff growth")
}

func TestDo_JitterDistributesSleep(t *testing.T) {
	// With factor 0.5 and base 20ms, sleep falls in [10ms, 30ms].
	// Run several times and check the observed total varies (jitter
	// is non-deterministic but at +/-50% it's hard to land on the
	// same value for 4 sleeps in a row).
	sentinel := errors.New("boom")
	totals := make([]time.Duration, 5)
	for i := range totals {
		start := time.Now()
		_ = retry.Do(context.Background(), func(ctx context.Context) error {
			return sentinel
		},
			retry.WithMaxAttempts(5),
			retry.WithInitInterval(20*time.Millisecond),
			retry.WithJitter(0.5),
		)
		totals[i] = time.Since(start)
	}
	allSame := true
	for i := 1; i < len(totals); i++ {
		if totals[i] != totals[0] {
			allSame = false
			break
		}
	}
	assert.False(t, allSame, "jitter must produce non-deterministic timings")
}

func TestDo_OnRetryHookFiresOnEachFailureExceptLast(t *testing.T) {
	sentinel := errors.New("boom")
	var hookAttempts []int
	var hookErrs []error

	err := retry.Do(context.Background(), func(ctx context.Context) error {
		return sentinel
	},
		retry.WithMaxAttempts(3),
		retry.WithInitInterval(time.Millisecond),
		retry.WithOnRetry(func(ctx context.Context, attempt int, err error) {
			hookAttempts = append(hookAttempts, attempt)
			hookErrs = append(hookErrs, err)
		}),
	)

	assert.ErrorIs(t, err, sentinel)
	// 3 attempts -> 2 retries -> 2 hook calls (final failure has no retry)
	assert.Equal(t, []int{1, 2}, hookAttempts)
	assert.Len(t, hookErrs, 2)
	for _, e := range hookErrs {
		assert.ErrorIs(t, e, sentinel)
	}
}

func TestDo_OnRetryHookPanicIsSwallowed(t *testing.T) {
	sentinel := errors.New("boom")
	var hookCalls atomic.Int32

	err := retry.Do(context.Background(), func(ctx context.Context) error {
		return sentinel
	},
		retry.WithMaxAttempts(3),
		retry.WithInitInterval(time.Millisecond),
		retry.WithOnRetry(func(ctx context.Context, attempt int, err error) {
			hookCalls.Add(1)
			panic("hook boom")
		}),
	)

	assert.ErrorIs(t, err, sentinel)
	assert.Equal(t, int32(2), hookCalls.Load(), "hook must run despite earlier panics")
}

func TestDo_CollectErrorsJoinsAll(t *testing.T) {
	e1 := errors.New("err1")
	e2 := errors.New("err2")
	e3 := errors.New("err3")
	seq := []error{e1, e2, e3}
	idx := 0

	err := retry.Do(context.Background(), func(ctx context.Context) error {
		e := seq[idx]
		idx++
		return e
	},
		retry.WithMaxAttempts(3),
		retry.WithInitInterval(time.Millisecond),
		retry.WithCollectErrors(),
	)

	assert.ErrorIs(t, err, e1)
	assert.ErrorIs(t, err, e2)
	assert.ErrorIs(t, err, e3)
}

func TestDo_CollectErrorsWithCtxCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	e1 := errors.New("err1")

	err := retry.Do(ctx, func(ctx context.Context) error {
		go func() {
			time.Sleep(5 * time.Millisecond)
			cancel()
		}()
		return e1
	},
		retry.WithMaxAttempts(10),
		retry.WithInitInterval(time.Hour),
		retry.WithCollectErrors(),
	)

	assert.ErrorIs(t, err, e1)
	assert.ErrorIs(t, err, context.Canceled)
}

// --- Option validation ---

func TestOption_WithMaxAttemptsZeroPanics(t *testing.T) {
	assert.Panics(t, func() {
		_ = retry.WithMaxAttempts(0)
	})
}

func TestOption_WithInitIntervalZeroPanics(t *testing.T) {
	assert.Panics(t, func() {
		_ = retry.WithInitInterval(0)
	})
}

func TestOption_WithMaxElapsedNegativeAPanics(t *testing.T) {
	assert.Panics(t, func() {
		_ = retry.WithMaxElapsed(-time.Second)
	})
}

func TestOption_WithMaxIntervalNegativePanics(t *testing.T) {
	assert.Panics(t, func() {
		_ = retry.WithMaxInterval(-time.Second)
	})
}

func TestOption_WithJitterOutOfRangePanics(t *testing.T) {
	assert.Panics(t, func() {
		_ = retry.WithJitter(-0.1)
	})
	assert.Panics(t, func() {
		_ = retry.WithJitter(1.1)
	})
}
