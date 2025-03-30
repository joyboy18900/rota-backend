package errors

import "fmt"

// Error represents a custom error type
type Error struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap implements the unwrap interface
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError creates a new error
func NewError(code string, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
