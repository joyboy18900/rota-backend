package routes

import (
	"rota-api/handlers"
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes configures all user management routes
func SetupUserRoutes(
	app *fiber.App,
	userHandler *handler.UserHandler,
	authService services.AuthService,
) {
	userGroup := app.Group("/api/v1/users")
	userGroup.Use(middleware.AuthMiddleware(authService))

	// Routes for admins only
	userGroup.Get("/", middleware.AdminMiddleware(), userHandler.GetAllUsers)
	userGroup.Post("/", middleware.AdminMiddleware(), userHandler.CreateUser)
	userGroup.Delete("/:id", middleware.AdminMiddleware(), userHandler.DeleteUser)
	userGroup.Put("/:id/role", middleware.AdminMiddleware(), userHandler.UpdateUserRole)

	// Routes for users to manage their own profile or for admins
	userGroup.Get("/:id", middleware.OwnResourceMiddleware(), userHandler.GetUserByID)
	userGroup.Put("/:id", middleware.OwnResourceMiddleware(), userHandler.UpdateUser)
}
