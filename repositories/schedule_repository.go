package repositories

import (
	"context"

	"rota-api/models"

	"gorm.io/gorm"
)

// ScheduleRepository defines the interface for schedule-related database operations
type ScheduleRepository interface {
	FindByID(ctx context.Context, id uint) (*models.Schedule, error)
	FindAll(ctx context.Context) ([]*models.Schedule, error)
	Search(ctx context.Context, params models.ScheduleSearchParams) (models.PagedResult, error)
	Create(ctx context.Context, schedule *models.Schedule) error
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

func (r *scheduleRepository) FindByID(ctx context.Context, id uint) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Route.StartStation").
		Preload("Route.EndStation").
		Preload("Vehicle").
		Preload("Station").
		First(&schedule, id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) FindAll(ctx context.Context) ([]*models.Schedule, error) {
	var schedules []*models.Schedule
	if err := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Route.StartStation").
		Preload("Route.EndStation").
		Preload("Vehicle").
		Preload("Station").
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *scheduleRepository) Create(ctx context.Context, schedule *models.Schedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

func (r *scheduleRepository) Update(ctx context.Context, schedule *models.Schedule) error {
	return r.db.WithContext(ctx).Save(schedule).Error
}

func (r *scheduleRepository) Search(ctx context.Context, params models.ScheduleSearchParams) (models.PagedResult, error) {
	var result models.PagedResult
	var schedules []*models.Schedule
	var totalCount int64

	// Initialize query
	query := r.db.WithContext(ctx).Model(&models.Schedule{})

	// Apply filters
	if params.RouteID != nil {
		query = query.Where("route_id = ?", *params.RouteID)
	}
	if params.VehicleID != nil {
		query = query.Where("vehicle_id = ?", *params.VehicleID)
	}
	if params.StationID != nil {
		query = query.Where("station_id = ?", *params.StationID)
	}
	if params.Status != nil && *params.Status != "" {
		query = query.Where("status = ?", *params.Status)
	}
	if params.Round != nil {
		query = query.Where("round = ?", *params.Round)
	}

	// Date range filters
	if params.StartDateFrom != nil {
		query = query.Where("departure_time >= ?", *params.StartDateFrom)
	}
	if params.StartDateTo != nil {
		query = query.Where("departure_time <= ?", *params.StartDateTo)
	}

	// Count total results (before pagination)
	if err := query.Count(&totalCount).Error; err != nil {
		return result, err
	}

	// Set default pagination if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Apply sorting
	if params.SortBy != "" {
		sortDirection := "asc"
		if params.SortDesc {
			sortDirection = "desc"
		}
		query = query.Order(params.SortBy + " " + sortDirection)
	} else {
		// Default sorting by departure_time
		query = query.Order("departure_time asc")
	}

	// Get paginated results with preloaded relations
	if err := query.
		Preload("Route").
		Preload("Route.StartStation").
		Preload("Route.EndStation").
		Preload("Vehicle").
		Preload("Station").
		Offset(offset).
		Limit(params.PageSize).
		Find(&schedules).Error; err != nil {
		return result, err
	}

	// Calculate total pages
	totalPages := int(totalCount) / params.PageSize
	if int(totalCount)%params.PageSize > 0 {
		totalPages++
	}

	// Prepare result
	result = models.PagedResult{
		TotalCount: totalCount,
		TotalPages: totalPages,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Data:       schedules,
	}

	return result, nil
}

func (r *scheduleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Schedule{}, id).Error
}
