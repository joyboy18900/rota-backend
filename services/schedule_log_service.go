package services

import (
	"context"
	"fmt"
	"time"

	"rota-api/models"
	"rota-api/repositories"
)

// ScheduleLogService interface defines methods for schedule log service
type ScheduleLogService interface {
	CreateScheduleLog(ctx context.Context, scheduleID, staffID uint, changeDescription string) (*models.ScheduleLog, error)
	GetScheduleLogByID(ctx context.Context, id uint) (*models.ScheduleLog, error)
	GetAllScheduleLogs(ctx context.Context) ([]models.ScheduleLog, error)
	GetScheduleLogsBySchedule(ctx context.Context, scheduleID uint) ([]models.ScheduleLog, error)
	GetScheduleLogsByStaff(ctx context.Context, staffID uint) ([]models.ScheduleLog, error)
}

// scheduleLogService implements ScheduleLogService
type scheduleLogService struct {
	scheduleLogRepo repositories.ScheduleLogRepository
}

// NewScheduleLogService creates a new schedule log service
func NewScheduleLogService(scheduleLogRepo repositories.ScheduleLogRepository) ScheduleLogService {
	return &scheduleLogService{scheduleLogRepo}
}

// CreateScheduleLog creates a new schedule log
func (s *scheduleLogService) CreateScheduleLog(ctx context.Context, scheduleID, staffID uint, changeDescription string) (*models.ScheduleLog, error) {
	log := &models.ScheduleLog{
		ScheduleID:        scheduleID,
		StaffID:           staffID,
		ChangeDescription: changeDescription,
		UpdatedAt:         time.Now(),
	}

	if err := s.scheduleLogRepo.Create(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to create schedule log: %w", err)
	}

	return log, nil
}

// GetScheduleLogByID retrieves a schedule log by ID
func (s *scheduleLogService) GetScheduleLogByID(ctx context.Context, id uint) (*models.ScheduleLog, error) {
	return s.scheduleLogRepo.FindByID(ctx, id)
}

// GetAllScheduleLogs retrieves all schedule logs
func (s *scheduleLogService) GetAllScheduleLogs(ctx context.Context) ([]models.ScheduleLog, error) {
	return s.scheduleLogRepo.FindAll(ctx)
}

// GetScheduleLogsBySchedule retrieves all schedule logs for a specific schedule
func (s *scheduleLogService) GetScheduleLogsBySchedule(ctx context.Context, scheduleID uint) ([]models.ScheduleLog, error) {
	return s.scheduleLogRepo.FindBySchedule(ctx, scheduleID)
}

// GetScheduleLogsByStaff retrieves all schedule logs for a specific staff
func (s *scheduleLogService) GetScheduleLogsByStaff(ctx context.Context, staffID uint) ([]models.ScheduleLog, error) {
	return s.scheduleLogRepo.FindByStaff(ctx, staffID)
}
