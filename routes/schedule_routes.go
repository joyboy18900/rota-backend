package routes

import (
	"rota-api/handlers"
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

func SetupScheduleRoutes(
	app *fiber.App,
	scheduleHandler *handler.ScheduleHandler,
	authService services.AuthService,
) {
	scheduleGroup := app.Group("/api/v1/schedules")
	scheduleGroup.Use(middleware.AuthMiddleware(authService))

	// Advanced search endpoint - ต้องอยู่ก่อนเส้นทาง /:id เพื่อป้องกันการจับคู่ผิดพลาด
	scheduleGroup.Get("/search", scheduleHandler.SearchSchedules)
	
	// CRUD operations
	scheduleGroup.Post("/", scheduleHandler.CreateSchedule)
	scheduleGroup.Get("/", scheduleHandler.GetAllSchedules)
	scheduleGroup.Get("/:id", scheduleHandler.GetScheduleByID)
	scheduleGroup.Put("/:id", scheduleHandler.UpdateSchedule)
	scheduleGroup.Delete("/:id", scheduleHandler.DeleteSchedule)
}
