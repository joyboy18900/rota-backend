package repositories

import (
	"context"
	"errors"
	"fmt"
	"rota-api/models"

	"gorm.io/gorm"
)

// VehicleRepository interface defines methods for vehicle database operations
type VehicleRepository interface {
	Create(ctx context.Context, vehicle *models.Vehicle) error
	FindByID(ctx context.Context, id uint) (*models.Vehicle, error)
	FindByLicensePlate(ctx context.Context, licensePlate string) (*models.Vehicle, error)
	FindAll(ctx context.Context) ([]*models.Vehicle, error)
	FindByRoute(ctx context.Context, routeID uint) ([]models.Vehicle, error)
	Update(ctx context.Context, vehicle *models.Vehicle) error
	Delete(ctx context.Context, id uint) error
}

// vehicleRepository implements VehicleRepository
type vehicleRepository struct {
	db *gorm.DB
}

// NewVehicleRepository creates a new vehicle repository
func NewVehicleRepository(db *gorm.DB) VehicleRepository {
	return &vehicleRepository{db}
}

// Create stores a new vehicle in the database
func (r *vehicleRepository) Create(ctx context.Context, vehicle *models.Vehicle) error {
	if err := r.db.WithContext(ctx).Create(vehicle).Error; err != nil {
		return fmt.Errorf("failed to create vehicle: %w", err)
	}
	return nil
}

// FindByID retrieves a vehicle by ID
func (r *vehicleRepository) FindByID(ctx context.Context, id uint) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	if err := r.db.WithContext(ctx).Preload("Route").First(&vehicle, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("vehicle not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find vehicle: %w", err)
	}
	return &vehicle, nil
}

// FindByLicensePlate retrieves a vehicle by license plate
func (r *vehicleRepository) FindByLicensePlate(ctx context.Context, licensePlate string) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	if err := r.db.WithContext(ctx).Where("license_plate = ?", licensePlate).First(&vehicle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("vehicle not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find vehicle: %w", err)
	}
	return &vehicle, nil
}

// FindAll retrieves all vehicles
func (r *vehicleRepository) FindAll(ctx context.Context) ([]*models.Vehicle, error) {
	var vehicles []*models.Vehicle
	if err := r.db.WithContext(ctx).Preload("Route").Find(&vehicles).Error; err != nil {
		return nil, fmt.Errorf("failed to find vehicles: %w", err)
	}
	return vehicles, nil
}

// FindByRoute retrieves all vehicles for a specific route
func (r *vehicleRepository) FindByRoute(ctx context.Context, routeID uint) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	if err := r.db.WithContext(ctx).Preload("Route").Where("route_id = ?", routeID).Find(&vehicles).Error; err != nil {
		return nil, fmt.Errorf("failed to find vehicles: %w", err)
	}
	return vehicles, nil
}

// Update updates a vehicle
func (r *vehicleRepository) Update(ctx context.Context, vehicle *models.Vehicle) error {
	if err := r.db.WithContext(ctx).Save(vehicle).Error; err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}
	return nil
}

// Delete removes a vehicle
func (r *vehicleRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Vehicle{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}
	return nil
}
