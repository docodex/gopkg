// Package recovery provides helpers for catching and logging panics in
// deferred calls.
//
// Both [Recover] and [Recovery] must be invoked via defer. They use the
// built-in recover to swallow a panic, log a structured error including
// a stack trace via slog.Default, then run any user-supplied
// [CleanupFunc] cleanups.
//
// [Recover] additionally writes the panic into the *error you pass in
// so the surrounding function can return that error to its caller;
// [Recovery] is the variant that does not need an error pointer.
//
// Each cleanup is itself shielded by recover, so a panicking cleanup
// will not propagate out of this package.
//
// Logging uses slog.Default. To redirect output, install your logger
// via slog.SetDefault before the deferred call runs.
package recovery

import (
	"context"
	"log/slog"
	"runtime/debug"

	"github.com/docodex/gopkg/internal"
)

// CleanupFunc is the signature of a cleanup invoked after a panic
// has been recovered. Cleanups receive the same context passed into
// [Recover] / [Recovery]. They are executed in the order given,
// nil entries are skipped, and a panic inside one cleanup is
// logged but does not abort subsequent cleanups.
type CleanupFunc func(ctx context.Context)

// Recover must be called via defer. If the surrounding function is
// panicking, Recover swallows the panic, writes a wrapping error into
// *err when err is non-nil, logs the panic and stack trace via
// slog.Default, then invokes each non-nil cleanup in cleanups in
// order.
//
// When the panic value is itself an error, the wrapping uses %w so
// callers can use errors.Is / errors.As against the original value.
// Otherwise the value is formatted with %v.
//
// The write into *err happens BEFORE cleanups run, so cleanups
// observe the recovered error via the same pointer. Callers that
// need to inspect the pre-existing *err value should snapshot it
// before the deferred call. Note that this also overwrites any err
// the surrounding function had already set on its return values.
//
// If no panic is in flight, Recover is a no-op: cleanups are not
// invoked and *err is not touched.
func Recover(ctx context.Context, err *error, cleanups ...CleanupFunc) {
	if r := recover(); r != nil {
		recovered := internal.WrapPanic(r)
		if err != nil {
			*err = recovered
		}
		slog.ErrorContext(ctx, "panic recovered",
			"error", recovered,
			"stack", string(debug.Stack()),
		)
		invoke(ctx, cleanups...)
	}
}

// Recovery must be called via defer. If the surrounding function is
// panicking, Recovery swallows the panic, logs the panic and stack
// trace via slog.Default, then invokes each non-nil cleanup in
// cleanups in order. Use [Recover] if you also need to capture the
// panic into an *error for the surrounding function to return.
//
// If no panic is in flight, Recovery is a no-op: cleanups are not
// invoked.
func Recovery(ctx context.Context, cleanups ...CleanupFunc) {
	if r := recover(); r != nil {
		recovered := internal.WrapPanic(r)
		slog.ErrorContext(ctx, "panic recovered",
			"error", recovered,
			"stack", string(debug.Stack()),
		)
		invoke(ctx, cleanups...)
	}
}

// invoke runs each non-nil cleanup in turn, isolating each call so
// a panic in one cleanup does not abort subsequent cleanups nor
// leak out of the package.
func invoke(ctx context.Context, cleanups ...CleanupFunc) {
	for _, cleanup := range cleanups {
		if cleanup != nil {
			safeInvoke(ctx, cleanup)
		}
	}
}

// safeInvoke runs cleanup in a nested deferred recover so a panic
// inside the cleanup is logged but cannot leak out.
func safeInvoke(ctx context.Context, cleanup CleanupFunc) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "panic in recovery cleanup",
				"error", internal.WrapPanic(r),
				"stack", string(debug.Stack()),
			)
		}
	}()
	cleanup(ctx)
}
