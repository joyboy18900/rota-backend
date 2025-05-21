package models

import (
	"time"

	"gorm.io/gorm"
)

type Stop struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Latitude    float64        `json:"latitude"`
	Longitude   float64        `json:"longitude"`
	Routes      []Route        `json:"routes" gorm:"many2many:route_stops"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Stop) TableName() string {
	return "stops"
}
