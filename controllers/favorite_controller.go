package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"rota-api/services"
)

// FavoriteController interface defines methods for favorite controller
type FavoriteController interface {
	AddFavorite(c *fiber.Ctx) error
	GetUserFavorites(c *fiber.Ctx) error
	RemoveFavorite(c *fiber.Ctx) error
}

// AddFavoriteRequest represents the request body for adding a favorite
type AddFavoriteRequest struct {
	StationID uint `json:"station_id"`
}

// favoriteController implements FavoriteController
type favoriteController struct {
	favoriteService services.FavoriteService
}

// NewFavoriteController creates a new favorite controller
func NewFavoriteController(favoriteService services.FavoriteService) FavoriteController {
	return &favoriteController{favoriteService}
}

// AddFavorite adds a station to user's favorites
func (c *favoriteController) AddFavorite(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userID").(uint)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	var req AddFavoriteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	favorite, err := c.favoriteService.AddFavorite(ctx.Context(), userID, req.StationID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    favorite,
	})
}

// GetUserFavorites retrieves all favorites for the authenticated user
func (c *favoriteController) GetUserFavorites(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userID").(uint)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	favorites, err := c.favoriteService.GetUserFavorites(ctx.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    favorites,
	})
}

// RemoveFavorite removes a station from user's favorites
func (c *favoriteController) RemoveFavorite(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userID").(uint)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid favorite ID")
	}

	if err := c.favoriteService.RemoveFavorite(ctx.Context(), uint(id), userID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Favorite removed successfully",
	})
}
