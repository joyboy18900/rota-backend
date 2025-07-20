package handler

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"rota-api/services"
	"rota-api/utils"
)

type GoogleOAuthHandler struct {
	authService   services.AuthService
	googleOAuth   services.GoogleOAuthService
}

func NewGoogleOAuthHandler(authService services.AuthService, googleOAuth services.GoogleOAuthService) *GoogleOAuthHandler {
	return &GoogleOAuthHandler{
		authService: authService,
		googleOAuth: googleOAuth,
	}
}

func (h *GoogleOAuthHandler) GoogleLogin(c *fiber.Ctx) error {
	if !h.isValidOrigin(c) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Invalid origin")
	}
	
	// Get redirect_uri from query parameter, fallback to default if not provided
	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		// Use default redirect URI from environment
		redirectURI = os.Getenv("GOOGLE_REDIRECT_URL")
	}
	
	state := h.googleOAuth.GenerateState()
	authURL := h.googleOAuth.GetAuthURLWithRedirect(state, redirectURI)
	
	return c.JSON(fiber.Map{
		"auth_url": authURL,
		"state":    state,
	})
}

func (h *GoogleOAuthHandler) GoogleCallback(c *fiber.Ctx) error {
	// Validate request origin for security
	if !h.isValidOrigin(c) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Invalid origin")
	}
	
	// Get required parameters from request
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")
	
	// Check for errors returned from Google OAuth
	if errorParam != "" {
		// Log and return clear error message
		logError := fmt.Sprintf("Google OAuth error: %s", errorParam)
		fmt.Println(logError)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "OAuth error: "+errorParam)
	}
	
	// Verify authorization code is provided
	if code == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Authorization code not provided")
	}
	
	// Verify state parameter is provided
	if state == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "State parameter not provided")
	}
	
	// Validate state to prevent CSRF attacks
	if !h.googleOAuth.ValidateState(state) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid or expired state parameter")
	}
	
	// Determine redirect_uri used for this callback
	// If the request came from frontend, use frontend callback URL
	// Otherwise, use the default backend callback URL
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URL") // default
	referer := c.Get("Referer")
	if strings.Contains(referer, "localhost:3000") || strings.Contains(referer, os.Getenv("FRONTEND_URL")) {
		// This callback was initiated from frontend
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}
		redirectURI = frontendURL + "/auth/callback"
	}
	
	// Process Google login with correct redirect URI
	user, accessToken, err := h.authService.LoginWithGoogleAndRedirect(c.Context(), code, state, redirectURI)
	if err != nil {
		// Log and return authentication error
		fmt.Printf("Error authenticating with Google: %v\n", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to authenticate with Google: "+err.Error())
	}
	
	// Generate refresh token
	refreshToken, err := h.authService.GenerateRefreshToken()
	if err != nil {
		fmt.Printf("Error generating refresh token: %v\n", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate refresh token")
	}
	
	// Return user data and tokens
	return utils.SuccessResponse(c, fiber.StatusOK, "Google login successful", fiber.Map{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    3600, // 1 hour by default
	})
}

func (h *GoogleOAuthHandler) isValidOrigin(c *fiber.Ctx) bool {
	env := os.Getenv("APP_ENV")
	if env == "development" || env == "" {
		return true
	}
	
	origin := c.Get("Origin")
	referer := c.Get("Referer")
	
	// Check allowed origins from environment variable
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	allowedOriginsList := strings.Split(allowedOriginsEnv, ",")
	
	// Add FRONTEND_URL if configured
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL != "" {
		allowedOriginsList = append(allowedOriginsList, frontendURL)
	}
	
	// Check both Origin and Referer headers
	for _, allowed := range allowedOriginsList {
		if allowed == "" {
			continue
		}
		
		// Validate Origin header
		if origin == allowed {
			return true
		}
		
		// Validate Referer header
		if referer != "" && strings.HasPrefix(referer, allowed) {
			return true
		}
	}
	
	return false
}
