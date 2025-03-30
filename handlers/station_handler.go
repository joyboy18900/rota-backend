package handler

import (
	"rota-api/models"
	"rota-api/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type StationHandler struct {
	stationService services.StationService
}

func NewStationHandler(stationService services.StationService) *StationHandler {
	return &StationHandler{
		stationService: stationService,
	}
}

func (h *StationHandler) GetStationByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid station ID",
		})
	}

	station, err := h.stationService.GetStationByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Station not found",
		})
	}

	return c.JSON(fiber.Map{
		"station": station,
	})
}

func (h *StationHandler) GetAllStations(c *fiber.Ctx) error {
	stations, err := h.stationService.GetAllStations(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"stations": stations,
	})
}

func (h *StationHandler) CreateStation(c *fiber.Ctx) error {
	var station models.Station
	if err := c.BodyParser(&station); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.stationService.CreateStation(c.Context(), &station); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Station created successfully",
		"station": station,
	})
}

func (h *StationHandler) UpdateStation(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid station ID",
		})
	}

	var station models.Station
	if err := c.BodyParser(&station); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	station.ID = uint(id)
	if err := h.stationService.UpdateStation(c.Context(), &station); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Station updated successfully",
		"station": station,
	})
}

func (h *StationHandler) DeleteStation(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid station ID",
		})
	}

	if err := h.stationService.DeleteStation(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Station deleted successfully",
	})
}
