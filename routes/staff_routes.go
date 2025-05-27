package routes

import (
	"rota-api/handlers"
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

// SetupStaffRoutes sets up all staff-related routes
func SetupStaffRoutes(app *fiber.App, staffHandler *handler.StaffHandler, authService services.AuthService) {
	// Create public group for routes that don't need authentication
	publicRoutes := app.Group("/api/v1/staff")
	
	// Public endpoints
	publicRoutes.Get("/", staffHandler.GetAllStaff)

	// Protected routes
	protectedRoutes := app.Group("/api/v1/staff")
	protectedRoutes.Use(middleware.AuthMiddleware(authService))
	protectedRoutes.Get("/:id", staffHandler.GetStaffByID)
	
	// Admin-only routes for staff management
	adminRoutes := app.Group("/api/v1/staff")
	adminRoutes.Use(middleware.AuthMiddleware(authService), middleware.AdminMiddleware())
	adminRoutes.Post("/", staffHandler.CreateStaff)
	adminRoutes.Put("/:id", staffHandler.UpdateStaff)
	adminRoutes.Delete("/:id", staffHandler.DeleteStaff)
}
