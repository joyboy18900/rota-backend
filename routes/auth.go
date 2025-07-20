package routes

import (
	handler "rota-api/handlers"
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes configures all auth routes
func SetupAuthRoutes(app *fiber.App, authHandler *handler.AuthHandler, googleHandler *handler.GoogleOAuthHandler, authService services.AuthService) {
	auth := app.Group("/api/v1/auth")

	// Public routes
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	
	// Google OAuth routes
	auth.Get("/google", googleHandler.GoogleLogin)
	auth.Get("/google/callback", googleHandler.GoogleCallback)

	// Protected routes
	auth.Get("/me", middleware.AuthMiddleware(authService), authHandler.GetCurrentUser)
	auth.Post("/logout", middleware.AuthMiddleware(authService), authHandler.Logout)
	auth.Post("/refresh", authHandler.RefreshToken)
}
