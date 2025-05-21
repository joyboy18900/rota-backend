package response

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Response represents a standard API response
type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Success sends a successful JSON response with status code 200
func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response[interface{}]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a success response with status code 201
func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(Response[interface{}]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response[interface{}]{
		Success: false,
		Message: "Bad Request",
		Error:   message,
	})
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response[interface{}]{
		Success: false,
		Message: "Unauthorized",
		Error:   message,
	})
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(Response[interface{}]{
		Success: false,
		Message: "Forbidden",
		Error:   message,
	})
}

// NotFound sends a 404 Not Found response
func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response[interface{}]{
		Success: false,
		Message: "Not Found",
		Error:   message,
	})
}

// Conflict sends a 409 Conflict response
func Conflict(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(Response[interface{}]{
		Success: false,
		Message: "Conflict",
		Error:   message,
	})
}

// ValidationError sends a 422 Unprocessable Entity response with validation errors
func ValidationError(c *fiber.Ctx, err error) error {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return BadRequest(c, err.Error())
	}

	// Convert validation errors to a map
	errors := make(map[string]string)
	for _, e := range errs {
		errors[e.Field()] = e.Tag()
	}

	return c.Status(fiber.StatusUnprocessableEntity).JSON(Response[map[string]string]{
		Success: false,
		Message: "Validation failed",
		Error:   "One or more validation errors occurred",
		Data:    errors,
	})
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response[interface{}]{
		Success: false,
		Message: "Internal Server Error",
		Error:   message,
	})
}
