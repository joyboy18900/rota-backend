package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"rota-api/services"
)

// StaffController interface defines methods for staff controller
type StaffController interface {
	CreateStaff(c *fiber.Ctx) error
	GetStaffByID(c *fiber.Ctx) error
	GetAllStaff(c *fiber.Ctx) error
	UpdateStaff(c *fiber.Ctx) error
	DeleteStaff(c *fiber.Ctx) error
}

// CreateStaffRequest represents the request body for creating a staff
type CreateStaffRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	StationID uint   `json:"station_id"`
}

// UpdateStaffRequest represents the request body for updating a staff
type UpdateStaffRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	StationID uint   `json:"station_id"`
}

// staffController implements StaffController
type staffController struct {
	staffService services.StaffService
}

// NewStaffController creates a new staff controller
func NewStaffController(staffService services.StaffService) StaffController {
	return &staffController{staffService}
}

// CreateStaff creates a new staff
func (c *staffController) CreateStaff(ctx *fiber.Ctx) error {
	var req CreateStaffRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	staff, err := c.staffService.CreateStaff(
		ctx.Context(),
		req.Username,
		req.Email,
		req.Password,
		req.StationID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    staff,
	})
}

// GetStaffByID retrieves a staff by ID
func (c *staffController) GetStaffByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid staff ID")
	}

	staff, err := c.staffService.GetStaffByID(ctx.Context(), uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    staff,
	})
}

// GetAllStaff retrieves all staff
func (c *staffController) GetAllStaff(ctx *fiber.Ctx) error {
	staff, err := c.staffService.GetAllStaff(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    staff,
	})
}

// UpdateStaff updates a staff
func (c *staffController) UpdateStaff(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid staff ID")
	}

	var req UpdateStaffRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	staff, err := c.staffService.UpdateStaff(
		ctx.Context(),
		uint(id),
		req.Username,
		req.Email,
		req.Password,
		req.StationID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Staff updated successfully",
		"data":    staff,
	})
}

// DeleteStaff deletes a staff
func (c *staffController) DeleteStaff(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid staff ID")
	}

	if err := c.staffService.DeleteStaff(ctx.Context(), uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Staff deleted successfully",
	})
}
