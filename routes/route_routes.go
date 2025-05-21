package routes

import (
	"rota-api/handlers"
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

// SetupRouteRoutes sets up all route-related routes
func SetupRouteRoutes(app *fiber.App, routeHandler *handlers.RouteHandler, authService services.AuthService) {
	routes := app.Group("/api/routes")

	// Apply auth middleware
	routes.Use(middleware.AuthMiddleware(authService))

	// Route endpoints
	routes.Get("/", routeHandler.GetAllRoutes)
	routes.Get("/:id", routeHandler.GetRoute)
	routes.Post("/", routeHandler.CreateRoute)
	routes.Put("/:id", routeHandler.UpdateRoute)
	routes.Delete("/:id", routeHandler.DeleteRoute)

	// Route-Stop relationship endpoints
	routes.Post("/:id/stops", routeHandler.AddStops)
	routes.Delete("/:id/stops", routeHandler.RemoveStops)

	// Route-Vehicle relationship endpoints
	routes.Post("/:id/vehicles", routeHandler.AddVehicles)
	routes.Delete("/:id/vehicles", routeHandler.RemoveVehicles)
}
