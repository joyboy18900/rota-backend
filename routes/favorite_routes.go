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

	// ดึงรายการโปรดของผู้ใช้ปัจจุบัน
	favoriteGroup.Get("/user", favoriteHandler.GetUserFavorites)
	
	// เพิ่ม endpoint สำหรับการเพิ่มและลบสถานีโปรดแบบคลิกเดียว
	favoriteGroup.Post("/stations/:stationId", favoriteHandler.AddStationToFavorites) // เพิ่มสถานีเข้ารายการโปรด
	favoriteGroup.Delete("/:id/remove", favoriteHandler.RemoveStationFromFavorites) // ลบสถานีออกจากรายการโปรด
	
	// API เดิมสำหรับ CRUD ทั่วไป
	favoriteGroup.Get("/", middleware.OwnFavoriteMiddleware(favoriteService), favoriteHandler.GetAllFavorites)
	favoriteGroup.Post("/", favoriteHandler.CreateFavorite)
	favoriteGroup.Get("/:id", middleware.OwnFavoriteMiddleware(favoriteService), favoriteHandler.GetFavoriteByID)
	favoriteGroup.Put("/:id", middleware.OwnFavoriteMiddleware(favoriteService), favoriteHandler.UpdateFavorite)
	favoriteGroup.Delete("/:id", middleware.OwnFavoriteMiddleware(favoriteService), favoriteHandler.DeleteFavorite)
}
