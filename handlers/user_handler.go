package handler

import (
	"rota-api/models"
	"rota-api/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService services.UserService
	authService services.AuthService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService services.UserService, authService services.AuthService) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
	}
}

// GetAllUsers retrieves all users (admin only)
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch users",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    users,
	})
}

// GetUserByID retrieves a specific user by ID
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	user, err := h.authService.GetUserByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// UpdateUser updates a user's information
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	// Get existing user
	user, err := h.authService.GetUserByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	// Parse updated user data
	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
	}

	// Update user data
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Role != "" {
		user.Role = updateData.Role
	}
	if updateData.Username != nil {
		user.Username = updateData.Username
	}
	if updateData.ProfilePicture != nil {
		user.ProfilePicture = updateData.ProfilePicture
	}
	if updateData.Password != nil && *updateData.Password != "" {
		user.Password = updateData.Password
	}

	// Save updated user
	if err := h.userService.UpdateUser(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	if err := h.userService.DeleteUser(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User deleted successfully",
	})
}

// CreateUser creates a new user (admin only)
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var newUser models.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if newUser.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email is required",
		})
	}

	if newUser.Password == nil || *newUser.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Password is required",
		})
	}

	// Set default role if not provided
	if newUser.Role == "" {
		newUser.Role = models.RoleUser
	}

	// Register the user through auth service
	createdUser, err := h.authService.Register(c, &newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User created successfully",
		"data":    createdUser,
	})
}

// UpdateUserRole updates a user's role (admin only)
func (h *UserHandler) UpdateUserRole(c *fiber.Ctx) error {
	userID := c.Params("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	// Get existing user
	user, err := h.authService.GetUserByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	// Parse role update data
	var roleData struct {
		Role models.UserRole `json:"role"`
	}
	if err := c.BodyParser(&roleData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
	}

	// Validate role
	if roleData.Role != models.RoleUser && roleData.Role != models.RoleStaff && roleData.Role != models.RoleAdmin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid role value",
		})
	}

	// Update user role
	user.Role = roleData.Role
	if err := h.userService.UpdateUser(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update user role",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User role updated successfully",
		"data": fiber.Map{
			"id":   user.ID,
			"role": user.Role,
		},
	})
}
