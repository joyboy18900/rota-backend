package handler

import (
	"rota-api/dto"
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

func (h *FavoriteHandler) AddStationToFavorites(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		userIDInt, ok := c.Locals("userID").(int)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
		}
		userID = uint(userIDInt)
	}

	stationID, err := strconv.ParseUint(c.Params("stationId"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid station ID")
	}

	favorite, err := h.favoriteService.AddFavorite(c.Context(), userID, uint(stationID))
	if err != nil {
		if err.Error() == "station is already a favorite" {
			existingFavorite, findErr := h.favoriteService.GetFavoriteByUserAndStation(c.Context(), userID, uint(stationID))
			if findErr != nil {
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch favorite information")
			}
			return utils.SuccessResponse(c, fiber.StatusOK, "Station is already in your favorites", existingFavorite)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add favorite: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Station added to favorites", favorite)
}

func (h *FavoriteHandler) RemoveStationByStationId(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		userIDInt, ok := c.Locals("userID").(int)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
		}
		userID = uint(userIDInt)
	}

	stationID, err := strconv.ParseUint(c.Params("stationId"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid station ID")
	}

	favorite, err := h.favoriteService.GetFavoriteByUserAndStation(c.Context(), userID, uint(stationID))
	if err != nil {
		if err.Error() == "favorite not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Station not in favorites")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to find favorite: "+err.Error())
	}

	if err := h.favoriteService.RemoveFavorite(c.Context(), favorite.ID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to remove favorite: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Station removed from favorites", nil)
}

func (h *FavoriteHandler) RemoveStationFromFavorites(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		userIDInt, ok := c.Locals("userID").(int)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
		}
		userID = uint(userIDInt)
	}

	favoriteID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid favorite ID")
	}

	if err := h.favoriteService.RemoveFavorite(c.Context(), uint(favoriteID), userID); err != nil {
		if err.Error() == "favorite not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Favorite not found")
		}
		if err.Error() == "unauthorized to remove this favorite" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not authorized to remove this favorite")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to remove favorite: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Station removed from favorites", nil)
}

func (h *FavoriteHandler) GetUserFavorites(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		userIDInt, ok := c.Locals("userID").(int)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
		}
		userID = uint(userIDInt)
	}

	favorites, err := h.favoriteService.GetUserFavorites(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get favorites: "+err.Error())
	}

	response := make([]dto.FavoriteStationResponse, 0, len(favorites))
	for _, fav := range favorites {
		response = append(response, dto.FavoriteStationResponse{
			StationID: fav.StationID,
			CreatedAt: fav.CreatedAt,
			Name:      fav.Station.Name,
			Location:  fav.Station.Location,
		})
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Favorites retrieved successfully", response)
}
