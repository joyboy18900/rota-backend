
package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Username        string         `gorm:"unique" json:"username"`
	Password        string         `json:"-"` // Hide password in json responses
	Email           string         `gorm:"unique;not null" json:"email"`
	Provider        string         `gorm:"default:local" json:"provider"`
	ProviderID      string         `json:"provider_id,omitempty"`
	ProfilePicture  string         `json:"profile_picture,omitempty"`
	IsVerified      bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	OAuthTokens     []OAuthToken   `gorm:"foreignKey:UserID" json:"-"`
	Favorites       []Favorite     `gorm:"foreignKey:UserID" json:"-"`
}
