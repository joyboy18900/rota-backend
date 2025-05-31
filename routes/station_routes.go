package routes

import (
	"rota-api/handlers"
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

// SetupStationRoutes sets up all station-related routes
func SetupStationRoutes(app *fiber.App, stationHandler *handler.StationHandler, scheduleHandler *handler.ScheduleHandler, authService services.AuthService) {
	// Create public group for routes that don't need authentication
	publicRoutes := app.Group("/api/v1/stations")
	
	// Public endpoints
	publicRoutes.Get("/", stationHandler.GetAllStations)
	// Get station schedules (both inbound and outbound)
	publicRoutes.Get("/:id/schedules", scheduleHandler.GetSchedulesByStation)

	// Protected routes for viewing
	protectedRoutes := app.Group("/api/v1/stations")
	protectedRoutes.Use(middleware.AuthMiddleware(authService))
	protectedRoutes.Get("/:id", stationHandler.GetStationByID)
	
	// Admin-only routes for station management
	adminRoutes := app.Group("/api/v1/stations")
	adminRoutes.Use(middleware.AuthMiddleware(authService), middleware.AdminMiddleware())
	adminRoutes.Post("/", stationHandler.CreateStation)
	adminRoutes.Put("/:id", stationHandler.UpdateStation)
	adminRoutes.Delete("/:id", stationHandler.DeleteStation)
}
