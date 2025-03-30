package services

import (
	"context"
	"fmt"

	"rota-api/models"
	"rota-api/repositories"
)

// StationService interface defines methods for station service
type StationService interface {
	GetStationByID(ctx context.Context, id uint) (*models.Station, error)
	GetAllStations(ctx context.Context) ([]*models.Station, error)
	CreateStation(ctx context.Context, station *models.Station) error
	UpdateStation(ctx context.Context, station *models.Station) error
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
func (s *stationService) CreateStation(ctx context.Context, station *models.Station) error {
	if err := s.stationRepo.Create(ctx, station); err != nil {
		return fmt.Errorf("failed to create station: %w", err)
	}
	return nil
}

// GetStationByID retrieves a station by ID
func (s *stationService) GetStationByID(ctx context.Context, id uint) (*models.Station, error) {
	return s.stationRepo.FindByID(ctx, id)
}

// GetAllStations retrieves all stations
func (s *stationService) GetAllStations(ctx context.Context) ([]*models.Station, error) {
	return s.stationRepo.FindAll(ctx)
}

// UpdateStation updates a station
func (s *stationService) UpdateStation(ctx context.Context, station *models.Station) error {
	if err := s.stationRepo.Update(ctx, station); err != nil {
		return fmt.Errorf("failed to update station: %w", err)
	}
	return nil
}

// DeleteStation deletes a station
func (s *stationService) DeleteStation(ctx context.Context, id uint) error {
	return s.stationRepo.Delete(ctx, id)
}
