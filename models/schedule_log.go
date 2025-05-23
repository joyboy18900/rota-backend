
package models

import (
	"time"
)

// ScheduleLog represents a log entry for schedule changes
type ScheduleLog struct {
	ID                uint       `gorm:"primarykey" json:"id"`
	ScheduleID        uint       `json:"schedule_id"`
	StaffID           uint       `json:"staff_id"`
	ChangeDescription string     `gorm:"not null" json:"change_description"`
	ActualDeparture   *time.Time `gorm:"default:null" json:"actual_departure,omitempty"`
	ActualArrival     *time.Time `gorm:"default:null" json:"actual_arrival,omitempty"`
	Status            string     `json:"status,omitempty"`
	Notes             string     `json:"notes,omitempty"`
	UpdatedAt         *time.Time `gorm:"default:null" json:"updated_at,omitempty"`
	// Relations
	Schedule          Schedule  `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	Staff             Staff     `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
}
