package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"rota-api/services"
)

// StationController interface defines methods for station controller
type StationController interface {
	CreateStation(c *fiber.Ctx) error
	GetStationByID(c *fiber.Ctx) error
	GetAllStations(c *fiber.Ctx) error
	UpdateStation(c *fiber.Ctx) error
	DeleteStation(c *fiber.Ctx) error
}

// CreateStationRequest represents the request body for creating a station
type CreateStationRequest struct {
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"required"`
}

// UpdateStationRequest represents the request body for updating a station
type UpdateStationRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

// stationController implements StationController
type stationController struct {
	stationService services.StationService
}

// NewStationController creates a new station controller
func NewStationController(stationService services.StationService) StationController {
	return &stationController{stationService}
}

// CreateStation creates a new station
func (c *stationController) CreateStation(ctx *fiber.Ctx) error {
	var req CreateStationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Station name is required")
	}

	station, err := c.stationService.CreateStation(ctx.Context(), req.Name, req.Location)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Station created successfully",
		"data":    station,
	})
}

// GetStationByID retrieves a station by ID
func (c *stationController) GetStationByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid station ID")
	}

	station, err := c.stationService.GetStationByID(ctx.Context(), uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    station,
	})
}

// GetAllStations retrieves all stations
func (c *stationController) GetAllStations(ctx *fiber.Ctx) error {
	stations, err := c.stationService.GetAllStations(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    stations,
	})
}

// UpdateStation updates a station
func (c *stationController) UpdateStation(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid station ID")
	}

	var req UpdateStationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	station, err := c.stationService.UpdateStation(ctx.Context(), uint(id), req.Name, req.Location)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Station updated successfully",
		"data":    station,
	})
}

// DeleteStation deletes a station
func (c *stationController) DeleteStation(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid station ID")
	}

	if err := c.stationService.DeleteStation(ctx.Context(), uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Station deleted successfully",
	})
}
