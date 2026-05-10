package internal

import "fmt"

// WrapPanic turns a recover() value into an error. If the panic value
// is already an error, it is wrapped with %w so errors.Is / errors.As
// continue to work; otherwise it is formatted with %v.
//
// This is the shared panic-to-error conversion for packages in this
// module so they all produce identical "panic: <msg>" text and
// preserve error identity uniformly.
func WrapPanic(r any) error {
	if err, ok := r.(error); ok {
		return fmt.Errorf("panic: %w", err)
	}
	return fmt.Errorf("panic: %v", r)
}
