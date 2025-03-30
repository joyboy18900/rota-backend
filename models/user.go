package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID           string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username     string     `json:"username" gorm:"uniqueIndex;not null"`
	Email        string     `json:"email" gorm:"uniqueIndex;not null"`
	Password     string     `json:"-" gorm:"not null"`
	FirstName    string     `json:"firstName" gorm:"column:first_name"`
	LastName     string     `json:"lastName" gorm:"column:last_name"`
	PhoneNumber  string     `json:"phoneNumber" gorm:"column:phone_number"`
	Role         string     `json:"role" gorm:"default:user"`
	IsActive     bool       `json:"isActive" gorm:"column:is_active;default:true"`
	GoogleID     string     `json:"googleId" gorm:"column:google_id;uniqueIndex"`
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty" gorm:"column:last_login_at"`
	RefreshToken string     `json:"-" gorm:"column:refresh_token"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt    time.Time  `json:"-" gorm:"column:deleted_at;index"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// BeforeCreate is a hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Hash password before saving
	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

// BeforeUpdate is a hook that runs before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		// Hash password before updating
		hashedPassword, err := HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return nil
}
