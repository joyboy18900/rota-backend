package repositories

import (
	"context"
	"rota-api/models"

	"gorm.io/gorm"
)

// FavoriteRepository defines the interface for favorite-related database operations
type FavoriteRepository interface {
	FindByID(ctx context.Context, id uint) (*models.Favorite, error)
	FindAll(ctx context.Context) ([]models.Favorite, error)
	FindByUserAndStation(ctx context.Context, userID, stationID uint) (*models.Favorite, error)
	FindByUser(ctx context.Context, userID uint) ([]models.Favorite, error)
	Create(ctx context.Context, favorite *models.Favorite) error
	Update(ctx context.Context, favorite *models.Favorite) error
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

func (r *favoriteRepository) FindByID(ctx context.Context, id uint) (*models.Favorite, error) {
	var favorite models.Favorite
	if err := r.db.WithContext(ctx).Preload("Station").First(&favorite, id).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *favoriteRepository) FindAll(ctx context.Context) ([]models.Favorite, error) {
	var favorites []models.Favorite
	if err := r.db.WithContext(ctx).Preload("User").Preload("Station").Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func (r *favoriteRepository) FindByUserAndStation(ctx context.Context, userID, stationID uint) (*models.Favorite, error) {
	var favorite models.Favorite
	if err := r.db.WithContext(ctx).Preload("Station").Where("user_id = ? AND station_id = ?", userID, stationID).First(&favorite).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *favoriteRepository) FindByUser(ctx context.Context, userID uint) ([]models.Favorite, error) {
	var favorites []models.Favorite
	if err := r.db.WithContext(ctx).Preload("Station").Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func (r *favoriteRepository) Create(ctx context.Context, favorite *models.Favorite) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

func (r *favoriteRepository) Update(ctx context.Context, favorite *models.Favorite) error {
	return r.db.WithContext(ctx).Save(favorite).Error
}

func (r *favoriteRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Favorite{}, id).Error
}
