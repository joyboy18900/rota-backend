package models

import (
	"time"
)

// SearchParams defines common search parameters
type SearchParams struct {
	Page     int `json:"page" query:"page"`
	PageSize int `json:"page_size" query:"page_size"`
	SortBy   string `json:"sort_by" query:"sort_by"`
	SortDesc bool `json:"sort_desc" query:"sort_desc"`
}

// ScheduleSearchParams defines parameters for searching schedules
type ScheduleSearchParams struct {
	SearchParams
	RouteID       *uint      `json:"route_id" query:"route_id"`
	VehicleID     *uint      `json:"vehicle_id" query:"vehicle_id"`
	StationID     *uint      `json:"station_id" query:"station_id"`
	Status        *string    `json:"status" query:"status"`
	StartDateFrom *time.Time `json:"start_date_from" query:"start_date_from"`
	StartDateTo   *time.Time `json:"start_date_to" query:"start_date_to"`
	Round         *int       `json:"round" query:"round"`
}

// RouteSearchParams defines parameters for searching routes
type RouteSearchParams struct {
	SearchParams
	StartStationID *uint    `json:"start_station_id" query:"start_station_id"`
	EndStationID   *uint    `json:"end_station_id" query:"end_station_id"`
	MinDistance    *float64 `json:"min_distance" query:"min_distance"`
	MaxDistance    *float64 `json:"max_distance" query:"max_distance"`
}

// VehicleSearchParams defines parameters for searching vehicles
type VehicleSearchParams struct {
	SearchParams
	RouteID      *uint `json:"route_id" query:"route_id"`
	MinCapacity  *int  `json:"min_capacity" query:"min_capacity"`
	MaxCapacity  *int  `json:"max_capacity" query:"max_capacity"`
	LicensePlate *string `json:"license_plate" query:"license_plate"`
}

// StaffSearchParams defines parameters for searching staff
type StaffSearchParams struct {
	SearchParams
	StationID *uint   `json:"station_id" query:"station_id"`
	Position  *string `json:"position" query:"position"`
	Name      *string `json:"name" query:"name"`
}

// PagedResult represents a generic paged result
type PagedResult struct {
	TotalCount int64       `json:"total_count"`
	TotalPages int         `json:"total_pages"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Data       interface{} `json:"data"`
}
