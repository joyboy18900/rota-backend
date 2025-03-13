
package models

import (
	"time"

	"gorm.io/gorm"
)

// Staff represents a staff member
type Staff struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"unique;not null" json:"username"`
	Password  string         `json:"-"` // Hide password in json responses
	Email     string         `gorm:"unique;not null" json:"email"`
	StationID uint           `json:"station_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// Relations
	Station      Station       `gorm:"foreignKey:StationID" json:"station,omitempty"`
	ScheduleLogs []ScheduleLog `gorm:"foreignKey:StaffID" json:"-"`
}
