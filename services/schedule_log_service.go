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
	// ดึงข้อมูลเดิมก่อน
	existingLog, err := s.scheduleLogRepo.FindByID(ctx, scheduleLog.ID)
	if err != nil {
		return err
	}

	// อัพเดทเฉพาะฟิลด์ที่ระบุ
	if scheduleLog.ScheduleID != 0 {
		existingLog.ScheduleID = scheduleLog.ScheduleID
	}
	if scheduleLog.StaffID != 0 {
		existingLog.StaffID = scheduleLog.StaffID
	}
	if scheduleLog.ChangeDescription != "" {
		existingLog.ChangeDescription = scheduleLog.ChangeDescription
	}
	if !scheduleLog.ActualDeparture.IsZero() {
		existingLog.ActualDeparture = scheduleLog.ActualDeparture
	}
	if !scheduleLog.ActualArrival.IsZero() {
		existingLog.ActualArrival = scheduleLog.ActualArrival
	}
	if scheduleLog.Status != "" {
		existingLog.Status = scheduleLog.Status
	}
	if scheduleLog.Notes != "" {
		existingLog.Notes = scheduleLog.Notes
	}

	// บันทึกการอัพเดท
	err = s.scheduleLogRepo.Update(ctx, existingLog)
	if err != nil {
		return err
	}

	// คัดลอกข้อมูลที่อัพเดทแล้วกลับไปยังพารามิเตอร์
	*scheduleLog = *existingLog

	return nil
}

// DeleteScheduleLog deletes a schedule log
func (s *scheduleLogService) DeleteScheduleLog(ctx context.Context, id uint) error {
	return s.scheduleLogRepo.Delete(ctx, id)
}
