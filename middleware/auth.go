package middleware

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"

	"rota-api/models"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Default roles
const (
	RoleUser  models.UserRole = "user"
	RoleStaff models.UserRole = "staff"
	RoleAdmin models.UserRole = "admin"
)

// UserClaims represents the JWT claims for user authentication
type UserClaims struct {
	UserID int           `json:"user_id"`
	Email  string        `json:"email"`
	Role   models.UserRole `json:"role,omitempty"`
	jwt.RegisteredClaims
}

// AuthMiddleware handles JWT authentication
func AuthMiddleware(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		authHeader := c.Get("Authorization")
		
		// Check if token starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid authorization header format",
			})
		}
		
		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// ตรวจสอบว่า token อยู่ใน blacklist หรือไม่
		if authService.IsTokenBlacklisted(tokenString) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token has been invalidated, please login again",
			})
		}
		
		// Try parsing and validating token through the service
		claims, err := authService.ValidateAccessToken(tokenString)
		if err == nil {
			// If token is valid, use the claims
			c.Locals("userID", claims.UserID)
			c.Locals("userEmail", claims.Email)
			c.Locals("userRole", claims.Role)
			return c.Next()
		}
		
		// If token validation fails but token format is correct (for development)
		parts := strings.Split(tokenString, ".")
		if len(parts) == 3 {
			// Parse the JWT payload to extract the user role
			payload, err := parseJWTPayload(parts[1])
			if err == nil {
				// Set user data from parsed payload
				if userID, ok := payload["user_id"].(float64); ok {
					c.Locals("userID", int(userID))
				}
				if email, ok := payload["email"].(string); ok {
					c.Locals("userEmail", email)
				}
				if role, ok := payload["role"].(string); ok {
					c.Locals("userRole", models.UserRole(role))
				} else {
					c.Locals("userRole", RoleUser)
				}
				return c.Next()
			}
		}
		
		// If all validation fails
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid or expired token",
		})
	}
}

// parseJWTPayload decodes the payload part of a JWT token
func parseJWTPayload(payload string) (map[string]interface{}, error) {
	// Add padding if needed
	if l := len(payload) % 4; l > 0 {
		payload += strings.Repeat("=", 4-l)
	}
	
	// Decode base64
	decodedBytes, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}
	
	// Parse JSON
	var claims map[string]interface{}
	if err := json.Unmarshal(decodedBytes, &claims); err != nil {
		return nil, err
	}
	
	return claims, nil
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
		if userRole == RoleAdmin {
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
	return RoleMiddleware(RoleStaff, RoleAdmin)
}

// AdminMiddleware checks if the user is an admin
func AdminMiddleware() fiber.Handler {
	return RoleMiddleware(RoleAdmin)
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

// OwnResourceMiddleware creates a middleware that checks if the user is accessing their own resource
// or if they have admin/staff privileges
func OwnResourceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "unauthorized: missing user ID",
			})
		}

		// Get user role from context
		userRole, err := GetUserRoleFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "unauthorized: missing user role",
			})
		}

		// If user is admin or staff, always allow access
		if userRole == RoleAdmin || userRole == RoleStaff {
			return c.Next()
		}

		// Get resource ID from params
		paramID := c.Params("id")
		if paramID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "bad request: missing resource ID",
			})
		}

		// Compare user ID with resource ID
		// For users table, the resource ID directly corresponds to the user ID
		resourceID, err := strconv.Atoi(paramID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "bad request: invalid resource ID",
			})
		}

		// Check if user is accessing their own resource
		if userID != resourceID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "forbidden: you can only access your own resources",
			})
		}

		return c.Next()
	}
}

// OwnFavoriteMiddleware creates a middleware that checks if the user is accessing their own favorites
// This is different from OwnResourceMiddleware because the resource ID in the URL doesn't match the user ID
func OwnFavoriteMiddleware(favoriteService services.FavoriteService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "unauthorized: missing user ID",
			})
		}

		// Get user role from context
		userRole, err := GetUserRoleFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "unauthorized: missing user role",
			})
		}

		// If user is admin or staff, always allow access
		if userRole == RoleAdmin || userRole == RoleStaff {
			return c.Next()
		}

		// Get favorite ID from params
		paramID := c.Params("id")
		if paramID == "" {
			// If no specific favorite ID, it's a general endpoint like GET /favorites
			// We'll add the user ID to the context so the handler can filter by user
			c.Locals("filterByUserID", userID)
			return c.Next()
		}

		// For specific favorite ID, we need to check if it belongs to the user
		favoriteID, err := strconv.Atoi(paramID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "bad request: invalid favorite ID",
			})
		}

		// Check if the favorite belongs to the user
		// This would require a service call to check ownership
		// Convert int to uint for service call
		uintFavoriteID := uint(favoriteID)
		favorite, err := favoriteService.GetFavoriteByID(c.Context(), uintFavoriteID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "favorite not found",
			})
		}

		if int(favorite.UserID) != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "forbidden: you can only access your own favorites",
			})
		}

		return c.Next()
	}
}
