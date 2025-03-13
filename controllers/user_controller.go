package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"rota-api/services"
)

// UserController interface defines methods for user controller
type UserController interface {
	GetUserByID(c *fiber.Ctx) error
	GetAllUsers(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// userController implements UserController
type userController struct {
	userService services.UserService
}

// NewUserController creates a new user controller
func NewUserController(userService services.UserService) UserController {
	return &userController{userService}
}

// GetUserByID retrieves a user by ID
func (c *userController) GetUserByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := c.userService.GetUserByID(ctx.Context(), uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// GetAllUsers retrieves all users
func (c *userController) GetAllUsers(ctx *fiber.Ctx) error {
	users, err := c.userService.GetAllUsers(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    users,
	})
}

// UpdateUser updates a user
func (c *userController) UpdateUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	var req UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Get current user ID from context
	currentUserID, ok := ctx.Locals("userID").(uint)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Check if user is updating their own profile
	if currentUserID != uint(id) {
		return fiber.NewError(fiber.StatusForbidden, "You can only update your own profile")
	}

	user, err := c.userService.UpdateUser(ctx.Context(), uint(id), req.Username, req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser deletes a user
func (c *userController) DeleteUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	// Get current user ID from context
	currentUserID, ok := ctx.Locals("userID").(uint)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Check if user is deleting their own account
	if currentUserID != uint(id) {
		return fiber.NewError(fiber.StatusForbidden, "You can only delete your own account")
	}

	if err := c.userService.DeleteUser(ctx.Context(), uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "User deleted successfully",
	})
}
