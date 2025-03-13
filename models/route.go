
package models

// Route represents a transit route between two stations
type Route struct {
	ID             uint    `gorm:"primarykey" json:"id"`
	StartStationID uint    `json:"start_station_id"`
	EndStationID   uint    `json:"end_station_id"`
	Distance       float64 `gorm:"not null" json:"distance"`
	Duration       string  `gorm:"not null" json:"duration"`
	// Relations
	StartStation   Station   `gorm:"foreignKey:StartStationID" json:"start_station,omitempty"`
	EndStation     Station   `gorm:"foreignKey:EndStationID" json:"end_station,omitempty"`
	Schedules      []Schedule `gorm:"foreignKey:RouteID" json:"-"`
	Vehicles       []Vehicle  `gorm:"foreignKey:RouteID" json:"-"`
}
