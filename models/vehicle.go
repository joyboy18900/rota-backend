
package models

// Vehicle represents a transportation vehicle
type Vehicle struct {
	ID           uint   `gorm:"primarykey" json:"id"`
	LicensePlate string `gorm:"unique;not null" json:"license_plate"`
	Capacity     int    `gorm:"not null" json:"capacity"`
	DriverName   string `gorm:"not null" json:"driver_name"`
	RouteID      uint   `json:"route_id"`
	// Relations
	Route       Route  `gorm:"foreignKey:RouteID" json:"route,omitempty"`
}
