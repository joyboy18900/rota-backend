package routes

import (
	"rota-api/handlers" // Import as 'handler' (package name in source files)
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

func SetupFavoriteRoutes(
	app *fiber.App,
	favoriteHandler *handler.FavoriteHandler, // Use 'handler' as declared in package
	authService services.AuthService,
	favoriteService services.FavoriteService,
) {
	favoriteGroup := app.Group("/api/v1/favorites")
	favoriteGroup.Use(middleware.AuthMiddleware(authService))

	// Use OwnFavoriteMiddleware to filter favorites by user ID
	favoriteGroup.Get("/", middleware.OwnFavoriteMiddleware(favoriteService), favoriteHandler.GetAllFavorites)
	
	// Users can create their own favorites
	favoriteGroup.Post("/", favoriteHandler.CreateFavorite)
	
	// Users can only access their own favorites
	favoriteGroup.Get("/:id", middleware.OwnFavoriteMiddleware(favoriteService), favoriteHandler.GetFavoriteByID)
	favoriteGroup.Put("/:id", middleware.OwnFavoriteMiddleware(favoriteService), favoriteHandler.UpdateFavorite)
	favoriteGroup.Delete("/:id", middleware.OwnFavoriteMiddleware(favoriteService), favoriteHandler.DeleteFavorite)
}
