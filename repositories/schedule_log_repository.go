package repositories

import (
	"context"

	"rota-api/models"

	"gorm.io/gorm"
)

// ScheduleLogRepository defines the interface for schedule log-related database operations
type ScheduleLogRepository interface {
	FindByID(ctx context.Context, id uint) (*models.ScheduleLog, error)
	FindAll(ctx context.Context) ([]*models.ScheduleLog, error)
	Create(ctx context.Context, scheduleLog *models.ScheduleLog) error
	Update(ctx context.Context, scheduleLog *models.ScheduleLog) error
	Delete(ctx context.Context, id uint) error
}

// scheduleLogRepository implements ScheduleLogRepository
type scheduleLogRepository struct {
	db *gorm.DB
}

// NewScheduleLogRepository creates a new schedule log repository
func NewScheduleLogRepository(db *gorm.DB) ScheduleLogRepository {
	return &scheduleLogRepository{db}
}

func (r *scheduleLogRepository) FindByID(ctx context.Context, id uint) (*models.ScheduleLog, error) {
	var scheduleLog models.ScheduleLog
	if err := r.db.WithContext(ctx).First(&scheduleLog, id).Error; err != nil {
		return nil, err
	}
	return &scheduleLog, nil
}

func (r *scheduleLogRepository) FindAll(ctx context.Context) ([]*models.ScheduleLog, error) {
	var scheduleLogs []*models.ScheduleLog
	if err := r.db.WithContext(ctx).Find(&scheduleLogs).Error; err != nil {
		return nil, err
	}
	return scheduleLogs, nil
}

func (r *scheduleLogRepository) Create(ctx context.Context, scheduleLog *models.ScheduleLog) error {
	return r.db.WithContext(ctx).Create(scheduleLog).Error
}

func (r *scheduleLogRepository) Update(ctx context.Context, scheduleLog *models.ScheduleLog) error {
	return r.db.WithContext(ctx).Save(scheduleLog).Error
}

func (r *scheduleLogRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ScheduleLog{}, id).Error
}
