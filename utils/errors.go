
package utils

import (
	"errors"
	"fmt"
)

// Common errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrDuplicateEntry    = errors.New("duplicate entry")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrPermissionDenied  = errors.New("permission denied")
)

// AuthError represents an authentication error
type AuthError struct {
	Message string
	Err     error
}

// Error implements the error interface
func (e *AuthError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AuthError) Unwrap() error {
	return e.Err
}

// NewAuthError creates a new AuthError
func NewAuthError(message string, err error) *AuthError {
	return &AuthError{
		Message: message,
		Err:     err,
	}
}
