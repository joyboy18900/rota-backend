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
	Search(ctx context.Context, params models.ScheduleSearchParams) (models.PagedResult, error)
	Create(ctx context.Context, schedule *models.Schedule) error
	Update(ctx context.Context, schedule *models.Schedule) error
	Delete(ctx context.Context, id uint) error
	FindSchedulesByStation(ctx context.Context, stationID uint, limit int) (models.StationSchedulesResponse, error)
	FindSimpleSchedulesByStation(ctx context.Context, stationID uint) (*models.SimpleStationScheduleResponse, error)
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
	if err := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Route.StartStation").
		Preload("Route.EndStation").
		Preload("Vehicle").
		Preload("Station").
		First(&schedule, id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) FindAll(ctx context.Context) ([]*models.Schedule, error) {
	var schedules []*models.Schedule
	if err := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Route.StartStation").
		Preload("Route.EndStation").
		Preload("Vehicle").
		Preload("Station").
		Find(&schedules).Error; err != nil {
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

func (r *scheduleRepository) Search(ctx context.Context, params models.ScheduleSearchParams) (models.PagedResult, error) {
	var result models.PagedResult
	var schedules []*models.Schedule
	var totalCount int64

	// Initialize query
	query := r.db.WithContext(ctx).Model(&models.Schedule{})

	// Apply filters
	if params.RouteID != nil {
		query = query.Where("route_id = ?", *params.RouteID)
	}
	if params.VehicleID != nil {
		query = query.Where("vehicle_id = ?", *params.VehicleID)
	}
	if params.StationID != nil {
		query = query.Where("station_id = ?", *params.StationID)
	}
	if params.Status != nil && *params.Status != "" {
		query = query.Where("status = ?", *params.Status)
	}
	if params.Round != nil {
		query = query.Where("round = ?", *params.Round)
	}

	// Date range filters
	if params.StartDateFrom != nil {
		query = query.Where("departure_time >= ?", *params.StartDateFrom)
	}
	if params.StartDateTo != nil {
		query = query.Where("departure_time <= ?", *params.StartDateTo)
	}

	// Count total results (before pagination)
	if err := query.Count(&totalCount).Error; err != nil {
		return result, err
	}

	// Set default pagination if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Apply sorting
	if params.SortBy != "" {
		sortDirection := "asc"
		if params.SortDesc {
			sortDirection = "desc"
		}
		query = query.Order(params.SortBy + " " + sortDirection)
	} else {
		// Default sorting by departure_time
		query = query.Order("departure_time asc")
	}

	// Get paginated results with preloaded relations
	if err := query.
		Preload("Route").
		Preload("Route.StartStation").
		Preload("Route.EndStation").
		Preload("Vehicle").
		Preload("Station").
		Offset(offset).
		Limit(params.PageSize).
		Find(&schedules).Error; err != nil {
		return result, err
	}

	// Calculate total pages
	totalPages := int(totalCount) / params.PageSize
	if int(totalCount)%params.PageSize > 0 {
		totalPages++
	}

	// Prepare result
	result = models.PagedResult{
		TotalCount: totalCount,
		TotalPages: totalPages,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Data:       schedules,
	}

	return result, nil
}

func (r *scheduleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Schedule{}, id).Error
}

// FindSchedulesByStation returns both inbound and outbound schedules for a specific station
// with detailed station information and limited to a specific number of schedules per direction
// limit parameter specifies how many schedules to return for EACH direction (e.g., limit=10 means 10 outbound + 10 inbound)
func (r *scheduleRepository) FindSchedulesByStation(ctx context.Context, stationID uint, limit int) (models.StationSchedulesResponse, error) {
	var response models.StationSchedulesResponse
	
	// Get station details
	var station models.Station
	if err := r.db.WithContext(ctx).First(&station, stationID).Error; err != nil {
		return response, err
	}
	response.Station = station

	// Find all routes where this station is either start or end
	var routeIDs []uint
	if err := r.db.WithContext(ctx).
		Model(&models.Route{}).
		Where("start_station_id = ? OR end_station_id = ?", stationID, stationID).
		Pluck("id", &routeIDs).Error; err != nil {
		return response, err
	}

	// Get outbound schedules (from this station to others)
	var outboundSchedules []*models.Schedule
	outboundQuery := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Route.StartStation").
		Preload("Route.EndStation").
		Preload("Vehicle").
		Preload("Station").
		Where("station_id = ?", stationID).
		Where("route_id IN ?", routeIDs).
		Order("departure_time asc").
		Limit(limit)

	if err := outboundQuery.Find(&outboundSchedules).Error; err != nil {
		return response, err
	}

	// Get inbound schedules (from others to this station)
	var inboundSchedules []*models.Schedule
	// First get routes where this station is the end station
	var inboundRouteIDs []uint
	if err := r.db.WithContext(ctx).
		Model(&models.Route{}).
		Where("end_station_id = ?", stationID).
		Pluck("id", &inboundRouteIDs).Error; err != nil {
		return response, err
	}

	// Then get schedules for those routes, where the station is the end station's paired station
	inboundQuery := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Route.StartStation").
		Preload("Route.EndStation").
		Preload("Vehicle").
		Preload("Station").
		Where("route_id IN ?", inboundRouteIDs).
		Order("departure_time asc").
		Limit(limit)

	if err := inboundQuery.Find(&inboundSchedules).Error; err != nil {
		return response, err
	}

	response.OutboundSchedules = outboundSchedules
	response.InboundSchedules = inboundSchedules

	// Add station details (example description)
	response.StationDetails = "สถานีขนส่งผู้โดยสารสถานี" + station.Name + " ให้บริการเดินรถระหว่างเมือง ตั้งอยู่ที่ " + station.Location

	return response, nil
}

// FindSimpleSchedulesByStation returns a simplified version of schedules for a station
// with only departure times and destinations, limited to 10 schedules in each direction
func (r *scheduleRepository) FindSimpleSchedulesByStation(ctx context.Context, stationID uint) (*models.SimpleStationScheduleResponse, error) {
	response := &models.SimpleStationScheduleResponse{}
	
	// Get station details
	var station models.Station
	if err := r.db.WithContext(ctx).First(&station, stationID).Error; err != nil {
		return nil, err
	}
	response.StationName = station.Name
	
	// Get outbound schedules (from this station to others)
	var outboundSchedules []*models.Schedule
	outboundQuery := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Route.EndStation").
		Where("station_id = ?", stationID).
		Order("departure_time asc").
		Limit(10)
	
	if err := outboundQuery.Find(&outboundSchedules).Error; err != nil {
		return nil, err
	}
	
	// Convert to simple format
	outboundTimes := make([]models.SimpleScheduleInfo, 0, len(outboundSchedules))
	for _, schedule := range outboundSchedules {
		outboundTimes = append(outboundTimes, models.SimpleScheduleInfo{
			DepartureTime: schedule.DepartureTime.Format("15:04"),
			Destination:   schedule.Route.EndStation.Name,
		})
	}
	response.OutboundTimes = outboundTimes
	
	// Get inbound schedules (from others to this station)
	var inboundRouteIDs []uint
	if err := r.db.WithContext(ctx).
		Model(&models.Route{}).
		Where("end_station_id = ?", stationID).
		Pluck("id", &inboundRouteIDs).Error; err != nil {
		return nil, err
	}
	
	var inboundSchedules []*models.Schedule
	inboundQuery := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Route.StartStation").
		Where("route_id IN ?", inboundRouteIDs).
		Order("departure_time asc").
		Limit(10)
	
	if err := inboundQuery.Find(&inboundSchedules).Error; err != nil {
		return nil, err
	}
	
	// Convert to simple format
	inboundTimes := make([]models.SimpleScheduleInfo, 0, len(inboundSchedules))
	for _, schedule := range inboundSchedules {
		inboundTimes = append(inboundTimes, models.SimpleScheduleInfo{
			DepartureTime: schedule.DepartureTime.Format("15:04"),
			Destination:   schedule.Route.StartStation.Name,
		})
	}
	response.InboundTimes = inboundTimes
	
	// Add station details
	response.StationDetails = "สถานีขนส่งผู้โดยสารสถานี" + station.Name + " ให้บริการเดินรถระหว่างเมือง ตั้งอยู่ที่ " + station.Location
	
	return response, nil
}
