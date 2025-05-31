package models

import (
	"time"

	"gorm.io/gorm"
)

// Schedule represents a transit schedule for a specific route and station
type Schedule struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	RouteID       uint           `json:"route_id"`
	VehicleID     uint           `json:"vehicle_id"`
	StationID     uint           `json:"station_id"`
	Round         int            `gorm:"default:1" json:"round"`
	DepartureTime time.Time      `json:"departure_time"`
	ArrivalTime   time.Time      `json:"arrival_time"`
	Status        string         `gorm:"default:'scheduled'" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	// Relations
	Route        Route         `gorm:"foreignKey:RouteID" json:"route,omitempty"`
	Vehicle      Vehicle       `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Station      Station       `gorm:"foreignKey:StationID" json:"station,omitempty"`
	ScheduleLogs []ScheduleLog `gorm:"foreignKey:ScheduleID" json:"-"`
}

// StationSchedulesResponse represents the response for a station's schedule inquiry
// Contains station details and both outbound and inbound schedules
type StationSchedulesResponse struct {
	Station           Station      `json:"station"`             // The station details
	OutboundSchedules []*Schedule  `json:"outbound_schedules"` // Schedules leaving from this station
	InboundSchedules  []*Schedule  `json:"inbound_schedules"`  // Schedules arriving at this station
	StationDetails    string       `json:"station_details"`    // Additional information about the station
}

// SimpleScheduleInfo represents a simplified view of a schedule with just essential information
type SimpleScheduleInfo struct {
	DepartureTime string `json:"departure_time"` // Departure time in readable format
	Destination   string `json:"destination"`    // Destination station name
}

// SimpleStationScheduleResponse represents a simplified response for station schedules
// Contains only the essential information needed for display
type SimpleStationScheduleResponse struct {
	StationName     string              `json:"station_name"`     // Name of the station
	OutboundTimes   []SimpleScheduleInfo `json:"outbound_times"`  // Simplified outbound schedules
	InboundTimes    []SimpleScheduleInfo `json:"inbound_times"`   // Simplified inbound schedules
	StationDetails  string              `json:"station_details"` // Brief description of the station
}
