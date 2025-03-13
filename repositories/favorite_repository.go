package repositories

import (
	"context"
	"errors"
	"fmt"
	"rota-api/models"

	"gorm.io/gorm"
)

// FavoriteRepository interface defines methods for favorite database operations
type FavoriteRepository interface {
	Create(ctx context.Context, favorite *models.Favorite) error
	FindByID(ctx context.Context, id uint) (*models.Favorite, error)
	FindByUserAndStation(ctx context.Context, userID, stationID uint) (*models.Favorite, error)
	FindByUser(ctx context.Context, userID uint) ([]models.Favorite, error)
	Delete(ctx context.Context, id uint) error
}

// favoriteRepository implements FavoriteRepository
type favoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository creates a new favorite repository
func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db}
}

// Create stores a new favorite in the database
func (r *favoriteRepository) Create(ctx context.Context, favorite *models.Favorite) error {
	if err := r.db.WithContext(ctx).Create(favorite).Error; err != nil {
		return fmt.Errorf("failed to create favorite: %w", err)
	}
	return nil
}

// FindByID retrieves a favorite by ID
func (r *favoriteRepository) FindByID(ctx context.Context, id uint) (*models.Favorite, error) {
	var favorite models.Favorite
	if err := r.db.WithContext(ctx).Preload("Station").First(&favorite, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("favorite not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find favorite: %w", err)
	}
	return &favorite, nil
}

// FindByUserAndStation retrieves a favorite by user ID and station ID
func (r *favoriteRepository) FindByUserAndStation(ctx context.Context, userID, stationID uint) (*models.Favorite, error) {
	var favorite models.Favorite
	if err := r.db.WithContext(ctx).Where("user_id = ? AND station_id = ?", userID, stationID).First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("favorite not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find favorite: %w", err)
	}
	return &favorite, nil
}

// FindByUser retrieves all favorites for a specific user
func (r *favoriteRepository) FindByUser(ctx context.Context, userID uint) ([]models.Favorite, error) {
	var favorites []models.Favorite
	if err := r.db.WithContext(ctx).Preload("Station").Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		return nil, fmt.Errorf("failed to find favorites: %w", err)
	}
	return favorites, nil
}

// Delete removes a favorite
func (r *favoriteRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Favorite{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete favorite: %w", err)
	}
	return nil
}
