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

	scheduleLogGroup.Post("/", scheduleLogHandler.CreateScheduleLog)
	scheduleLogGroup.Get("/", scheduleLogHandler.GetAllScheduleLogs)
	scheduleLogGroup.Get("/:id", scheduleLogHandler.GetScheduleLogByID)
	scheduleLogGroup.Put("/:id", scheduleLogHandler.UpdateScheduleLog)
	scheduleLogGroup.Delete("/:id", scheduleLogHandler.DeleteScheduleLog)
}
