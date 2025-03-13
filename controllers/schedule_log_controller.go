package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"rota-api/services"
)

// ScheduleLogController interface defines methods for schedule log controller
type ScheduleLogController interface {
	CreateScheduleLog(c *fiber.Ctx) error
	GetScheduleLogByID(c *fiber.Ctx) error
	GetAllScheduleLogs(c *fiber.Ctx) error
}

// CreateScheduleLogRequest represents the request body for creating a schedule log
type CreateScheduleLogRequest struct {
	ScheduleID        uint   `json:"schedule_id"`
	StaffID           uint   `json:"staff_id"`
	ChangeDescription string `json:"change_description"`
}

// scheduleLogController implements ScheduleLogController
type scheduleLogController struct {
	scheduleLogService services.ScheduleLogService
}

// NewScheduleLogController creates a new schedule log controller
func NewScheduleLogController(scheduleLogService services.ScheduleLogService) ScheduleLogController {
	return &scheduleLogController{scheduleLogService}
}

// CreateScheduleLog creates a new schedule log
func (c *scheduleLogController) CreateScheduleLog(ctx *fiber.Ctx) error {
	var req CreateScheduleLogRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	log, err := c.scheduleLogService.CreateScheduleLog(
		ctx.Context(),
		req.ScheduleID,
		req.StaffID,
		req.ChangeDescription,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    log,
	})
}

// GetScheduleLogByID retrieves a schedule log by ID
func (c *scheduleLogController) GetScheduleLogByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid schedule log ID")
	}

	log, err := c.scheduleLogService.GetScheduleLogByID(ctx.Context(), uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    log,
	})
}

// GetAllScheduleLogs retrieves all schedule logs
func (c *scheduleLogController) GetAllScheduleLogs(ctx *fiber.Ctx) error {
	logs, err := c.scheduleLogService.GetAllScheduleLogs(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    logs,
	})
}
