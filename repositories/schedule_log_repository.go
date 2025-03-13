package repositories

import (
	"context"
	"errors"
	"fmt"
	"rota-api/models"

	"gorm.io/gorm"
)

// ScheduleLogRepository interface defines methods for schedule log database operations
type ScheduleLogRepository interface {
	Create(ctx context.Context, log *models.ScheduleLog) error
	FindByID(ctx context.Context, id uint) (*models.ScheduleLog, error)
	FindAll(ctx context.Context) ([]models.ScheduleLog, error)
	FindBySchedule(ctx context.Context, scheduleID uint) ([]models.ScheduleLog, error)
	FindByStaff(ctx context.Context, staffID uint) ([]models.ScheduleLog, error)
}

// scheduleLogRepository implements ScheduleLogRepository
type scheduleLogRepository struct {
	db *gorm.DB
}

// NewScheduleLogRepository creates a new schedule log repository
func NewScheduleLogRepository(db *gorm.DB) ScheduleLogRepository {
	return &scheduleLogRepository{db}
}

// Create stores a new schedule log in the database
func (r *scheduleLogRepository) Create(ctx context.Context, log *models.ScheduleLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("failed to create schedule log: %w", err)
	}
	return nil
}

// FindByID retrieves a schedule log by ID
func (r *scheduleLogRepository) FindByID(ctx context.Context, id uint) (*models.ScheduleLog, error) {
	var log models.ScheduleLog
	if err := r.db.WithContext(ctx).Preload("Schedule").Preload("Staff").First(&log, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("schedule log not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find schedule log: %w", err)
	}
	return &log, nil
}

// FindAll retrieves all schedule logs
func (r *scheduleLogRepository) FindAll(ctx context.Context) ([]models.ScheduleLog, error) {
	var logs []models.ScheduleLog
	if err := r.db.WithContext(ctx).Preload("Schedule").Preload("Staff").Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to find schedule logs: %w", err)
	}
	return logs, nil
}

// FindBySchedule retrieves all schedule logs for a specific schedule
func (r *scheduleLogRepository) FindBySchedule(ctx context.Context, scheduleID uint) ([]models.ScheduleLog, error) {
	var logs []models.ScheduleLog
	if err := r.db.WithContext(ctx).Preload("Schedule").Preload("Staff").Where("schedule_id = ?", scheduleID).Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to find schedule logs: %w", err)
	}
	return logs, nil
}

// FindByStaff retrieves all schedule logs for a specific staff
func (r *scheduleLogRepository) FindByStaff(ctx context.Context, staffID uint) ([]models.ScheduleLog, error) {
	var logs []models.ScheduleLog
	if err := r.db.WithContext(ctx).Preload("Schedule").Preload("Staff").Where("staff_id = ?", staffID).Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to find schedule logs: %w", err)
	}
	return logs, nil
}
