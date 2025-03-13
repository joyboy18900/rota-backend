package models

// Schedule represents a transit schedule for a specific route and station
type Schedule struct {
	ID            uint   `gorm:"primarykey" json:"id"`
	RouteID       uint   `json:"route_id"`
	StationID     uint   `json:"station_id"`
	Round         int    `gorm:"not null" json:"round"`
	DepartureTime string `gorm:"not null" json:"departure_time"` // Using string for time representation
	ArrivalTime   string `gorm:"not null" json:"arrival_time"`   // Using string for time representation
	// Relations
	Route        Route         `gorm:"foreignKey:RouteID" json:"route,omitempty"`
	Station      Station       `gorm:"foreignKey:StationID" json:"station,omitempty"`
	ScheduleLogs []ScheduleLog `gorm:"foreignKey:ScheduleID" json:"-"`
}
