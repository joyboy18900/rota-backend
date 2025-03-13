
package models

// Station represents a transit station
type Station struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Location string `json:"location"`
	// Relations
	StartRoutes      []Route       `gorm:"foreignKey:StartStationID" json:"-"`
	EndRoutes        []Route       `gorm:"foreignKey:EndStationID" json:"-"`
	Schedules        []Schedule    `gorm:"foreignKey:StationID" json:"-"`
	Staff            []Staff       `gorm:"foreignKey:StationID" json:"-"`
	Favorites        []Favorite    `gorm:"foreignKey:StationID" json:"-"`
}
