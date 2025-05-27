package routes

import (
	"rota-api/handlers"
	"rota-api/middleware"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
)

func SetupScheduleLogRoutes(
	app *fiber.App,
	scheduleLogHandler *handler.ScheduleLogHandler,
	authService services.AuthService,
) {
	scheduleLogGroup := app.Group("/api/v1/schedule-logs")
	scheduleLogGroup.Use(middleware.AuthMiddleware(authService))

	// Read operations - accessible to all authenticated users
	scheduleLogGroup.Get("/", scheduleLogHandler.GetAllScheduleLogs)
	scheduleLogGroup.Get("/:id", scheduleLogHandler.GetScheduleLogByID)
	
	// Admin-only operations for schedule log management
	adminScheduleLogGroup := app.Group("/api/v1/schedule-logs")
	adminScheduleLogGroup.Use(middleware.AuthMiddleware(authService), middleware.AdminMiddleware())
	adminScheduleLogGroup.Post("/", scheduleLogHandler.CreateScheduleLog)
	adminScheduleLogGroup.Put("/:id", scheduleLogHandler.UpdateScheduleLog)
	adminScheduleLogGroup.Delete("/:id", scheduleLogHandler.DeleteScheduleLog)
}
