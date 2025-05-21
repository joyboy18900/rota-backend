package routes

import (
	"rota-api/handlers"
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

// SetupRouteRoutes sets up all route-related routes
func SetupRouteRoutes(app *fiber.App, routeHandler *handler.RouteHandler, authService services.AuthService) {
	// Create public group for routes that don't need authentication
	publicRoutes := app.Group("/api/v1/routes")
	
	// Public endpoints
	publicRoutes.Get("/", routeHandler.GetAllRoutes)

	// Protected routes
	protectedRoutes := app.Group("/api/v1/routes")
	protectedRoutes.Use(middleware.AuthMiddleware(authService))
	protectedRoutes.Get("/:id", routeHandler.GetRouteByID)
	protectedRoutes.Post("/", routeHandler.CreateRoute)
	protectedRoutes.Put("/:id", routeHandler.UpdateRoute)
	protectedRoutes.Delete("/:id", routeHandler.DeleteRoute)

	// Route-Stop relationship endpoints
	// TODO: Implement these methods in RouteHandler
	// routes.Post("/:id/stops", routeHandler.AddStops)
	// routes.Delete("/:id/stops", routeHandler.RemoveStops)

	// Route-Vehicle relationship endpoints
	// TODO: Implement these methods in RouteHandler
	// routes.Post("/:id/vehicles", routeHandler.AddVehicles)
	// routes.Delete("/:id/vehicles", routeHandler.RemoveVehicles)
}
