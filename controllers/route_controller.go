package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"rota-api/services"
)

// RouteController interface defines methods for route controller
type RouteController interface {
	CreateRoute(c *fiber.Ctx) error
	GetRouteByID(c *fiber.Ctx) error
	GetAllRoutes(c *fiber.Ctx) error
	UpdateRoute(c *fiber.Ctx) error
	DeleteRoute(c *fiber.Ctx) error
}

// CreateRouteRequest represents the request body for creating a route
type CreateRouteRequest struct {
	StartStationID uint    `json:"start_station_id"`
	EndStationID   uint    `json:"end_station_id"`
	Distance       float64 `json:"distance"`
	Duration       string  `json:"duration"`
}

// UpdateRouteRequest represents the request body for updating a route
type UpdateRouteRequest struct {
	StartStationID uint    `json:"start_station_id"`
	EndStationID   uint    `json:"end_station_id"`
	Distance       float64 `json:"distance"`
	Duration       string  `json:"duration"`
}

// routeController implements RouteController
type routeController struct {
	routeService services.RouteService
}

// NewRouteController creates a new route controller
func NewRouteController(routeService services.RouteService) RouteController {
	return &routeController{routeService}
}

// CreateRoute creates a new route
func (c *routeController) CreateRoute(ctx *fiber.Ctx) error {
	var req CreateRouteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	route, err := c.routeService.CreateRoute(
		ctx.Context(),
		req.StartStationID,
		req.EndStationID,
		req.Distance,
		req.Duration,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    route,
	})
}

// GetRouteByID retrieves a route by ID
func (c *routeController) GetRouteByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid route ID")
	}

	route, err := c.routeService.GetRouteByID(ctx.Context(), uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    route,
	})
}

// GetAllRoutes retrieves all routes
func (c *routeController) GetAllRoutes(ctx *fiber.Ctx) error {
	routes, err := c.routeService.GetAllRoutes(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    routes,
	})
}

// UpdateRoute updates a route
func (c *routeController) UpdateRoute(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid route ID")
	}

	var req UpdateRouteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	route, err := c.routeService.UpdateRoute(
		ctx.Context(),
		uint(id),
		req.StartStationID,
		req.EndStationID,
		req.Distance,
		req.Duration,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Route updated successfully",
		"data":    route,
	})
}

// DeleteRoute deletes a route
func (c *routeController) DeleteRoute(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid route ID")
	}

	if err := c.routeService.DeleteRoute(ctx.Context(), uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Route deleted successfully",
	})
}
