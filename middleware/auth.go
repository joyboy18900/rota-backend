package middleware

import (
	"strings"

	"rota-api/repositories"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware handles authentication
func AuthMiddleware(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		// Check if token starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token format",
			})
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Check if token is blacklisted
		if authService.IsTokenBlacklisted(token) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token has been invalidated",
			})
		}

		// Validate token and get user ID
		claims, err := authService.ValidateAccessToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		// Set user ID in context
		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}

// AdminMiddleware handles admin authorization
func AdminMiddleware(userRepo repositories.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from context
		userID := c.Locals("userID").(string)

		// Get user from database
		user, err := userRepo.FindByID(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "user not found",
			})
		}

		// Check if user is admin
		if user.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "admin access required",
			})
		}

		return c.Next()
	}
}
