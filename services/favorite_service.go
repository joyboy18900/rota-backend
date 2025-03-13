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

	return favorite, nil
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
