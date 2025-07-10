package dto

import (
	"time"
)

// FavoriteStationResponse represents the simplified favorite station data sent to clients
type FavoriteStationResponse struct {
	StationID uint      `json:"station_id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
}
