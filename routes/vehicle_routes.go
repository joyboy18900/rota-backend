package routes

import (
	"rota-api/handlers"
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

// SetupVehicleRoutes sets up all vehicle-related routes
func SetupVehicleRoutes(
	app *fiber.App,
	vehicleHandler *handler.VehicleHandler,
	authService services.AuthService,
) {
	vehicles := app.Group("/api/v1/vehicles")

	// Routes for viewing vehicles (available to all authenticated users)
	vehicles.Use(middleware.AuthMiddleware(authService))
	vehicles.Get("/", vehicleHandler.GetAllVehicles)
	vehicles.Get("/:id", vehicleHandler.GetVehicleByID)

	// Admin-only routes for vehicle management
	adminVehicles := app.Group("/api/v1/vehicles")
	adminVehicles.Use(middleware.AuthMiddleware(authService), middleware.AdminMiddleware())
	adminVehicles.Post("/", vehicleHandler.CreateVehicle)
	adminVehicles.Put("/:id", vehicleHandler.UpdateVehicle)
	adminVehicles.Delete("/:id", vehicleHandler.DeleteVehicle)
}
