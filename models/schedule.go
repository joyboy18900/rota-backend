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
