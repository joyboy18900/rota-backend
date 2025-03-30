package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// TokenConfig represents the configuration for token generation
type TokenConfig struct {
	Secret     string
	ExpiryTime time.Duration
}
