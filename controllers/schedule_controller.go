package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"rota-api/services"
)

// ScheduleController interface defines methods for schedule controller
type ScheduleController interface {
	CreateSchedule(c *fiber.Ctx) error
	GetScheduleByID(c *fiber.Ctx) error
	GetAllSchedules(c *fiber.Ctx) error
	UpdateSchedule(c *fiber.Ctx) error
	DeleteSchedule(c *fiber.Ctx) error
}

// CreateScheduleRequest represents the request body for creating a schedule
type CreateScheduleRequest struct {
	RouteID       uint   `json:"route_id"`
	StationID     uint   `json:"station_id"`
	Round         int    `json:"round"`
	DepartureTime string `json:"departure_time"`
	ArrivalTime   string `json:"arrival_time"`
}

// UpdateScheduleRequest represents the request body for updating a schedule
type UpdateScheduleRequest struct {
	RouteID       uint   `json:"route_id"`
	StationID     uint   `json:"station_id"`
	Round         int    `json:"round"`
	DepartureTime string `json:"departure_time"`
	ArrivalTime   string `json:"arrival_time"`
}

// scheduleController implements ScheduleController
type scheduleController struct {
	scheduleService services.ScheduleService
}

// NewScheduleController creates a new schedule controller
func NewScheduleController(scheduleService services.ScheduleService) ScheduleController {
	return &scheduleController{scheduleService}
}

// CreateSchedule creates a new schedule
func (c *scheduleController) CreateSchedule(ctx *fiber.Ctx) error {
	var req CreateScheduleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	schedule, err := c.scheduleService.CreateSchedule(
		ctx.Context(),
		req.RouteID,
		req.StationID,
		req.Round,
		req.DepartureTime,
		req.ArrivalTime,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    schedule,
	})
}

// GetScheduleByID retrieves a schedule by ID
func (c *scheduleController) GetScheduleByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid schedule ID")
	}

	schedule, err := c.scheduleService.GetScheduleByID(ctx.Context(), uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    schedule,
	})
}

// GetAllSchedules retrieves all schedules
func (c *scheduleController) GetAllSchedules(ctx *fiber.Ctx) error {
	schedules, err := c.scheduleService.GetAllSchedules(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    schedules,
	})
}

// UpdateSchedule updates a schedule
func (c *scheduleController) UpdateSchedule(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid schedule ID")
	}

	var req UpdateScheduleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	schedule, err := c.scheduleService.UpdateSchedule(
		ctx.Context(),
		uint(id),
		req.RouteID,
		req.StationID,
		req.Round,
		req.DepartureTime,
		req.ArrivalTime,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Schedule updated successfully",
		"data":    schedule,
	})
}

// DeleteSchedule deletes a schedule
func (c *scheduleController) DeleteSchedule(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid schedule ID")
	}

	if err := c.scheduleService.DeleteSchedule(ctx.Context(), uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Schedule deleted successfully",
	})
}
