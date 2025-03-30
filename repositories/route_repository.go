package repositories

import (
	"context"
	"errors"
	"fmt"
	"rota-api/models"

	"gorm.io/gorm"
)

// RouteRepository interface defines methods for route database operations
type RouteRepository interface {
	Create(ctx context.Context, route *models.Route) error
	FindByID(ctx context.Context, id uint) (*models.Route, error)
	FindAll(ctx context.Context) ([]*models.Route, error)
	FindByStation(ctx context.Context, stationID uint) ([]models.Route, error)
	Update(ctx context.Context, route *models.Route) error
	Delete(ctx context.Context, id uint) error
}

// routeRepository implements RouteRepository
type routeRepository struct {
	db *gorm.DB
}

// NewRouteRepository creates a new route repository
func NewRouteRepository(db *gorm.DB) RouteRepository {
	return &routeRepository{db}
}

// Create stores a new route in the database
func (r *routeRepository) Create(ctx context.Context, route *models.Route) error {
	if err := r.db.WithContext(ctx).Create(route).Error; err != nil {
		return fmt.Errorf("failed to create route: %w", err)
	}
	return nil
}

// FindByID retrieves a route by ID
func (r *routeRepository) FindByID(ctx context.Context, id uint) (*models.Route, error) {
	var route models.Route
	if err := r.db.WithContext(ctx).Preload("StartStation").Preload("EndStation").First(&route, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("route not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find route: %w", err)
	}
	return &route, nil
}

// FindAll retrieves all routes
func (r *routeRepository) FindAll(ctx context.Context) ([]*models.Route, error) {
	var routes []*models.Route
	if err := r.db.WithContext(ctx).Preload("StartStation").Preload("EndStation").Find(&routes).Error; err != nil {
		return nil, fmt.Errorf("failed to find routes: %w", err)
	}
	return routes, nil
}

// FindByStation retrieves all routes that start or end at a station
func (r *routeRepository) FindByStation(ctx context.Context, stationID uint) ([]models.Route, error) {
	var routes []models.Route
	if err := r.db.WithContext(ctx).
		Preload("StartStation").
		Preload("EndStation").
		Where("start_station_id = ? OR end_station_id = ?", stationID, stationID).
		Find(&routes).Error; err != nil {
		return nil, fmt.Errorf("failed to find routes: %w", err)
	}
	return routes, nil
}

// Update updates a route
func (r *routeRepository) Update(ctx context.Context, route *models.Route) error {
	if err := r.db.WithContext(ctx).Save(route).Error; err != nil {
		return fmt.Errorf("failed to update route: %w", err)
	}
	return nil
}

// Delete removes a route
func (r *routeRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Route{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete route: %w", err)
	}
	return nil
}
