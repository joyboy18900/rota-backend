package services

import (
	"context"

	"rota-api/models"
	"rota-api/repositories"
)

// ScheduleLogService interface defines methods for schedule log service
type ScheduleLogService interface {
	GetScheduleLogByID(ctx context.Context, id uint) (*models.ScheduleLog, error)
	GetAllScheduleLogs(ctx context.Context) ([]*models.ScheduleLog, error)
	CreateScheduleLog(ctx context.Context, scheduleLog *models.ScheduleLog) error
	UpdateScheduleLog(ctx context.Context, scheduleLog *models.ScheduleLog) error
	DeleteScheduleLog(ctx context.Context, id uint) error
}

// scheduleLogService implements ScheduleLogService
type scheduleLogService struct {
	scheduleLogRepo repositories.ScheduleLogRepository
}

// NewScheduleLogService creates a new schedule log service
func NewScheduleLogService(scheduleLogRepo repositories.ScheduleLogRepository) ScheduleLogService {
	return &scheduleLogService{scheduleLogRepo}
}

// GetScheduleLogByID retrieves a schedule log by ID
func (s *scheduleLogService) GetScheduleLogByID(ctx context.Context, id uint) (*models.ScheduleLog, error) {
	return s.scheduleLogRepo.FindByID(ctx, id)
}

// GetAllScheduleLogs retrieves all schedule logs
func (s *scheduleLogService) GetAllScheduleLogs(ctx context.Context) ([]*models.ScheduleLog, error) {
	return s.scheduleLogRepo.FindAll(ctx)
}

// CreateScheduleLog creates a new schedule log
func (s *scheduleLogService) CreateScheduleLog(ctx context.Context, scheduleLog *models.ScheduleLog) error {
	return s.scheduleLogRepo.Create(ctx, scheduleLog)
}

// UpdateScheduleLog updates a schedule log
func (s *scheduleLogService) UpdateScheduleLog(ctx context.Context, scheduleLog *models.ScheduleLog) error {
	return s.scheduleLogRepo.Update(ctx, scheduleLog)
}

// DeleteScheduleLog deletes a schedule log
func (s *scheduleLogService) DeleteScheduleLog(ctx context.Context, id uint) error {
	return s.scheduleLogRepo.Delete(ctx, id)
}
