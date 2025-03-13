package services

import (
	"context"
	"fmt"

	"rota-api/models"
	"rota-api/repositories"
)

// RouteService interface defines methods for route service
type RouteService interface {
	CreateRoute(ctx context.Context, startStationID, endStationID uint, distance float64, duration string) (*models.Route, error)
	GetRouteByID(ctx context.Context, id uint) (*models.Route, error)
	GetAllRoutes(ctx context.Context) ([]models.Route, error)
	GetRoutesByStation(ctx context.Context, stationID uint) ([]models.Route, error)
	UpdateRoute(ctx context.Context, id, startStationID, endStationID uint, distance float64, duration string) (*models.Route, error)
	DeleteRoute(ctx context.Context, id uint) error
}

// routeService implements RouteService
type routeService struct {
	routeRepo   repositories.RouteRepository
	stationRepo repositories.StationRepository
}

// NewRouteService creates a new route service
func NewRouteService(routeRepo repositories.RouteRepository) RouteService {
	return &routeService{routeRepo: routeRepo}
}

// CreateRoute creates a new route
func (s *routeService) CreateRoute(ctx context.Context, startStationID, endStationID uint, distance float64, duration string) (*models.Route, error) {
	route := &models.Route{
		StartStationID: startStationID,
		EndStationID:   endStationID,
		Distance:       distance,
		Duration:       duration,
	}

	if err := s.routeRepo.Create(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	return route, nil
}

// GetRouteByID retrieves a route by ID
func (s *routeService) GetRouteByID(ctx context.Context, id uint) (*models.Route, error) {
	return s.routeRepo.FindByID(ctx, id)
}

// GetAllRoutes retrieves all routes
func (s *routeService) GetAllRoutes(ctx context.Context) ([]models.Route, error) {
	return s.routeRepo.FindAll(ctx)
}

// GetRoutesByStation retrieves all routes that start or end at a station
func (s *routeService) GetRoutesByStation(ctx context.Context, stationID uint) ([]models.Route, error) {
	return s.routeRepo.FindByStation(ctx, stationID)
}

// UpdateRoute updates a route
func (s *routeService) UpdateRoute(ctx context.Context, id, startStationID, endStationID uint, distance float64, duration string) (*models.Route, error) {
	route, err := s.routeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find route: %w", err)
	}

	// Update fields if provided
	if startStationID != 0 {
		route.StartStationID = startStationID
	}
	if endStationID != 0 {
		route.EndStationID = endStationID
	}
	if distance != 0 {
		route.Distance = distance
	}
	if duration != "" {
		route.Duration = duration
	}

	if err := s.routeRepo.Update(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to update route: %w", err)
	}

	return route, nil
}

// DeleteRoute deletes a route
func (s *routeService) DeleteRoute(ctx context.Context, id uint) error {
	return s.routeRepo.Delete(ctx, id)
}
