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

	// Apply auth middleware
	vehicles.Use(middleware.AuthMiddleware(authService))

	// Vehicle endpoints
	vehicles.Get("/", vehicleHandler.GetAllVehicles)
	vehicles.Get("/:id", vehicleHandler.GetVehicleByID)
	vehicles.Post("/", vehicleHandler.CreateVehicle)
	vehicles.Put("/:id", vehicleHandler.UpdateVehicle)
	vehicles.Delete("/:id", vehicleHandler.DeleteVehicle)
}
