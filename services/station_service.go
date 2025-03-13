package services

import (
	"context"
	"fmt"

	"rota-api/models"
	"rota-api/repositories"
)

// StationService interface defines methods for station service
type StationService interface {
	CreateStation(ctx context.Context, name, location string) (*models.Station, error)
	GetStationByID(ctx context.Context, id uint) (*models.Station, error)
	GetAllStations(ctx context.Context) ([]models.Station, error)
	UpdateStation(ctx context.Context, id uint, name, location string) (*models.Station, error)
	DeleteStation(ctx context.Context, id uint) error
}

// stationService implements StationService
type stationService struct {
	stationRepo repositories.StationRepository
}

// NewStationService creates a new station service
func NewStationService(stationRepo repositories.StationRepository) StationService {
	return &stationService{stationRepo}
}

// CreateStation creates a new station
func (s *stationService) CreateStation(ctx context.Context, name, location string) (*models.Station, error) {
	station := &models.Station{
		Name:     name,
		Location: location,
	}

	if err := s.stationRepo.Create(ctx, station); err != nil {
		return nil, fmt.Errorf("failed to create station: %w", err)
	}

	return station, nil
}

// GetStationByID retrieves a station by ID
func (s *stationService) GetStationByID(ctx context.Context, id uint) (*models.Station, error) {
	return s.stationRepo.FindByID(ctx, id)
}

// GetAllStations retrieves all stations
func (s *stationService) GetAllStations(ctx context.Context) ([]models.Station, error) {
	return s.stationRepo.FindAll(ctx)
}

// UpdateStation updates a station
func (s *stationService) UpdateStation(ctx context.Context, id uint, name, location string) (*models.Station, error) {
	station, err := s.stationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find station: %w", err)
	}

	// Update fields if provided
	if name != "" {
		station.Name = name
	}
	if location != "" {
		station.Location = location
	}

	if err := s.stationRepo.Update(ctx, station); err != nil {
		return nil, fmt.Errorf("failed to update station: %w", err)
	}

	return station, nil
}

// DeleteStation deletes a station
func (s *stationService) DeleteStation(ctx context.Context, id uint) error {
	return s.stationRepo.Delete(ctx, id)
}
