
package models

import (
	"time"
)

// OAuthToken represents a token used for OAuth authentication
type OAuthToken struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `json:"user_id"`
	Provider     string    `gorm:"not null" json:"provider"`
	AccessToken  string    `gorm:"not null" json:"-"` // Hide token in json responses
	RefreshToken string    `json:"-"`                 // Hide token in json responses
	ExpiresAt    time.Time `json:"expires_at"`
	User         User      `gorm:"foreignKey:UserID" json:"-"`
}
