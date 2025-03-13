package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"rota-api/services"
)

// VehicleController interface defines methods for vehicle controller
type VehicleController interface {
	CreateVehicle(c *fiber.Ctx) error
	GetVehicleByID(c *fiber.Ctx) error
	GetAllVehicles(c *fiber.Ctx) error
	UpdateVehicle(c *fiber.Ctx) error
	DeleteVehicle(c *fiber.Ctx) error
}

// CreateVehicleRequest represents the request body for creating a vehicle
type CreateVehicleRequest struct {
	LicensePlate string `json:"license_plate"`
	Capacity     int    `json:"capacity"`
	DriverName   string `json:"driver_name"`
	RouteID      uint   `json:"route_id"`
}

// UpdateVehicleRequest represents the request body for updating a vehicle
type UpdateVehicleRequest struct {
	LicensePlate string `json:"license_plate"`
	Capacity     int    `json:"capacity"`
	DriverName   string `json:"driver_name"`
	RouteID      uint   `json:"route_id"`
}

// vehicleController implements VehicleController
type vehicleController struct {
	vehicleService services.VehicleService
}

// NewVehicleController creates a new vehicle controller
func NewVehicleController(vehicleService services.VehicleService) VehicleController {
	return &vehicleController{vehicleService}
}

// CreateVehicle creates a new vehicle
func (c *vehicleController) CreateVehicle(ctx *fiber.Ctx) error {
	var req CreateVehicleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	vehicle, err := c.vehicleService.CreateVehicle(
		ctx.Context(),
		req.LicensePlate,
		req.Capacity,
		req.DriverName,
		req.RouteID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    vehicle,
	})
}

// GetVehicleByID retrieves a vehicle by ID
func (c *vehicleController) GetVehicleByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid vehicle ID")
	}

	vehicle, err := c.vehicleService.GetVehicleByID(ctx.Context(), uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    vehicle,
	})
}

// GetAllVehicles retrieves all vehicles
func (c *vehicleController) GetAllVehicles(ctx *fiber.Ctx) error {
	vehicles, err := c.vehicleService.GetAllVehicles(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    vehicles,
	})
}

// UpdateVehicle updates a vehicle
func (c *vehicleController) UpdateVehicle(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid vehicle ID")
	}

	var req UpdateVehicleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	vehicle, err := c.vehicleService.UpdateVehicle(
		ctx.Context(),
		uint(id),
		req.LicensePlate,
		req.Capacity,
		req.DriverName,
		req.RouteID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Vehicle updated successfully",
		"data":    vehicle,
	})
}

// DeleteVehicle deletes a vehicle
func (c *vehicleController) DeleteVehicle(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid vehicle ID")
	}

	if err := c.vehicleService.DeleteVehicle(ctx.Context(), uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Vehicle deleted successfully",
	})
}
