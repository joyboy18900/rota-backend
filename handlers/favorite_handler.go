package handler

import (
	"rota-api/models"
	"rota-api/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type FavoriteHandler struct {
	favoriteService services.FavoriteService
}

func NewFavoriteHandler(favoriteService services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
	}
}

func (h *FavoriteHandler) GetFavoriteByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid favorite ID",
		})
	}

	favorite, err := h.favoriteService.GetFavoriteByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Favorite not found",
		})
	}

	return c.JSON(fiber.Map{
		"favorite": favorite,
	})
}

func (h *FavoriteHandler) GetAllFavorites(c *fiber.Ctx) error {
	favorites, err := h.favoriteService.GetAllFavorites(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"favorites": favorites,
	})
}

func (h *FavoriteHandler) CreateFavorite(c *fiber.Ctx) error {
	var favorite models.Favorite
	if err := c.BodyParser(&favorite); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.favoriteService.CreateFavorite(c.Context(), &favorite); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Favorite created successfully",
		"favorite": favorite,
	})
}

func (h *FavoriteHandler) UpdateFavorite(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid favorite ID",
		})
	}

	var favorite models.Favorite
	if err := c.BodyParser(&favorite); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	favorite.ID = uint(id)
	if err := h.favoriteService.UpdateFavorite(c.Context(), &favorite); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Favorite updated successfully",
		"favorite": favorite,
	})
}

func (h *FavoriteHandler) DeleteFavorite(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid favorite ID",
		})
	}

	if err := h.favoriteService.DeleteFavorite(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Favorite deleted successfully",
	})
}
