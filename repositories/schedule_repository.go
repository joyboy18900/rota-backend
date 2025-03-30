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
	if err := r.db.WithContext(ctx).First(&schedule, id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) FindAll(ctx context.Context) ([]*models.Schedule, error) {
	var schedules []*models.Schedule
	if err := r.db.WithContext(ctx).Find(&schedules).Error; err != nil {
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

func (r *scheduleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Schedule{}, id).Error
}
