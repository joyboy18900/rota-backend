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
) {
	favoriteGroup := app.Group("/api/v1/favorites")
	favoriteGroup.Use(middleware.AuthMiddleware(authService))

	favoriteGroup.Get("/", favoriteHandler.GetAllFavorites)
	favoriteGroup.Post("/", favoriteHandler.CreateFavorite)
	favoriteGroup.Get("/:id", favoriteHandler.GetFavoriteByID)
	favoriteGroup.Put("/:id", favoriteHandler.UpdateFavorite)
	favoriteGroup.Delete("/:id", favoriteHandler.DeleteFavorite)
}
