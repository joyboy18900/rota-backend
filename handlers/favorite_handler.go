package handler

import (
	"fmt"
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

// AddStationToFavorites เพิ่มสถานีเข้าไปในรายการโปรดของผู้ใช้ด้วยการคลิกเดียว
func (h *FavoriteHandler) AddStationToFavorites(c *fiber.Ctx) error {
	// ดึง user ID จาก context (ตั้งค่าโดย AuthMiddleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		// ถ้า userID จาก context ไม่ใช่ uint ให้ลองแปลงจาก int
		userIDInt, ok := c.Locals("userID").(int)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
		}
		userID = uint(userIDInt)
	}

	// ดึง station ID จาก URL parameter และตรวจสอบความถูกต้อง
	stationID, err := strconv.ParseUint(c.Params("stationId"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid station ID")
	}

	// เรียกใช้ service method เพื่อเพิ่มสถานีโปรด
	favorite, err := h.favoriteService.AddFavorite(c.Context(), userID, uint(stationID))
	if err != nil {
		// ตรวจสอบข้อความ error เพื่อให้การตอบกลับที่เหมาะสม
		if err.Error() == "station is already a favorite" {
			// ส่ง success response แทน error เมื่อสถานีอยู่ในรายการโปรดอยู่แล้ว
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

// RemoveStationFromFavorites ลบสถานีออกจากรายการโปรดของผู้ใช้ด้วยการคลิกเดียว
func (h *FavoriteHandler) RemoveStationFromFavorites(c *fiber.Ctx) error {
	// ดึง user ID จาก context (ตั้งค่าโดย AuthMiddleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		// ถ้า userID จาก context ไม่ใช่ uint ให้ลองแปลงจาก int
		userIDInt, ok := c.Locals("userID").(int)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
		}
		userID = uint(userIDInt)
	}

	// ดึง favorite ID จาก URL parameter และตรวจสอบความถูกต้อง
	favoriteID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid favorite ID")
	}

	// เรียกใช้ service method เพื่อลบสถานีโปรด
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

// GetUserFavorites ดึงรายการสถานีโปรดทั้งหมดของผู้ใช้ปัจจุบัน
func (h *FavoriteHandler) GetUserFavorites(c *fiber.Ctx) error {
	// ดึง user ID จาก context (ตั้งค่าโดย AuthMiddleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		// ถ้า userID จาก context ไม่ใช่ uint ให้ลองแปลงจาก int
		userIDInt, ok := c.Locals("userID").(int)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated")
		}
		userID = uint(userIDInt)
	}

	// เรียกใช้ service method เพื่อดึงรายการสถานีโปรด
	favorites, err := h.favoriteService.GetUserFavorites(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get favorites: "+err.Error())
	}

	// แปลงข้อมูลให้อยู่ในรูปแบบที่กระชับและตรงตามความต้องการ
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
