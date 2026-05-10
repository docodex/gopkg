package recovery_test

import (
	"context"
	"errors"
	"testing"

	"github.com/docodex/gopkg/recovery"
	"github.com/stretchr/testify/assert"
)

func TestRecover_WritesError(t *testing.T) {
	var err error
	func() {
		defer recovery.Recover(context.Background(), &err)
		panic("boom")
	}()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic: boom")
}

func TestRecover_PreservesErrorIdentity(t *testing.T) {
	sentinel := errors.New("sentinel")
	var err error
	func() {
		defer recovery.Recover(context.Background(), &err)
		panic(sentinel)
	}()

	assert.Error(t, err)
	assert.True(t, errors.Is(err, sentinel),
		"errors.Is must find the original panic value through the wrapping")
	assert.Contains(t, err.Error(), "panic: sentinel")
}

func TestRecover_PanicNil(t *testing.T) {
	// Go 1.21+ converts panic(nil) into *runtime.PanicNilError, so
	// recover() returns a non-nil value and Recover must still log
	// and produce a non-nil *err.
	var err error
	func() {
		defer recovery.Recover(context.Background(), &err)
		panic(nil)
	}()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic:")
}

func TestRecover_OverwritesPreviousErr(t *testing.T) {
	// Documented contract: when a panic is in flight, *err is
	// overwritten with the recovered error even if the surrounding
	// function had already set it.
	preset := errors.New("preset")
	err := preset
	func() {
		defer recovery.Recover(context.Background(), &err)
		panic("boom")
	}()

	assert.NotSame(t, preset, err, "*err must be overwritten on panic")
	assert.Contains(t, err.Error(), "panic: boom")
}

func TestRecover_NilErrPointer(t *testing.T) {
	// Passing nil for the error pointer must not panic; the recover
	// path should still log and run callbacks.
	called := false
	func() {
		defer recovery.Recover(context.Background(), nil, func(ctx context.Context) {
			called = true
		})
		panic("boom")
	}()

	assert.True(t, called)
}

func TestRecover_NoPanicIsNoop(t *testing.T) {
	preset := errors.New("preset")
	err := preset
	called := false

	func() {
		defer recovery.Recover(context.Background(), &err, func(ctx context.Context) {
			called = true
		})
		// no panic
	}()

	assert.Same(t, preset, err, "Recover must not touch *err when no panic is in flight")
	assert.False(t, called, "callbacks must not run when no panic is in flight")
}

func TestRecover_MultipleCallbacksInOrder(t *testing.T) {
	var order []int
	func() {
		defer recovery.Recover(context.Background(), nil,
			func(ctx context.Context) {
				order = append(order, 1)
			},
			nil, // skipped silently
			func(ctx context.Context) {
				order = append(order, 2)
			}, func(ctx context.Context) {
				order = append(order, 3)
			},
		)
		panic("boom")
	}()

	assert.Equal(t, []int{1, 2, 3}, order)
}

func TestRecover_PanickingCallbackIsIsolated(t *testing.T) {
	// A panicking callback must not stop subsequent callbacks nor
	// leak out of the recovery package.
	var ran []int
	assert.NotPanics(t, func() {
		func() {
			defer recovery.Recover(context.Background(), nil,
				func(ctx context.Context) {
					ran = append(ran, 1)
				}, func(ctx context.Context) {
					panic("callback boom")
				}, func(ctx context.Context) {
					ran = append(ran, 3)
				},
			)
			panic("boom")
		}()
	})

	assert.Equal(t, []int{1, 3}, ran)
}

func TestRecovery_SwallowsPanic(t *testing.T) {
	called := false
	assert.NotPanics(t, func() {
		func() {
			defer recovery.Recovery(context.Background(), func(ctx context.Context) {
				called = true
			})
			panic("boom")
		}()
	})

	assert.True(t, called)
}

func TestRecovery_NoPanicIsNoop(t *testing.T) {
	called := false
	func() {
		defer recovery.Recovery(context.Background(), func(ctx context.Context) {
			called = true
		})
		// no panic
	}()

	assert.False(t, called)
}
