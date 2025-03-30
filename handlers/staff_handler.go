package handler

import (
	"rota-api/models"
	"rota-api/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type StaffHandler struct {
	staffService services.StaffService
}

func NewStaffHandler(staffService services.StaffService) *StaffHandler {
	return &StaffHandler{
		staffService: staffService,
	}
}

func (h *StaffHandler) GetStaffByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid staff ID",
		})
	}

	staff, err := h.staffService.GetStaffByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Staff not found",
		})
	}

	return c.JSON(fiber.Map{
		"staff": staff,
	})
}

func (h *StaffHandler) GetAllStaff(c *fiber.Ctx) error {
	staff, err := h.staffService.GetAllStaff(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"staff": staff,
	})
}

func (h *StaffHandler) CreateStaff(c *fiber.Ctx) error {
	var staff models.Staff
	if err := c.BodyParser(&staff); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.staffService.CreateStaff(c.Context(), &staff); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Staff created successfully",
		"staff":   staff,
	})
}

func (h *StaffHandler) UpdateStaff(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid staff ID",
		})
	}

	var staff models.Staff
	if err := c.BodyParser(&staff); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	staff.ID = uint(id)
	if err := h.staffService.UpdateStaff(c.Context(), &staff); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Staff updated successfully",
		"staff":   staff,
	})
}

func (h *StaffHandler) DeleteStaff(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid staff ID",
		})
	}

	if err := h.staffService.DeleteStaff(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Staff deleted successfully",
	})
}
