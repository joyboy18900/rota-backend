package routes

import (
	"rota-api/handlers"
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

func SetupFavoriteRoutes(
	app *fiber.App,
	favoriteHandler *handlers.FavoriteHandler,
	authService services.AuthService,
) {
	favoriteGroup := app.Group("/api/favorites")
	favoriteGroup.Use(middleware.AuthMiddleware(authService))

	favoriteGroup.Get("/", favoriteHandler.GetFavoritesByUserID)
	favoriteGroup.Post("/", favoriteHandler.AddFavoriteByUserID)
	favoriteGroup.Get("/:id", favoriteHandler.GetFavoriteByID)
	favoriteGroup.Put("/:id", favoriteHandler.UpdateFavoriteByID)
	favoriteGroup.Delete("/:id", favoriteHandler.DeleteFavoriteByID)
}
