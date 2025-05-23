package services

import (
	"context"
	"fmt"

	"rota-api/models"
	"rota-api/repositories"
)

// VehicleService interface defines methods for vehicle service
type VehicleService interface {
	GetVehicleByID(ctx context.Context, id uint) (*models.Vehicle, error)
	GetAllVehicles(ctx context.Context) ([]*models.Vehicle, error)
	CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error
	UpdateVehicle(ctx context.Context, vehicle *models.Vehicle) error
	DeleteVehicle(ctx context.Context, id uint) error
}

// vehicleService implements VehicleService
type vehicleService struct {
	vehicleRepo repositories.VehicleRepository
}

// NewVehicleService creates a new vehicle service
func NewVehicleService(vehicleRepo repositories.VehicleRepository) VehicleService {
	return &vehicleService{vehicleRepo}
}

// CreateVehicle creates a new vehicle
func (s *vehicleService) CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error {
	// Check if license plate already exists
	existingVehicle, err := s.vehicleRepo.FindByLicensePlate(ctx, vehicle.LicensePlate)
	if err == nil && existingVehicle != nil {
		return fmt.Errorf("vehicle with license plate %s already exists", vehicle.LicensePlate)
	}

	if err := s.vehicleRepo.Create(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to create vehicle: %w", err)
	}

	return nil
}

// GetVehicleByID retrieves a vehicle by ID
func (s *vehicleService) GetVehicleByID(ctx context.Context, id uint) (*models.Vehicle, error) {
	return s.vehicleRepo.FindByID(ctx, id)
}

// GetAllVehicles retrieves all vehicles
func (s *vehicleService) GetAllVehicles(ctx context.Context) ([]*models.Vehicle, error) {
	return s.vehicleRepo.FindAll(ctx)
}

// UpdateVehicle updates a vehicle
func (s *vehicleService) UpdateVehicle(ctx context.Context, vehicle *models.Vehicle) error {
	existingVehicle, err := s.vehicleRepo.FindByID(ctx, vehicle.ID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	// Update fields if provided
	if vehicle.LicensePlate != "" && vehicle.LicensePlate != existingVehicle.LicensePlate {
		// Check if license plate already exists
		vehicleWithSamePlate, err := s.vehicleRepo.FindByLicensePlate(ctx, vehicle.LicensePlate)
		if err == nil && vehicleWithSamePlate != nil && vehicleWithSamePlate.ID != vehicle.ID {
			return fmt.Errorf("vehicle with license plate %s already exists", vehicle.LicensePlate)
		}
		existingVehicle.LicensePlate = vehicle.LicensePlate
	}
	if vehicle.Capacity != 0 {
		existingVehicle.Capacity = vehicle.Capacity
	}
	if vehicle.DriverName != "" {
		existingVehicle.DriverName = vehicle.DriverName
	}
	if vehicle.RouteID != 0 {
		existingVehicle.RouteID = vehicle.RouteID
	}

	if err := s.vehicleRepo.Update(ctx, existingVehicle); err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}

	return nil
}

// DeleteVehicle deletes a vehicle
func (s *vehicleService) DeleteVehicle(ctx context.Context, id uint) error {
	return s.vehicleRepo.Delete(ctx, id)
}
