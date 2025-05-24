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
	SearchSchedules(ctx context.Context, params models.ScheduleSearchParams) (models.PagedResult, error)
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
	// ดึงข้อมูลเดิมก่อน
	existingSchedule, err := s.scheduleRepo.FindByID(ctx, schedule.ID)
	if err != nil {
		return err
	}

	// อัพเดทเฉพาะฟิลด์ที่ระบุ
	if schedule.RouteID != 0 {
		existingSchedule.RouteID = schedule.RouteID
	}
	if schedule.VehicleID != 0 {
		existingSchedule.VehicleID = schedule.VehicleID
	}
	if schedule.StationID != 0 {
		existingSchedule.StationID = schedule.StationID
	}
	if schedule.Round != 0 {
		existingSchedule.Round = schedule.Round
	}
	if !schedule.DepartureTime.IsZero() {
		existingSchedule.DepartureTime = schedule.DepartureTime
	}
	if !schedule.ArrivalTime.IsZero() {
		existingSchedule.ArrivalTime = schedule.ArrivalTime
	}
	if schedule.Status != "" {
		existingSchedule.Status = schedule.Status
	}

	// บันทึกการอัพเดท
	err = s.scheduleRepo.Update(ctx, existingSchedule)
	if err != nil {
		return err
	}

	// คัดลอกข้อมูลที่อัพเดทแล้วกลับไปยังพารามิเตอร์
	*schedule = *existingSchedule

	return nil
}

// SearchSchedules searches for schedules with advanced filtering
func (s *scheduleService) SearchSchedules(ctx context.Context, params models.ScheduleSearchParams) (models.PagedResult, error) {
	return s.scheduleRepo.Search(ctx, params)
}

// DeleteSchedule deletes a schedule
func (s *scheduleService) DeleteSchedule(ctx context.Context, id uint) error {
	return s.scheduleRepo.Delete(ctx, id)
}
