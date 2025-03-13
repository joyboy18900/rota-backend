package controllers

import (
	"github.com/gofiber/fiber/v2"

	"rota-api/services"
)

// AuthController interface defines methods for authentication controller
type AuthController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	GoogleLogin(c *fiber.Ctx) error
	GoogleCallback(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// TokenResponse represents the response body for token requests
type TokenResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

// authController implements AuthController
type authController struct {
	authService services.AuthService
}

// NewAuthController creates a new authentication controller
func NewAuthController(authService services.AuthService) AuthController {
	return &authController{authService}
}

// Register handles user registration
func (c *authController) Register(ctx *fiber.Ctx) error {
	var req RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Username, email, and password are required")
	}

	user, err := c.authService.Register(ctx.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User registered successfully",
		"data": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login handles user login
func (c *authController) Login(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email and password are required")
	}

	token, err := c.authService.Login(ctx.Context(), req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return ctx.JSON(TokenResponse{
		Success: true,
		Token:   token,
	})
}

// GoogleLogin redirects to Google OAuth login page
func (c *authController) GoogleLogin(ctx *fiber.Ctx) error {
	url := c.authService.GoogleLogin()
	return ctx.Redirect(url)
}

// GoogleCallback handles the callback from Google OAuth
func (c *authController) GoogleCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	if code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Authorization code is required")
	}

	token, err := c.authService.GoogleCallback(ctx.Context(), code)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// In a real application, you might redirect to a frontend URL with the token
	return ctx.JSON(TokenResponse{
		Success: true,
		Token:   token,
	})
}

// RefreshToken refreshes a token
func (c *authController) RefreshToken(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Refresh token is required")
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	newToken, err := c.authService.RefreshToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return ctx.JSON(TokenResponse{
		Success: true,
		Token:   newToken,
	})
}
