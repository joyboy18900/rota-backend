package routes

import (
	handler "rota-api/handlers"
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes configures all auth routes
func SetupAuthRoutes(app *fiber.App, handler *handler.AuthHandler, authService services.AuthService) {
	auth := app.Group("/api/auth")

	// Public routes
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	// TODO: Implement Google OAuth routes in the future
	// auth.Get("/google", handler.GoogleLogin)
	// auth.Get("/google/callback", handler.GoogleCallback)

	// Protected routes
	auth.Get("/me", middleware.AuthMiddleware(authService), handler.GetCurrentUser)
	auth.Post("/logout", middleware.AuthMiddleware(authService), handler.Logout)
	auth.Post("/refresh", handler.RefreshToken)
}
