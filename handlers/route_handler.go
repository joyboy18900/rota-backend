package handler

import (
	"rota-api/models"
	"rota-api/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RouteHandler struct {
	routeService services.RouteService
}

func NewRouteHandler(routeService services.RouteService) *RouteHandler {
	return &RouteHandler{
		routeService: routeService,
	}
}

func (h *RouteHandler) GetRouteByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid route ID",
		})
	}

	route, err := h.routeService.GetRouteByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Route not found",
		})
	}

	return c.JSON(fiber.Map{
		"route": route,
	})
}

func (h *RouteHandler) GetAllRoutes(c *fiber.Ctx) error {
	routes, err := h.routeService.GetAllRoutes(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"routes": routes,
	})
}

func (h *RouteHandler) CreateRoute(c *fiber.Ctx) error {
	var route models.Route
	if err := c.BodyParser(&route); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.routeService.CreateRoute(c.Context(), &route); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Route created successfully",
		"route":   route,
	})
}

func (h *RouteHandler) UpdateRoute(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid route ID",
		})
	}

	var route models.Route
	if err := c.BodyParser(&route); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	route.ID = uint(id)
	if err := h.routeService.UpdateRoute(c.Context(), &route); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Route updated successfully",
		"route":   route,
	})
}

func (h *RouteHandler) DeleteRoute(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid route ID",
		})
	}

	if err := h.routeService.DeleteRoute(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Route deleted successfully",
	})
}
