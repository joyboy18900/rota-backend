package utils

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

// Validation errors
var (
	ErrEmptyUsername    = errors.New("username cannot be empty")
	ErrInvalidUsername  = errors.New("username must be 3-50 characters and contain only letters, numbers, and underscores")
	ErrEmptyEmail       = errors.New("email cannot be empty")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrEmptyPassword    = errors.New("password cannot be empty")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
)

// ValidateUsername validates a username
func ValidateUsername(username string) error {
	if username == "" {
		return ErrEmptyUsername
	}

	if len(username) < 3 || len(username) > 50 {
		return ErrInvalidUsername
	}

	// Check if username contains only letters, numbers, and underscores
	matched, err := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	if err != nil {
		return fmt.Errorf("error validating username: %w", err)
	}
	if !matched {
		return ErrInvalidUsername
	}

	return nil
}

// ValidateEmail validates an email
func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmptyEmail
	}

	// Trim whitespace
	email = strings.TrimSpace(email)

	// Check if email is valid
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}

	return nil
}

// ValidatePassword validates a password
func ValidatePassword(password string) error {
	if password == "" {
		return ErrEmptyPassword
	}

	if len(password) < 6 {
		return ErrPasswordTooShort
	}

	return nil
}
