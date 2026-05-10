// Package retry runs a fallible operation under a configurable retry
// policy: max attempts, initial interval, max interval cap, backoff
// strategy, jitter, panic handler, retryable-error predicate, max
// elapsed time, on-retry hook, and optional error aggregation.
//
// [Do] reruns the operation while it returns a non-nil error, stopping
// when any of the following holds:
//   - the operation succeeds (returns nil)
//   - the configured maximum attempts has been reached
//   - the configured retryable predicate says the error is not retryable
//   - the configured maximum elapsed time has been reached
//   - the context is canceled (the sleep between attempts is also
//     canceled when ctx is done)
//
// If fn returns nil, [Do] returns nil regardless of the ctx state -
// success is reported even when ctx was canceled during fn. Callers
// that need ctx cancellation to take precedence over a successful fn
// should check ctx.Err() after [Do] returns.
//
// By default a panic in the [RetryableFunc] propagates to the caller.
// Pass [WithPanicHandler] to convert it into a returned error and stop
// retrying. Functions installed via [WithBackoff] and [WithRetryable]
// are NOT shielded by the package - they are expected to be pure and
// must not panic; a panic in them propagates to the caller of [Do].
// The hook installed via [WithOnRetry] IS shielded: a panic inside it
// is logged via slog.Default and swallowed so retry can continue.
package retry

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"runtime/debug"
	"time"

	"github.com/docodex/gopkg/internal"
)

// RetryableFunc is the operation re-run by [Do]. The ctx is the same
// ctx passed to [Do]. Returning nil ends the retry loop with success;
// returning a non-nil error triggers retry handling, subject to the
// configured retryable predicate, max attempts, and max elapsed time.
type RetryableFunc func(ctx context.Context) error

// PanicHandler converts a recovered panic value into an error. The
// returned error is treated like any other error returned by the
// operation: it is reported as the [Do] result and ends retrying. If
// no PanicHandler is configured, a panic propagates to the caller.
//
// The handler is expected to be pure and must not panic; a panic
// inside the handler propagates to the caller of [Do].
type PanicHandler func(ctx context.Context, recovered any) error

// Backoff returns the wait time before the next attempt. attempt is
// the number of failed attempts so far (1 after the first failure,
// 2 after the second, and so on). initInterval is the value
// configured via [WithInitInterval]. A non-positive return value
// falls back to initInterval.
type Backoff func(initInterval time.Duration, attempt int) time.Duration

// OnRetryFunc is invoked just before sleeping between attempts. The
// attempt argument is 1 after the first failure, 2 after the second,
// and so on; that is, "attempt N has just failed and a sleep for the
// (N+1)-th attempt is about to start". The hook is not invoked after
// the final failure (no further sleep follows). A panic raised inside
// the hook is recovered, logged via slog.Default, and swallowed so
// retry can continue.
//
// The hook is called synchronously and counts toward the wall-clock
// time observed by the retry loop. Keep it fast; offload slow work
// (network I/O, blocking metric emission) to a goroutine if needed.
//
// The hook may have already run for an attempt that does not actually
// execute, when ctx is canceled during the sleep that follows the
// hook. Treat the hook as "attempt N failed, sleep about to begin"
// rather than "attempt N+1 will run".
type OnRetryFunc func(ctx context.Context, attempt int, err error)

// ConstantBackoff returns initInterval regardless of attempt. It
// matches the default sleep behavior when no [Backoff] is configured
// but lets callers opt in explicitly via [WithConstantBackoff].
func ConstantBackoff(initInterval time.Duration, attempt int) time.Duration {
	return initInterval
}

// LinearBackoff returns initInterval * attempt, clamped at [maxBackoff]
// (24h) to avoid int64 overflow. attempt <= 1 returns initInterval
// as-is. initInterval <= 0 returns maxBackoff (consistent with
// [ExponentialBackoff]).
func LinearBackoff(initInterval time.Duration, attempt int) time.Duration {
	if attempt <= 1 {
		return initInterval
	}
	if initInterval <= 0 {
		return maxBackoff
	}
	// Overflow guard: if attempt exceeds the multiplier that would
	// produce maxBackoff, clamp directly without performing the
	// multiplication. The guard is sufficient on its own because
	// initInterval > 0 and attempt > 1 imply the product is > 0.
	if time.Duration(attempt) > maxBackoff/initInterval {
		return maxBackoff
	}
	return initInterval * time.Duration(attempt)
}

// ExponentialBackoff returns initInterval * 2^(attempt-1), clamped at
// [maxBackoff] (24h) to avoid int64 overflow. attempt <= 1 returns
// initInterval as-is. initInterval <= 0 returns maxBackoff.
func ExponentialBackoff(initInterval time.Duration, attempt int) time.Duration {
	if attempt <= 1 {
		return initInterval
	}
	if initInterval <= 0 {
		return maxBackoff
	}
	shift := attempt - 1
	if shift >= 62 {
		return maxBackoff
	}
	d := initInterval << shift
	if d <= 0 || d > maxBackoff {
		return maxBackoff
	}
	return d
}

const (
	defaultMaxAttempts  = 3
	defaultInitInterval = 500 * time.Millisecond
	maxBackoff          = 24 * time.Hour
)

// Do runs fn under the retry policy assembled from opts. See package
// doc for the stop conditions. Do panics if fn is nil.
func Do(ctx context.Context, fn RetryableFunc, opts ...Option) error {
	if fn == nil {
		panic("retry: RetryableFunc must not be nil")
	}
	o := options{
		maxAttempts:  defaultMaxAttempts,
		initInterval: defaultInitInterval,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}

	start := time.Now()
	var (
		lastErr error
		errs    []error
	)
	// joinCtxErr returns the error to surface when ctx is canceled in
	// the middle of the loop. last is the most recent attempt error
	// (which may already be in errs when collectErrors is on).
	joinCtxErr := func(last error) error {
		err := ctx.Err()
		if o.collectErrors {
			return errors.Join(append(errs, err)...)
		}
		return errors.Join(last, err)
	}

	for attempt := 1; ; attempt++ {
		if err := ctx.Err(); err != nil {
			if lastErr == nil {
				return err
			}
			return joinCtxErr(lastErr)
		}

		err := safeCall(ctx, fn, o.panicHandler)
		if err == nil {
			return nil
		}
		lastErr = err
		if o.collectErrors {
			errs = append(errs, err)
		}

		stop := attempt >= o.maxAttempts ||
			(o.retryable != nil && !o.retryable(err)) ||
			(o.maxElapsed > 0 && time.Since(start) >= o.maxElapsed)
		if stop {
			if o.collectErrors {
				return errors.Join(errs...)
			}
			return err
		}

		interval := o.initInterval
		if o.backoff != nil {
			if d := o.backoff(o.initInterval, attempt); d > 0 {
				interval = d
			}
		}
		if o.maxInterval > 0 && interval > o.maxInterval {
			interval = o.maxInterval
		}
		if o.jitterFactor > 0 {
			interval = applyJitter(interval, o.jitterFactor)
		}

		if o.onRetry != nil {
			safeOnRetry(ctx, o.onRetry, attempt, err)
		}

		if interval <= 0 {
			// Yield once so ctx cancellation and other goroutines
			// can be observed even when the configured backoff
			// returns 0 every iteration.
			select {
			case <-ctx.Done():
				return joinCtxErr(err)
			default:
			}
			continue
		}
		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return joinCtxErr(err)
		case <-timer.C:
		}
	}
}

// safeCall runs fn and converts a panic into an error via handler.
// If handler is nil, the panic propagates.
func safeCall(ctx context.Context, fn RetryableFunc, handler PanicHandler) (err error) {
	if handler != nil {
		defer func() {
			if r := recover(); r != nil {
				err = handler(ctx, r)
			}
		}()
	}
	return fn(ctx)
}

// safeOnRetry runs hook in a deferred recover so a panic inside the
// hook is logged via slog.Default and swallowed instead of aborting
// the retry loop.
func safeOnRetry(ctx context.Context, hook OnRetryFunc, attempt int, err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "panic in retry OnRetry hook",
				"error", internal.WrapPanic(r),
				"attempt", attempt,
				"stack", string(debug.Stack()),
			)
		}
	}()
	hook(ctx, attempt, err)
}

// applyJitter randomizes d by +/- factor of itself. The caller must
// pass factor in (0, 1]; the returned value is non-negative.
func applyJitter(d time.Duration, factor float64) time.Duration {
	delta := (rand.Float64()*2 - 1) * factor
	result := time.Duration(float64(d) * (1 + delta))
	if result < 0 {
		return 0
	}
	return result
}
