package retry

import "time"

// options holds the resolved retry policy assembled from Option
// values. Internal type, configured exclusively via [Option] setters.
type options struct {
	maxAttempts   int
	initInterval  time.Duration
	maxInterval   time.Duration
	backoff       Backoff
	jitterFactor  float64
	panicHandler  PanicHandler
	retryable     func(error) bool
	maxElapsed    time.Duration
	onRetry       OnRetryFunc
	collectErrors bool
}

// Option modifies the retry policy used by [Do]. Invalid values panic
// at Option construction time rather than being silently ignored, so
// misuse is caught loudly during development.
type Option func(o *options)

// WithMaxAttempts sets the maximum number of attempts. Default is 3.
func WithMaxAttempts(maxAttempts int) Option {
	if maxAttempts < 1 {
		panic("retry: WithMaxAttempts requires maxAttempts >= 1")
	}
	return func(o *options) {
		o.maxAttempts = maxAttempts
	}
}

// WithInitInterval sets the initial wait between attempts. Default is
// 500ms. Panics if initInterval <= 0.
func WithInitInterval(initInterval time.Duration) Option {
	if initInterval <= 0 {
		panic("retry: WithInitInterval requires initInterval > 0")
	}
	return func(o *options) {
		o.initInterval = initInterval
	}
}

// WithMaxInterval caps the per-attempt sleep at maxInterval. The cap
// is applied after [Backoff] but before jitter, so jitter still
// scatters callers around the cap. A value of 0 disables the cap (an
// internal 24h hard cap still applies inside [ExponentialBackoff] to
// prevent overflow). Panics if maxInterval < 0.
func WithMaxInterval(maxInterval time.Duration) Option {
	if maxInterval < 0 {
		panic("retry: WithMaxInterval requires maxInterval >= 0")
	}
	return func(o *options) {
		o.maxInterval = maxInterval
	}
}

// WithBackoff installs a custom backoff strategy. A nil backoff means
// every retry waits exactly initInterval (before maxInterval cap and
// jitter).
func WithBackoff(backoff Backoff) Option {
	return func(o *options) {
		o.backoff = backoff
	}
}

// WithConstantBackoff installs [ConstantBackoff] as the backoff
// strategy. The behavior is identical to installing no Backoff
// (every retry waits exactly initInterval), but the explicit form
// documents intent at the call site and pairs symmetrically with
// [WithLinearBackoff] and [WithExponentialBackoff].
func WithConstantBackoff() Option {
	return WithBackoff(ConstantBackoff)
}

// WithLinearBackoff installs [LinearBackoff] as the backoff strategy.
// Equivalent to WithBackoff(LinearBackoff).
func WithLinearBackoff() Option {
	return WithBackoff(LinearBackoff)
}

// WithExponentialBackoff installs [ExponentialBackoff] as the backoff
// strategy. Equivalent to WithBackoff(ExponentialBackoff).
func WithExponentialBackoff() Option {
	return WithBackoff(ExponentialBackoff)
}

// WithJitter randomizes the per-attempt sleep by +/- factor of the
// computed interval. For example WithJitter(0.5) on a 1s interval
// produces a uniform random sleep in [500ms, 1500ms]. Used to avoid
// thundering-herd effects when many clients retry at once. factor
// must be in [0, 1]; 0 disables jitter. Panics if factor is outside
// [0, 1].
func WithJitter(factor float64) Option {
	if factor < 0 || factor > 1 {
		panic("retry: WithJitter requires factor in [0, 1]")
	}
	return func(o *options) {
		o.jitterFactor = factor
	}
}

// WithPanicHandler installs a panic handler. Without it, a panic in
// the retried function propagates to the caller of [Do].
func WithPanicHandler(handler PanicHandler) Option {
	return func(o *options) {
		o.panicHandler = handler
	}
}

// WithRetryable installs a predicate that decides whether a non-nil
// error should be retried. Returning false stops the retry loop and
// surfaces the error to the caller. A nil predicate disables the
// check, so every error is retried until another stop condition fires.
func WithRetryable(retryable func(error) bool) Option {
	return func(o *options) {
		o.retryable = retryable
	}
}

// WithMaxElapsed bounds the total wall-clock time spent retrying.
// Once exceeded, [Do] returns the most recent error. A value of 0
// disables the bound. Panics if maxElapsed < 0.
func WithMaxElapsed(maxElapsed time.Duration) Option {
	if maxElapsed < 0 {
		panic("retry: WithMaxElapsed requires maxElapsed >= 0")
	}
	return func(o *options) {
		o.maxElapsed = maxElapsed
	}
}

// WithOnRetry installs a hook invoked on every failure that will be
// followed by another attempt. Useful for logging or metrics. The
// hook is not called after the final failure. A panic raised inside
// the hook is recovered, logged via slog.Default with a stack trace,
// and swallowed so retry can continue. See [OnRetryFunc] for the
// full contract.
func WithOnRetry(hook OnRetryFunc) Option {
	return func(o *options) {
		o.onRetry = hook
	}
}

// WithCollectErrors makes [Do] return all attempt errors joined via
// [errors.Join] instead of only the most recent one. Useful when
// every failure carries distinct context worth surfacing to the
// caller. Off by default.
func WithCollectErrors() Option {
	return func(o *options) {
		o.collectErrors = true
	}
}
