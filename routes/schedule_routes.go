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
	// Public group for read-only operations that don't need authentication
	publicGroup := app.Group("/api/v1/schedules")
	
	// Advanced search endpoint - ต้องอยู่ก่อนเส้นทาง /:id เพื่อป้องกันการจับคู่ผิดพลาด
	publicGroup.Get("/search", scheduleHandler.SearchSchedules)
	
	// Read operations - accessible to all users without authentication
	publicGroup.Get("/", scheduleHandler.GetAllSchedules)
	publicGroup.Get("/:id", scheduleHandler.GetScheduleByID)
	
	// Admin-only operations for schedule management
	adminScheduleGroup := app.Group("/api/v1/schedules")
	adminScheduleGroup.Use(middleware.AuthMiddleware(authService), middleware.AdminMiddleware())
	adminScheduleGroup.Post("/", scheduleHandler.CreateSchedule)
	adminScheduleGroup.Put("/:id", scheduleHandler.UpdateSchedule)
	adminScheduleGroup.Delete("/:id", scheduleHandler.DeleteSchedule)
}
