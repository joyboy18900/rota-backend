package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRole defines the type for user roles
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleStaff UserRole = "staff"
	RoleAdmin UserRole = "admin"
)

// User represents a user in the system
type User struct {
	ID             int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Username       *string    `json:"username" gorm:"size:50;uniqueIndex"`
	Password       *string    `json:"-" gorm:"type:text"`
	Email          string     `json:"email" gorm:"size:100;not null;uniqueIndex"`
	Provider       string     `json:"provider" gorm:"size:20;not null;default:'local'"`
	ProviderID     *string    `json:"providerId,omitempty" gorm:"type:text"`
	ProfilePicture *string    `json:"profilePicture,omitempty" gorm:"type:text"`
	IsVerified     bool       `json:"isVerified" gorm:"not null;default:false"`
	// Role field is commented out as it doesn't exist in the database
	// Role           UserRole   `json:"role" gorm:"type:varchar(20);not null;default:'user'"`
	// LastLoginAt field is commented out as it doesn't exist in the database
	// LastLoginAt    *time.Time `json:"lastLoginAt" gorm:"default:null"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt      time.Time  `json:"updatedAt" gorm:"not null;default:now()"`
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

// CheckPassword verifies a password against the hashed password
func (u *User) CheckPassword(password string) error {
	if u.Password == nil {
		return bcrypt.ErrMismatchedHashAndPassword
	}
	return bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password))
}

// BeforeCreate is a hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Role field is removed from the model since it doesn't exist in the database
	// if u.Role == "" {
	// 	u.Role = RoleUser
	// }
	
	// Hash password if provided
	if u.Password != nil && *u.Password != "" {
		hashedPassword, err := HashPassword(*u.Password)
		if err != nil {
			return err
		}
		u.Password = &hashedPassword
	}
	
	// Set timestamps
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	
	return nil
}

// BeforeUpdate is a hook that runs before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") && u.Password != nil {
		// Hash password before updating
		hashedPassword, err := HashPassword(*u.Password)
		if err != nil {
			return err
		}
		u.Password = &hashedPassword
	}
	return nil
}
