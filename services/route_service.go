package services

import (
	"context"
	"fmt"

	"rota-api/models"
	"rota-api/repositories"
)

// RouteService interface defines methods for route service
type RouteService interface {
	GetRouteByID(ctx context.Context, id uint) (*models.Route, error)
	GetAllRoutes(ctx context.Context) ([]*models.Route, error)
	CreateRoute(ctx context.Context, route *models.Route) error
	UpdateRoute(ctx context.Context, route *models.Route) error
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
func (s *routeService) CreateRoute(ctx context.Context, route *models.Route) error {
	if err := s.routeRepo.Create(ctx, route); err != nil {
		return fmt.Errorf("failed to create route: %w", err)
	}
	return nil
}

// GetRouteByID retrieves a route by ID
func (s *routeService) GetRouteByID(ctx context.Context, id uint) (*models.Route, error) {
	return s.routeRepo.FindByID(ctx, id)
}

// GetAllRoutes retrieves all routes
func (s *routeService) GetAllRoutes(ctx context.Context) ([]*models.Route, error) {
	return s.routeRepo.FindAll(ctx)
}

// UpdateRoute updates a route
func (s *routeService) UpdateRoute(ctx context.Context, route *models.Route) error {

	existingRoute, err := s.routeRepo.FindByID(ctx, route.ID)
	if err != nil {
		return fmt.Errorf("failed to find route: %w", err)
	}


	if route.StartStationID != 0 {
		existingRoute.StartStationID = route.StartStationID
	}
	if route.EndStationID != 0 {
		existingRoute.EndStationID = route.EndStationID
	}
	if route.Distance != 0 {
		existingRoute.Distance = route.Distance
	}
	if route.Duration != "" {
		existingRoute.Duration = route.Duration
	}


	err = s.routeRepo.Update(ctx, existingRoute)
	if err != nil {
		return fmt.Errorf("failed to update route: %w", err)
	}
	

	*route = *existingRoute
	
	return nil
}

// DeleteRoute deletes a route
func (s *routeService) DeleteRoute(ctx context.Context, id uint) error {
	return s.routeRepo.Delete(ctx, id)
}
