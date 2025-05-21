package handler

import (
	"rota-api/models"
	"rota-api/response"
	"rota-api/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles HTTP requests related to authentication
type AuthHandler struct {
	authService services.AuthService
	validate    *validator.Validate
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=50,alphanum"`
	Email    string  `json:"email" validate:"required,email,max=100"`
	Password string  `json:"password" validate:"required,min=8,containsany=!@#$%^&*()_+=-,contains=ABCDEFGHIJKLMNOPQRSTUVWXYZ,contains=abcdefghijklmnopqrstuvwxyz,contains=0123456789"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest represents the request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration details"
// @Success 201 {object} response.SuccessResponse{data=AuthResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return response.ValidationError(c, err)
	}

	// Create user model
	user := &models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: &req.Password,
	}

	// Register user
	user, err := h.authService.Register(c, user)
	if err != nil {
		switch err.Error() {
		case "email already exists":
			return response.Conflict(c, "Email already exists")
		case "username already exists":
			return response.Conflict(c, "Username already exists")
		default:
			return response.InternalServerError(c, "Failed to register user")
		}
	}

	// Generate JWT token
	token, err := h.authService.GenerateAccessToken(user)
	if err != nil {
		return response.InternalServerError(c, "Failed to generate token")
	}

	// Generate refresh token
	refreshToken, err := h.authService.GenerateRefreshToken()
	if err != nil {
		return response.InternalServerError(c, "Failed to generate refresh token")
	}

	// Clear sensitive data
	user.Password = nil

	// Return success response with tokens
	return response.Created(c, AuthResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
	}, "User registered successfully")
}

// Login handles user login
// @Summary Login a user
// @Description Authenticate user and return access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User credentials"
// @Success 200 {object} response.SuccessResponse{data=AuthResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 429 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return response.ValidationError(c, err)
	}

	// Check rate limiting (pseudo-code, implement your rate limiter)
	// if rateLimited := checkRateLimit(c.IP()); rateLimited {
	// 	return response.TooManyRequests(c, "Too many login attempts. Please try again later.")
	// }


	// Authenticate user
	user, token, err := h.authService.Login(c, req.Email, req.Password)
	if err != nil {
		// Log failed login attempt
		// incrementFailedAttempts(c.IP())
		return response.Unauthorized(c, "Invalid email or password")
	}

	// Generate refresh token with expiration
	expiresIn := int(h.authService.GetJWTExpiration().Seconds())
	
	// Generate access token with expiration claim
	token, err = h.authService.GenerateAccessToken(user)
	if err != nil {
		return response.InternalServerError(c, "Failed to generate access token")
	}

	// Generate refresh token
	refreshToken, err := h.authService.GenerateRefreshToken()
	if err != nil {
		return response.InternalServerError(c, "Failed to generate refresh token")
	}

	// TODO: Store refresh token in database with user ID and expiration
	// if err := h.tokenRepo.StoreRefreshToken(user.ID, refreshToken, h.authService.GetRefreshExpiration()); err != nil {
	// 	return response.InternalServerError(c, "Failed to store refresh token")
	// }

	// Clear sensitive data
	user.Password = nil

	// Return tokens
	return response.Success(c, AuthResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}, "Login successful")
}

// GetCurrentUser returns the current authenticated user's profile
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=models.User}
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(int)
	if !ok || userID == 0 {
		return response.Unauthorized(c, "Unauthorized")
	}

	// Get user from database
	user, err := h.authService.GetUserByID(c.Context(), userID)
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	// Clear sensitive data
	user.Password = nil

	return response.Success(c, user, "User retrieved successfully")
}

// RefreshToken handles access token refresh
// @Summary Refresh access token
// @Description Get a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.SuccessResponse{data=AuthResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return response.ValidationError(c, err)
	}

	// Check if refresh token is valid and not revoked
	// userID, err := h.tokenRepo.ValidateRefreshToken(req.RefreshToken)
	// if err != nil {
	// 	return response.Unauthorized(c, "Invalid or expired refresh token")
	// }


	// Get user from database
	// user, err := h.authService.GetUserByID(c.Context(), userID)
	// if err != nil {
	// 	return response.NotFound(c, "User not found")
	// }


	// Generate new access token
	// token, err := h.authService.GenerateAccessToken(user)
	// if err != nil {
	// 	return response.InternalServerError(c, "Failed to generate access token")
	// }


	// Generate new refresh token (optional: implement refresh token rotation)
	// newRefreshToken, err := h.authService.GenerateRefreshToken()
	// if err != nil {
	// 	return response.InternalServerError(c, "Failed to generate refresh token")
	// }


	// Invalidate old refresh token
	// if err := h.tokenRepo.RevokeRefreshToken(req.RefreshToken); err != nil {
	// 	// Log error but continue
	// 	log.Printf("Failed to revoke refresh token: %v", err)
	// }


	// Store new refresh token
	// if err := h.tokenRepo.StoreRefreshToken(user.ID, newRefreshToken, h.authService.GetRefreshExpiration()); err != nil {
	// 	return response.InternalServerError(c, "Failed to store refresh token")
	// }


	// Return new tokens
	return response.Success(c, AuthResponse{
		AccessToken:  "new_access_token", // Replace with actual token
		RefreshToken: "new_refresh_token", // Replace with actual refresh token
		TokenType:    "Bearer",
		ExpiresIn:    int(h.authService.GetJWTExpiration().Seconds()),
	}, "Token refreshed successfully")
}

// Logout handles user logout
// @Summary Logout user
// @Description Invalidate the current access token and refresh token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return response.BadRequest(c, "Missing authorization header")
	}

	// Check if token starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return response.BadRequest(c, "Invalid authorization header format")
	}

	// Extract token (without "Bearer " prefix)
	token := authHeader[7:]

	// Get refresh token from request body if provided
	var refreshToken string
	if c.Get("Content-Type") == "application/json" {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BodyParser(&body); err == nil && body.RefreshToken != "" {
			refreshToken = body.RefreshToken
		}
	}

	// Invalidate the access token
	if err := h.authService.Logout(token); err != nil {
		return response.InternalServerError(c, "Failed to invalidate access token")
	}

	// If refresh token is provided, invalidate it as well
	if refreshToken != "" {
		// In a real implementation, you would also invalidate the refresh token
		// if err := h.tokenRepo.RevokeRefreshToken(refreshToken); err != nil {
		// 	// Log error but don't fail the request
		// 	log.Printf("Failed to revoke refresh token: %v", err)
		// }
	}

	return c.SendStatus(fiber.StatusNoContent)
}

/* TODO: Implement Google OAuth handlers in the future
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
{{ ... }}
	return c.JSON(fiber.Map{
		"url": url,
	})
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "code is required",
		})
	}

	token, err := h.authService.GoogleCallback(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
*/
