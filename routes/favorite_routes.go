package routes

import (
	handler "rota-api/handlers"
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

func SetupFavoriteRoutes(
	app *fiber.App,
	favoriteHandler *handler.FavoriteHandler,
	authService services.AuthService,
	favoriteService services.FavoriteService,
) {
	favoriteGroup := app.Group("/api/v1/favorites")
	favoriteGroup.Use(middleware.AuthMiddleware(authService))

	favoriteGroup.Get("/", favoriteHandler.GetUserFavorites)
	favoriteGroup.Post("/stations/:stationId", favoriteHandler.AddStationToFavorites)
	favoriteGroup.Delete("/stations/:stationId", favoriteHandler.RemoveStationByStationId)
	
	favoriteGroup.Get("/admin", middleware.OwnFavoriteMiddleware(favoriteService), favoriteHandler.GetAllFavorites)
	favoriteGroup.Post("/admin", favoriteHandler.CreateFavorite)
	favoriteGroup.Get("/admin/:id", middleware.OwnFavoriteMiddleware(favoriteService), favoriteHandler.GetFavoriteByID)
	favoriteGroup.Put("/admin/:id", middleware.OwnFavoriteMiddleware(favoriteService), favoriteHandler.UpdateFavorite)
	favoriteGroup.Delete("/admin/:id", middleware.OwnFavoriteMiddleware(favoriteService), favoriteHandler.DeleteFavorite)
}
