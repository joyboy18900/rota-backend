
package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ErrorResponse defines the structure for error responses
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ErrorHandler is a custom error handler for Fiber
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default status code is 500
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	errorDetail := ""

	// Retrieve the custom status code if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	} else {
		// Handle specific error types
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			code = fiber.StatusNotFound
			message = "Resource not found"
		case strings.Contains(err.Error(), "duplicate key"):
			code = fiber.StatusConflict
			message = "Resource already exists"
		}

		// Log non-fiber errors
		if code == fiber.StatusInternalServerError {
			log.Printf("Internal error: %v", err)
			// In production, don't expose the error detail
			if getEnv("ENVIRONMENT", "development") != "production" {
				errorDetail = err.Error()
			}
		} else {
			errorDetail = err.Error()
		}
	}

	// Return JSON response with error details
	return c.Status(code).JSON(ErrorResponse{
		Success: false,
		Message: message,
		Error:   errorDetail,
	})
}

// NewInternalError creates a new internal server error with the given message
func NewInternalError(format string, args ...interface{}) error {
	return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf(format, args...))
}

// NewNotFoundError creates a new not found error with the given message
func NewNotFoundError(format string, args ...interface{}) error {
	return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf(format, args...))
}

// NewBadRequestError creates a new bad request error with the given message
func NewBadRequestError(format string, args ...interface{}) error {
	return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf(format, args...))
}

// NewUnauthorizedError creates a new unauthorized error with the given message
func NewUnauthorizedError(format string, args ...interface{}) error {
	return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf(format, args...))
}

// NewForbiddenError creates a new forbidden error with the given message
func NewForbiddenError(format string, args ...interface{}) error {
	return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf(format, args...))
}
