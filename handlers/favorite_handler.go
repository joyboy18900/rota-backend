package handler

import (
	"fmt"
	"rota-api/models"
	"rota-api/services"
	"rota-api/utils"
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
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid favorite ID")
	}

	favorite, err := h.favoriteService.GetFavoriteByID(c.Context(), uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Favorite not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Favorite retrieved successfully", fiber.Map{
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
	// Parse request body
	var favorite models.Favorite
	if err := c.BodyParser(&favorite); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Get user ID and role from context (set by AuthMiddleware)
	userID, _ := c.Locals("userID").(int)
	userRole := c.Locals("userRole")

	// Log values for debugging
	fmt.Printf("CreateFavorite - userID from context: %v (type: %T), favorite.UserID: %v (type: %T)\n", 
		userID, userID, favorite.UserID, favorite.UserID)

	// Security check: ensure user can only create favorites for themselves unless they're an admin
	if userRole != models.RoleAdmin && int(favorite.UserID) != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "forbidden: you can only create favorites for your own account",
		})
	}

	// Create favorite
	if err := h.favoriteService.CreateFavorite(c.Context(), &favorite); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
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
