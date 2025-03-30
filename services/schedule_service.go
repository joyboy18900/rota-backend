package services

import (
	"context"

	"rota-api/models"
	"rota-api/repositories"
)

// ScheduleService interface defines methods for schedule service
type ScheduleService interface {
	GetScheduleByID(ctx context.Context, id uint) (*models.Schedule, error)
	GetAllSchedules(ctx context.Context) ([]*models.Schedule, error)
	CreateSchedule(ctx context.Context, schedule *models.Schedule) error
	UpdateSchedule(ctx context.Context, schedule *models.Schedule) error
	DeleteSchedule(ctx context.Context, id uint) error
}

// scheduleService implements ScheduleService
type scheduleService struct {
	scheduleRepo repositories.ScheduleRepository
}

// NewScheduleService creates a new schedule service
func NewScheduleService(scheduleRepo repositories.ScheduleRepository) ScheduleService {
	return &scheduleService{scheduleRepo}
}

// GetScheduleByID retrieves a schedule by ID
func (s *scheduleService) GetScheduleByID(ctx context.Context, id uint) (*models.Schedule, error) {
	return s.scheduleRepo.FindByID(ctx, id)
}

// GetAllSchedules retrieves all schedules
func (s *scheduleService) GetAllSchedules(ctx context.Context) ([]*models.Schedule, error) {
	return s.scheduleRepo.FindAll(ctx)
}

// CreateSchedule creates a new schedule
func (s *scheduleService) CreateSchedule(ctx context.Context, schedule *models.Schedule) error {
	return s.scheduleRepo.Create(ctx, schedule)
}

// UpdateSchedule updates a schedule
func (s *scheduleService) UpdateSchedule(ctx context.Context, schedule *models.Schedule) error {
	return s.scheduleRepo.Update(ctx, schedule)
}

// DeleteSchedule deletes a schedule
func (s *scheduleService) DeleteSchedule(ctx context.Context, id uint) error {
	return s.scheduleRepo.Delete(ctx, id)
}
