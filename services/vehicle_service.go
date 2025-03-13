package services

import (
	"context"
	"fmt"

	"rota-api/models"
	"rota-api/repositories"
)

// VehicleService interface defines methods for vehicle service
type VehicleService interface {
	CreateVehicle(ctx context.Context, licensePlate string, capacity int, driverName string, routeID uint) (*models.Vehicle, error)
	GetVehicleByID(ctx context.Context, id uint) (*models.Vehicle, error)
	GetAllVehicles(ctx context.Context) ([]models.Vehicle, error)
	GetVehiclesByRoute(ctx context.Context, routeID uint) ([]models.Vehicle, error)
	UpdateVehicle(ctx context.Context, id uint, licensePlate string, capacity int, driverName string, routeID uint) (*models.Vehicle, error)
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
func (s *vehicleService) CreateVehicle(ctx context.Context, licensePlate string, capacity int, driverName string, routeID uint) (*models.Vehicle, error) {
	// Check if license plate already exists
	existingVehicle, err := s.vehicleRepo.FindByLicensePlate(ctx, licensePlate)
	if err == nil && existingVehicle != nil {
		return nil, fmt.Errorf("vehicle with license plate %s already exists", licensePlate)
	}

	vehicle := &models.Vehicle{
		LicensePlate: licensePlate,
		Capacity:     capacity,
		DriverName:   driverName,
		RouteID:      routeID,
	}

	if err := s.vehicleRepo.Create(ctx, vehicle); err != nil {
		return nil, fmt.Errorf("failed to create vehicle: %w", err)
	}

	return vehicle, nil
}

// GetVehicleByID retrieves a vehicle by ID
func (s *vehicleService) GetVehicleByID(ctx context.Context, id uint) (*models.Vehicle, error) {
	return s.vehicleRepo.FindByID(ctx, id)
}

// GetAllVehicles retrieves all vehicles
func (s *vehicleService) GetAllVehicles(ctx context.Context) ([]models.Vehicle, error) {
	return s.vehicleRepo.FindAll(ctx)
}

// GetVehiclesByRoute retrieves all vehicles for a specific route
func (s *vehicleService) GetVehiclesByRoute(ctx context.Context, routeID uint) ([]models.Vehicle, error) {
	return s.vehicleRepo.FindByRoute(ctx, routeID)
}

// UpdateVehicle updates a vehicle
func (s *vehicleService) UpdateVehicle(ctx context.Context, id uint, licensePlate string, capacity int, driverName string, routeID uint) (*models.Vehicle, error) {
	vehicle, err := s.vehicleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find vehicle: %w", err)
	}

	// Update fields if provided
	if licensePlate != "" && licensePlate != vehicle.LicensePlate {
		// Check if license plate already exists
		existingVehicle, err := s.vehicleRepo.FindByLicensePlate(ctx, licensePlate)
		if err == nil && existingVehicle != nil && existingVehicle.ID != id {
			return nil, fmt.Errorf("vehicle with license plate %s already exists", licensePlate)
		}
		vehicle.LicensePlate = licensePlate
	}
	if capacity != 0 {
		vehicle.Capacity = capacity
	}
	if driverName != "" {
		vehicle.DriverName = driverName
	}
	if routeID != 0 {
		vehicle.RouteID = routeID
	}

	if err := s.vehicleRepo.Update(ctx, vehicle); err != nil {
		return nil, fmt.Errorf("failed to update vehicle: %w", err)
	}

	return vehicle, nil
}

// DeleteVehicle deletes a vehicle
func (s *vehicleService) DeleteVehicle(ctx context.Context, id uint) error {
	return s.vehicleRepo.Delete(ctx, id)
}
