// Package e provides a standardized error handling mechanism for the Core library.
// It allows for wrapping errors with contextual information, making it easier to
// trace the origin of an error and provide meaningful feedback.
//
// The design of this package is influenced by the need for a simple, yet powerful
// way to handle errors that can occur in different layers of the application,
// from low-level file operations to high-level service interactions.
//
// The key features of this package are:
//   - Error wrapping: The Op and an optional Msg field provide context about
//     where and why an error occurred.
//   - Stack traces: By wrapping errors, we can build a logical stack trace
//     that is more informative than a raw stack trace.
//   - Consistent error handling: Encourages a uniform approach to error
//     handling across the entire codebase.
package e

import (
	"fmt"
)

// Error represents a standardized error with operational context.
type Error struct {
	// Op is the operation being performed, e.g., "config.Load".
	Op string
	// Msg is a human-readable message explaining the error.
	Msg string
	// Err is the underlying error that was wrapped.
	Err error
}

// E is a helper function to create a new Error.
// This is the primary way to create errors that will be consumed by the system.
// For example:
//
//	return e.E("config.Load", "failed to load config file", err)
//
// The 'op' parameter should be in the format of 'package.function' or 'service.method'.
// The 'msg' parameter should be a human-readable message that can be displayed to the user.
// The 'err' parameter is the underlying error that is being wrapped.
func E(op, msg string, err error) error {
	if err == nil {
		return &Error{Op: op, Msg: msg}
	}
	return &Error{Op: op, Msg: msg, Err: err}
}

// Error returns the string representation of the error.
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Msg, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Msg)
}

// Unwrap provides compatibility for Go's errors.Is and errors.As functions.
func (e *Error) Unwrap() error {
	return e.Err
}
