package repositories

import (
	"context"
	"errors"
	"fmt"
	"rota-api/models"

	"gorm.io/gorm"
)

// StationRepository interface defines methods for station database operations
type StationRepository interface {
	Create(ctx context.Context, station *models.Station) error
	FindByID(ctx context.Context, id uint) (*models.Station, error)
	FindAll(ctx context.Context) ([]models.Station, error)
	Update(ctx context.Context, station *models.Station) error
	Delete(ctx context.Context, id uint) error
}

// stationRepository implements StationRepository
type stationRepository struct {
	db *gorm.DB
}

// NewStationRepository creates a new station repository
func NewStationRepository(db *gorm.DB) StationRepository {
	return &stationRepository{db}
}

// Create stores a new station in the database
func (r *stationRepository) Create(ctx context.Context, station *models.Station) error {
	if err := r.db.WithContext(ctx).Create(station).Error; err != nil {
		return fmt.Errorf("failed to create station: %w", err)
	}
	return nil
}

// FindByID retrieves a station by ID
func (r *stationRepository) FindByID(ctx context.Context, id uint) (*models.Station, error) {
	var station models.Station
	if err := r.db.WithContext(ctx).First(&station, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("station not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find station: %w", err)
	}
	return &station, nil
}

// FindAll retrieves all stations
func (r *stationRepository) FindAll(ctx context.Context) ([]models.Station, error) {
	var stations []models.Station
	if err := r.db.WithContext(ctx).Find(&stations).Error; err != nil {
		return nil, fmt.Errorf("failed to find stations: %w", err)
	}
	return stations, nil
}

// Update updates a station
func (r *stationRepository) Update(ctx context.Context, station *models.Station) error {
	if err := r.db.WithContext(ctx).Save(station).Error; err != nil {
		return fmt.Errorf("failed to update station: %w", err)
	}
	return nil
}

// Delete removes a station
func (r *stationRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Station{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete station: %w", err)
	}
	return nil
}
