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
	GetSchedulesByStation(ctx context.Context, stationID uint, limit int) (models.StationSchedulesResponse, error)
	GetSimpleSchedulesByStation(ctx context.Context, stationID uint) (*models.SimpleStationScheduleResponse, error)
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
	// Get original data first
	existingSchedule, err := s.scheduleRepo.FindByID(ctx, schedule.ID)
	if err != nil {
		return err
	}

	// Update only specified fields
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

	// Save the update
	err = s.scheduleRepo.Update(ctx, existingSchedule)
	if err != nil {
		return err
	}

	// Copy updated data back to the parameter
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

// GetSchedulesByStation retrieves schedules (both inbound and outbound) for a specific station
// with a limit of schedules per direction
func (s *scheduleService) GetSchedulesByStation(ctx context.Context, stationID uint, limit int) (models.StationSchedulesResponse, error) {
	response, err := s.scheduleRepo.FindSchedulesByStation(ctx, stationID, limit)
	
	// If station details is empty, provide a default description based on station name
	if response.StationDetails == "" && response.Station.ID > 0 {
		response.StationDetails = "สถานีขนส่งผู้โดยสาร" + response.Station.Name + " ให้บริการเดินรถระหว่างเมือง ตั้งอยู่ที่ " + response.Station.Location
	}
	
	return response, err
}

// GetSimpleSchedulesByStation retrieves a simplified version of schedules for a station
// with only departure times and destinations, limited to 10 schedules in each direction
func (s *scheduleService) GetSimpleSchedulesByStation(ctx context.Context, stationID uint) (*models.SimpleStationScheduleResponse, error) {
	return s.scheduleRepo.FindSimpleSchedulesByStation(ctx, stationID)
}
