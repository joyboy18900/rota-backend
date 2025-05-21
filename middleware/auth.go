package middleware

import (
	"fmt"
	"strings"

	"rota-api/models"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// UserClaims represents the JWT claims for user authentication
type UserClaims struct {
	UserID int           `json:"user_id"`
	Email  string        `json:"email"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware handles JWT authentication
func AuthMiddleware(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Missing authorization header",
			})
		}

		// Check if token starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid authorization header format",
			})
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(authService.GetJWTSecret()), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired token",
			})
		}

		// Check if token is valid
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			// Check if token is blacklisted
			if authService.IsTokenBlacklisted(tokenString) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "Token has been invalidated",
				})
			}

			// Set user data in context
			c.Locals("userID", claims.UserID)
			c.Locals("userEmail", claims.Email)
			c.Locals("userRole", claims.Role)

			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid token",
		})
	}
}

// RoleMiddleware creates a middleware that checks if the user has the required role
func RoleMiddleware(roles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context
		userRole, ok := c.Locals("userRole").(models.UserRole)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "unauthorized: missing user role",
			})
		}

		// Check if user has any of the required roles
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		// If user is admin, always allow access
		if userRole == models.RoleAdmin {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "forbidden: insufficient permissions",
		})
	}
}

// StaffMiddleware checks if the user is a staff member or admin
func StaffMiddleware() fiber.Handler {
	return RoleMiddleware(models.RoleStaff, models.RoleAdmin)
}

// AdminMiddleware checks if the user is an admin
func AdminMiddleware() fiber.Handler {
	return RoleMiddleware(models.RoleAdmin)
}

// GetUserIDFromContext retrieves the user ID from the context
func GetUserIDFromContext(c *fiber.Ctx) (int, error) {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "missing user ID in context")
	}
	return userID, nil
}

// GetUserRoleFromContext retrieves the user role from the context
func GetUserRoleFromContext(c *fiber.Ctx) (models.UserRole, error) {
	role, ok := c.Locals("userRole").(models.UserRole)
	if !ok {
		return "", fiber.NewError(fiber.StatusUnauthorized, "missing user role in context")
	}
	return role, nil
}
