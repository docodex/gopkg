package internal_test

import (
	"errors"
	"testing"

	"github.com/docodex/gopkg/internal"
	"github.com/stretchr/testify/assert"
)

func TestWrapPanic_StringValue(t *testing.T) {
	err := internal.WrapPanic("boom")
	assert.Error(t, err)
	assert.Equal(t, "panic: boom", err.Error())
}

func TestWrapPanic_IntValue(t *testing.T) {
	err := internal.WrapPanic(42)
	assert.Error(t, err)
	assert.Equal(t, "panic: 42", err.Error())
}

func TestWrapPanic_ErrorPreservesIdentity(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := internal.WrapPanic(sentinel)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, sentinel),
		"errors.Is must find the original error though %w wrapping")
	assert.Equal(t, "panic: sentinel", err.Error())
}

func TestWrapPanic_NilValue(t *testing.T) {
	// Defensive: WrapPanic should never be called with a nil value
	// (the only meaningful caller is `if r := recover(); r != nil`),
	// but if it is, the result must still be a usable error.
	err := internal.WrapPanic(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic:")
}

type myErrorWithCode struct {
	code int
}

func (e *myErrorWithCode) Error() string {
	return "my error"
}

func TestWrapPanic_CustomErrorType(t *testing.T) {
	myErr := &myErrorWithCode{code: 7}
	err := internal.WrapPanic(myErr)

	var target *myErrorWithCode
	assert.True(t, errors.As(err, &target),
		"errors.As must unwrap to the original concrete error type")
	assert.Equal(t, 7, target.code)
}
