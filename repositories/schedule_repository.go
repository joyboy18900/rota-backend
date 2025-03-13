package repositories

import (
	"context"
	"errors"
	"fmt"
	"rota-api/models"

	"gorm.io/gorm"
)

// ScheduleRepository interface defines methods for schedule database operations
type ScheduleRepository interface {
	Create(ctx context.Context, schedule *models.Schedule) error
	FindByID(ctx context.Context, id uint) (*models.Schedule, error)
	FindAll(ctx context.Context) ([]models.Schedule, error)
	FindByRoute(ctx context.Context, routeID uint) ([]models.Schedule, error)
	FindByStation(ctx context.Context, stationID uint) ([]models.Schedule, error)
	Update(ctx context.Context, schedule *models.Schedule) error
	Delete(ctx context.Context, id uint) error
}

// scheduleRepository implements ScheduleRepository
type scheduleRepository struct {
	db *gorm.DB
}

// NewScheduleRepository creates a new schedule repository
func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{db}
}

// Create stores a new schedule in the database
func (r *scheduleRepository) Create(ctx context.Context, schedule *models.Schedule) error {
	if err := r.db.WithContext(ctx).Create(schedule).Error; err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}
	return nil
}

// FindByID retrieves a schedule by ID
func (r *scheduleRepository) FindByID(ctx context.Context, id uint) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := r.db.WithContext(ctx).Preload("Route").Preload("Station").First(&schedule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("schedule not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find schedule: %w", err)
	}
	return &schedule, nil
}

// FindAll retrieves all schedules
func (r *scheduleRepository) FindAll(ctx context.Context) ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := r.db.WithContext(ctx).Preload("Route").Preload("Station").Find(&schedules).Error; err != nil {
		return nil, fmt.Errorf("failed to find schedules: %w", err)
	}
	return schedules, nil
}

// FindByRoute retrieves all schedules for a specific route
func (r *scheduleRepository) FindByRoute(ctx context.Context, routeID uint) ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := r.db.WithContext(ctx).Preload("Route").Preload("Station").Where("route_id = ?", routeID).Find(&schedules).Error; err != nil {
		return nil, fmt.Errorf("failed to find schedules: %w", err)
	}
	return schedules, nil
}

// FindByStation retrieves all schedules for a specific station
func (r *scheduleRepository) FindByStation(ctx context.Context, stationID uint) ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := r.db.WithContext(ctx).Preload("Route").Preload("Station").Where("station_id = ?", stationID).Find(&schedules).Error; err != nil {
		return nil, fmt.Errorf("failed to find schedules: %w", err)
	}
	return schedules, nil
}

// Update updates a schedule
func (r *scheduleRepository) Update(ctx context.Context, schedule *models.Schedule) error {
	if err := r.db.WithContext(ctx).Save(schedule).Error; err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}
	return nil
}

// Delete removes a schedule
func (r *scheduleRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Schedule{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}
	return nil
}
