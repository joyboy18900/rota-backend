package handler

import (
	"rota-api/models"
	"rota-api/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ScheduleHandler struct {
	scheduleService services.ScheduleService
}

// parseUintParam parses a string parameter to uint pointer
func parseUintParam(param string) (*uint, error) {
	if param == "" {
		return nil, nil
	}
	valueInt, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return nil, err
	}
	value := uint(valueInt)
	return &value, nil
}

// parseIntParam parses a string parameter to int pointer
func parseIntParam(param string) (*int, error) {
	if param == "" {
		return nil, nil
	}
	value, err := strconv.Atoi(param)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// parseTimeParam parses a string parameter to time.Time pointer
func parseTimeParam(param string) (*time.Time, error) {
	if param == "" {
		return nil, nil
	}
	// Support ISO8601 format (2025-06-02T08:00:00Z)
	value, err := time.Parse(time.RFC3339, param)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// parseStringParam returns a string pointer or nil if empty
func parseStringParam(param string) *string {
	if param == "" {
		return nil
	}
	return &param
}

func NewScheduleHandler(scheduleService services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

func (h *ScheduleHandler) GetScheduleByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schedule ID",
		})
	}

	schedule, err := h.scheduleService.GetScheduleByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Schedule not found",
		})
	}

	return c.JSON(fiber.Map{
		"schedule": schedule,
	})
}

func (h *ScheduleHandler) GetAllSchedules(c *fiber.Ctx) error {
	schedules, err := h.scheduleService.GetAllSchedules(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"schedules": schedules,
	})
}

func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	var schedule models.Schedule
	if err := c.BodyParser(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.scheduleService.CreateSchedule(c.Context(), &schedule); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Schedule created successfully",
		"schedule": schedule,
	})
}

func (h *ScheduleHandler) UpdateSchedule(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schedule ID",
		})
	}

	var schedule models.Schedule
	if err := c.BodyParser(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	schedule.ID = uint(id)
	if err := h.scheduleService.UpdateSchedule(c.Context(), &schedule); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Schedule updated successfully",
		"schedule": schedule,
	})
}

// SearchSchedules handles advanced search requests for schedules
func (h *ScheduleHandler) SearchSchedules(c *fiber.Ctx) error {
	// Extract search parameters from query params
	params := models.ScheduleSearchParams{}

	// Parse pagination parameters
	page, err := parseIntParam(c.Query("page"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page parameter",
		})
	}
	if page != nil {
		params.Page = *page
	}

	pageSize, err := parseIntParam(c.Query("page_size"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page_size parameter",
		})
	}
	if pageSize != nil {
		params.PageSize = *pageSize
	}

	// Sorting parameters
	params.SortBy = c.Query("sort_by")
	params.SortDesc = c.Query("sort_desc") == "true"

	// Filter parameters
	routeID, err := parseUintParam(c.Query("route_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid route_id parameter",
		})
	}
	params.RouteID = routeID

	vehicleID, err := parseUintParam(c.Query("vehicle_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vehicle_id parameter",
		})
	}
	params.VehicleID = vehicleID

	stationID, err := parseUintParam(c.Query("station_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid station_id parameter",
		})
	}
	params.StationID = stationID

	round, err := parseIntParam(c.Query("round"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid round parameter",
		})
	}
	params.Round = round

	// Date range parameters
	startDateFrom, err := parseTimeParam(c.Query("start_date_from"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start_date_from format. Use ISO8601 (e.g. 2025-06-02T08:00:00Z)",
		})
	}
	params.StartDateFrom = startDateFrom

	startDateTo, err := parseTimeParam(c.Query("start_date_to"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start_date_to format. Use ISO8601 (e.g. 2025-06-02T08:00:00Z)",
		})
	}
	params.StartDateTo = startDateTo

	// Status parameter
	params.Status = parseStringParam(c.Query("status"))

	// Perform search
	result, err := h.scheduleService.SearchSchedules(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return result
	return c.JSON(fiber.Map{
		"total_count": result.TotalCount,
		"total_pages": result.TotalPages,
		"page":       result.Page,
		"page_size":  result.PageSize,
		"schedules":  result.Data,
	})
}

func (h *ScheduleHandler) DeleteSchedule(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schedule ID",
		})
	}

	if err := h.scheduleService.DeleteSchedule(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Schedule deleted successfully",
	})
}

// GetSchedulesByStation retrieves schedules (both inbound and outbound) for a specific station
// Returns simplified schedule format with just station name, departure times and destinations
// Fixed at 10 schedules each direction
func (h *ScheduleHandler) GetSchedulesByStation(c *fiber.Ctx) error {
	// Parse station ID from path parameters
	stationID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid station ID",
		})
	}

	// Call service to get simplified schedules (10 per direction)
	response, err := h.scheduleService.GetSimpleSchedulesByStation(c.Context(), uint(stationID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Station or schedules not found: " + err.Error(),
		})
	}

	return c.JSON(response)
}

// GetSimpleSchedulesByStation retrieves a simplified version of schedules for a station
// with only departure times and destinations, limited to 10 schedules in each direction
func (h *ScheduleHandler) GetSimpleSchedulesByStation(c *fiber.Ctx) error {
	// Parse station ID from path parameters
	stationID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid station ID",
		})
	}

	// Call service to get simple schedules for this station
	// (Fixed at 10 schedules each direction as per requirement)
	response, err := h.scheduleService.GetSimpleSchedulesByStation(c.Context(), uint(stationID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Station or schedules not found: " + err.Error(),
		})
	}

	return c.JSON(response)
}
