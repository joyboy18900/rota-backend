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

	scheduleGroup.Post("/", scheduleHandler.CreateSchedule)
	scheduleGroup.Get("/", scheduleHandler.GetAllSchedules)
	scheduleGroup.Get("/:id", scheduleHandler.GetScheduleByID)
	scheduleGroup.Put("/:id", scheduleHandler.UpdateSchedule)
	scheduleGroup.Delete("/:id", scheduleHandler.DeleteSchedule)
}
