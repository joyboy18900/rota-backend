package services

import (
	"context"
	"fmt"
	"time"

	"rota-api/models"
	"rota-api/repositories"
)

// FavoriteService interface defines methods for favorite service
type FavoriteService interface {
	AddFavorite(ctx context.Context, userID, stationID uint) (*models.Favorite, error)
	GetUserFavorites(ctx context.Context, userID uint) ([]models.Favorite, error)
	RemoveFavorite(ctx context.Context, id, userID uint) error
	GetFavoriteByID(ctx context.Context, id uint) (*models.Favorite, error)
	GetFavoriteByUserAndStation(ctx context.Context, userID, stationID uint) (*models.Favorite, error)
	GetAllFavorites(ctx context.Context) ([]models.Favorite, error)
	CreateFavorite(ctx context.Context, favorite *models.Favorite) error
	UpdateFavorite(ctx context.Context, favorite *models.Favorite) error
	DeleteFavorite(ctx context.Context, id uint) error
}

// favoriteService implements FavoriteService
type favoriteService struct {
	favoriteRepo repositories.FavoriteRepository
}

// NewFavoriteService creates a new favorite service
func NewFavoriteService(favoriteRepo repositories.FavoriteRepository) FavoriteService {
	return &favoriteService{favoriteRepo}
}

// AddFavorite adds a station to user's favorites
func (s *favoriteService) AddFavorite(ctx context.Context, userID, stationID uint) (*models.Favorite, error) {
	// Check if already a favorite
	_, err := s.favoriteRepo.FindByUserAndStation(ctx, userID, stationID)
	if err == nil {
		return nil, fmt.Errorf("station is already a favorite")
	}

	favorite := &models.Favorite{
		UserID:    userID,
		StationID: stationID,
		CreatedAt: time.Now(),
	}

	if err := s.favoriteRepo.Create(ctx, favorite); err != nil {
		return nil, fmt.Errorf("failed to add favorite: %w", err)
	}

	// โหลดข้อมูล favorite ที่สร้างเสร็จแล้วจากฐานข้อมูล พร้อมความสัมพันธ์
	completeFavorite, err := s.favoriteRepo.FindByID(ctx, favorite.ID)
	if err != nil {
		return favorite, nil // ส่งคืนข้อมูลธรรมดาถ้าไม่สามารถโหลดข้อมูลเพิ่มเติม
	}

	return completeFavorite, nil
}

// GetUserFavorites retrieves all favorites for a user
func (s *favoriteService) GetUserFavorites(ctx context.Context, userID uint) ([]models.Favorite, error) {
	return s.favoriteRepo.FindByUser(ctx, userID)
}

// RemoveFavorite removes a station from user's favorites
func (s *favoriteService) RemoveFavorite(ctx context.Context, id, userID uint) error {
	// Verify ownership
	favorite, err := s.favoriteRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("favorite not found: %w", err)
	}

	if favorite.UserID != userID {
		return fmt.Errorf("unauthorized to remove this favorite")
	}

	return s.favoriteRepo.Delete(ctx, id)
}

func (s *favoriteService) GetFavoriteByID(ctx context.Context, id uint) (*models.Favorite, error) {
	return s.favoriteRepo.FindByID(ctx, id)
}

func (s *favoriteService) GetAllFavorites(ctx context.Context) ([]models.Favorite, error) {
	return s.favoriteRepo.FindAll(ctx)
}

func (s *favoriteService) CreateFavorite(ctx context.Context, favorite *models.Favorite) error {
	return s.favoriteRepo.Create(ctx, favorite)
}

func (s *favoriteService) UpdateFavorite(ctx context.Context, favorite *models.Favorite) error {
	return s.favoriteRepo.Update(ctx, favorite)
}

func (s *favoriteService) DeleteFavorite(ctx context.Context, id uint) error {
	return s.favoriteRepo.Delete(ctx, id)
}

// GetFavoriteByUserAndStation ดึงรายการโปรดตามผู้ใช้และสถานี
func (s *favoriteService) GetFavoriteByUserAndStation(ctx context.Context, userID, stationID uint) (*models.Favorite, error) {
	return s.favoriteRepo.FindByUserAndStation(ctx, userID, stationID)
}
