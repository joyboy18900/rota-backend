package models

import "time"

// RouteVehicle represents the many-to-many relationship between routes and vehicles
type RouteVehicle struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RouteID   string    `json:"route_id" gorm:"type:uuid;not null"`
	VehicleID string    `json:"vehicle_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Route   Route   `json:"route,omitempty" gorm:"foreignKey:RouteID"`
	Vehicle Vehicle `json:"vehicle,omitempty" gorm:"foreignKey:VehicleID"`
}

// TableName specifies the table name for RouteVehicle
func (RouteVehicle) TableName() string {
	return "route_vehicles"
}
