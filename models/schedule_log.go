
package models

import (
	"time"
)

// ScheduleLog represents a log entry for schedule changes
type ScheduleLog struct {
	ID                uint      `gorm:"primarykey" json:"id"`
	ScheduleID        uint      `json:"schedule_id"`
	StaffID           uint      `json:"staff_id"`
	ChangeDescription string    `gorm:"not null" json:"change_description"`
	UpdatedAt         time.Time `json:"updated_at"`
	// Relations
	Schedule          Schedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	Staff             Staff    `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
}
