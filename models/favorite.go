
package models

import (
	"time"
)

// Favorite represents a user's favorite station
type Favorite struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `json:"user_id"`
	StationID uint      `json:"station_id"`
	CreatedAt time.Time `json:"created_at"`
	// Relations
	User    User    `gorm:"foreignKey:UserID" json:"-"` // ไม่ส่งข้อมูล User กลับเนื่องจากไม่จำเป็น
	Station Station `gorm:"foreignKey:StationID" json:"station,omitempty"`
}
