package handler

import (
	"rota-api/models"
	"rota-api/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ScheduleLogHandler struct {
	scheduleLogService services.ScheduleLogService
}

func NewScheduleLogHandler(scheduleLogService services.ScheduleLogService) *ScheduleLogHandler {
	return &ScheduleLogHandler{
		scheduleLogService: scheduleLogService,
	}
}

func (h *ScheduleLogHandler) GetScheduleLogByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schedule log ID",
		})
	}

	scheduleLog, err := h.scheduleLogService.GetScheduleLogByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Schedule log not found",
		})
	}

	return c.JSON(fiber.Map{
		"schedule_log": scheduleLog,
	})
}

func (h *ScheduleLogHandler) GetAllScheduleLogs(c *fiber.Ctx) error {
	scheduleLogs, err := h.scheduleLogService.GetAllScheduleLogs(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"schedule_logs": scheduleLogs,
	})
}

func (h *ScheduleLogHandler) CreateScheduleLog(c *fiber.Ctx) error {
	var scheduleLog models.ScheduleLog
	if err := c.BodyParser(&scheduleLog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if c.Get("Content-Type") == "application/json" {
		body := make(map[string]interface{})
		if err := c.BodyParser(&body); err == nil {
			if _, ok := body["actual_departure"]; ok && body["actual_departure"] != "" {
				if t, err := time.Parse(time.RFC3339, body["actual_departure"].(string)); err == nil {
					scheduleLog.ActualDeparture = &t
				}
			}

			if _, ok := body["actual_arrival"]; ok && body["actual_arrival"] != "" {
				if t, err := time.Parse(time.RFC3339, body["actual_arrival"].(string)); err == nil {
					scheduleLog.ActualArrival = &t
				}
			}
		}
	}

	now := time.Now()
	scheduleLog.UpdatedAt = &now

	if err := h.scheduleLogService.CreateScheduleLog(c.Context(), &scheduleLog); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Schedule log created successfully",
		"schedule_log": scheduleLog,
	})
}

func (h *ScheduleLogHandler) UpdateScheduleLog(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schedule log ID",
		})
	}

	var scheduleLog models.ScheduleLog
	if err := c.BodyParser(&scheduleLog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if c.Get("Content-Type") == "application/json" {
		body := make(map[string]interface{})
		if err := c.BodyParser(&body); err == nil {
			if _, ok := body["actual_departure"]; ok && body["actual_departure"] != "" {
				if t, err := time.Parse(time.RFC3339, body["actual_departure"].(string)); err == nil {
					scheduleLog.ActualDeparture = &t
				}
			}

			if _, ok := body["actual_arrival"]; ok && body["actual_arrival"] != "" {
				if t, err := time.Parse(time.RFC3339, body["actual_arrival"].(string)); err == nil {
					scheduleLog.ActualArrival = &t
				}
			}
		}
	}

	now := time.Now()
	scheduleLog.UpdatedAt = &now
	scheduleLog.ID = uint(id)

	if err := h.scheduleLogService.UpdateScheduleLog(c.Context(), &scheduleLog); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":      "Schedule log updated successfully",
		"schedule_log": scheduleLog,
	})
}

func (h *ScheduleLogHandler) DeleteScheduleLog(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schedule log ID",
		})
	}

	if err := h.scheduleLogService.DeleteScheduleLog(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Schedule log deleted successfully",
	})
}
