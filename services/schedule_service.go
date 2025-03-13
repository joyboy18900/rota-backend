package services

import (
	"context"
	"fmt"

	"rota-api/models"
	"rota-api/repositories"
)

// ScheduleService interface defines methods for schedule service
type ScheduleService interface {
	CreateSchedule(ctx context.Context, routeID, stationID uint, round int, departureTime, arrivalTime string) (*models.Schedule, error)
	GetScheduleByID(ctx context.Context, id uint) (*models.Schedule, error)
	GetAllSchedules(ctx context.Context) ([]models.Schedule, error)
	GetSchedulesByRoute(ctx context.Context, routeID uint) ([]models.Schedule, error)
	GetSchedulesByStation(ctx context.Context, stationID uint) ([]models.Schedule, error)
	UpdateSchedule(ctx context.Context, id, routeID, stationID uint, round int, departureTime, arrivalTime string) (*models.Schedule, error)
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

// CreateSchedule creates a new schedule
func (s *scheduleService) CreateSchedule(ctx context.Context, routeID, stationID uint, round int, departureTime, arrivalTime string) (*models.Schedule, error) {
	schedule := &models.Schedule{
		RouteID:       routeID,
		StationID:     stationID,
		Round:         round,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
	}

	if err := s.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return schedule, nil
}

// GetScheduleByID retrieves a schedule by ID
func (s *scheduleService) GetScheduleByID(ctx context.Context, id uint) (*models.Schedule, error) {
	return s.scheduleRepo.FindByID(ctx, id)
}

// GetAllSchedules retrieves all schedules
func (s *scheduleService) GetAllSchedules(ctx context.Context) ([]models.Schedule, error) {
	return s.scheduleRepo.FindAll(ctx)
}

// GetSchedulesByRoute retrieves all schedules for a specific route
func (s *scheduleService) GetSchedulesByRoute(ctx context.Context, routeID uint) ([]models.Schedule, error) {
	return s.scheduleRepo.FindByRoute(ctx, routeID)
}

// GetSchedulesByStation retrieves all schedules for a specific station
func (s *scheduleService) GetSchedulesByStation(ctx context.Context, stationID uint) ([]models.Schedule, error) {
	return s.scheduleRepo.FindByStation(ctx, stationID)
}

// UpdateSchedule updates a schedule
func (s *scheduleService) UpdateSchedule(ctx context.Context, id, routeID, stationID uint, round int, departureTime, arrivalTime string) (*models.Schedule, error) {
	schedule, err := s.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find schedule: %w", err)
	}

	// Update fields if provided
	if routeID != 0 {
		schedule.RouteID = routeID
	}
	if stationID != 0 {
		schedule.StationID = stationID
	}
	if round != 0 {
		schedule.Round = round
	}
	if departureTime != "" {
		schedule.DepartureTime = departureTime
	}
	if arrivalTime != "" {
		schedule.ArrivalTime = arrivalTime
	}

	if err := s.scheduleRepo.Update(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	return schedule, nil
}

// DeleteSchedule deletes a schedule
func (s *scheduleService) DeleteSchedule(ctx context.Context, id uint) error {
	return s.scheduleRepo.Delete(ctx, id)
}
