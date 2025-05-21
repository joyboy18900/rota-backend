package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"rota-api/models"
	"rota-api/repositories"
)

// Custom errors
var (
	ErrEmailExists        = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenGeneration    = errors.New("failed to generate token")
	ErrUserNotFound       = errors.New("user not found")
)

// TokenClaims represents the JWT claims for authentication
type TokenClaims struct {
	UserID int             `json:"user_id"`
	Email  string          `json:"email"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// AuthService defines the interface for authentication operations
type AuthService interface {
	// User operations
	Register(ctx *fiber.Ctx, user *models.User) (*models.User, error)
	Login(ctx *fiber.Ctx, email, password string) (*models.User, string, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)

	// Token operations
	GenerateAccessToken(user *models.User) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateAccessToken(token string) (*TokenClaims, error)
	RefreshToken(refreshToken string) (string, string, error)
	Logout(token string) error

	// Token blacklist operations
	IsTokenBlacklisted(token string) bool
	AddToBlacklist(token string, expiry time.Duration) error

	// Utility methods
	GetJWTSecret() string
	GetJWTExpiration() time.Duration
	GetRefreshExpiration() time.Duration
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
}

// AuthConfig represents the configuration for AuthService
type AuthConfig struct {
	JWTSecret           string        `mapstructure:"jwt_secret"`
	JWTExpiration       time.Duration `mapstructure:"jwt_expiration"`
	RefreshExpiration   time.Duration `mapstructure:"refresh_expiration"`
	BlacklistExpiration time.Duration `mapstructure:"blacklist_expiration"`
	TokenConfig         models.TokenConfig
	RedisConfig         *repositories.RedisConfig
}

// AuthServiceImpl implements AuthService
type AuthServiceImpl struct {
	userRepo  repositories.UserRepository
	redisRepo repositories.RedisRepository
	config    AuthConfig
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(
	userRepo repositories.UserRepository,
	redisRepo repositories.RedisRepository,
	config AuthConfig,
) AuthService {
	return &AuthServiceImpl{
		userRepo:  userRepo,
		redisRepo: redisRepo,
		config:    config,
	}
}

// Register creates a new user with the provided information
func (s *AuthServiceImpl) Register(ctx *fiber.Ctx, user *models.User) (*models.User, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx.Context(), user.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailExists
	}

	// Check if username is provided and already exists
	if user.Username != nil && *user.Username != "" {
		existingUser, err = s.userRepo.FindByUsername(ctx.Context(), *user.Username)
		if err == nil && existingUser != nil {
			return nil, ErrUsernameExists
		}
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = models.RoleUser
	}

	// Hash password if provided
	if user.Password != nil && *user.Password != "" {
		hashedPassword, err := s.HashPassword(*user.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = &hashedPassword
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Create user in database
	if err := s.userRepo.Create(ctx.Context(), user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns the user and a JWT token
func (s *AuthServiceImpl) Login(ctx *fiber.Ctx, email, password string) (*models.User, string, error) {
	// Get user by email
	user, err := s.userRepo.FindByEmail(ctx.Context(), email)
	if err != nil || user == nil {
		return nil, "", ErrInvalidCredentials
	}

	// Check if user has a password (OAuth users might not have one)
	if user.Password == nil {
		return nil, "", errors.New("please use the appropriate login method")
	}

	// Verify password
	if !s.CheckPassword(password, *user.Password) {
		return nil, "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.GenerateAccessToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrTokenGeneration, err)
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	user.UpdatedAt = now
	if err := s.userRepo.Update(ctx.Context(), user); err != nil {
		log.Printf("Failed to update user last login time: %v", err)
		// Continue even if update fails - login is still successful
	}

	return user, token, nil
}

// GenerateAccessToken generates a new JWT access token for the user
func (s *AuthServiceImpl) GenerateAccessToken(user *models.User) (string, error) {
	// Set token claims
	claims := TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "rota-api",
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a new refresh token
func (s *AuthServiceImpl) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ValidateAccessToken validates the JWT token and returns the claims
func (s *AuthServiceImpl) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(s.config.JWTSecret), nil
	})

	// Check for errors
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	// Check if token is valid
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthServiceImpl) RefreshToken(refreshToken string) (string, string, error) {
	// Get user by refresh token
	user, err := s.userRepo.FindByRefreshToken(context.Background(), refreshToken)
	if err != nil || user == nil {
		return "", "", ErrInvalidToken
	}

	// Generate new access token
	accessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", ErrTokenGeneration, err)
	}

	// Generate new refresh token
	newRefreshToken, err := s.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", ErrTokenGeneration, err)
	}

	return accessToken, newRefreshToken, nil
}

// Logout invalidates the provided token
func (s *AuthServiceImpl) Logout(token string) error {
	return s.AddToBlacklist(token, s.config.BlacklistExpiration)
}

// IsTokenBlacklisted checks if a token is in the blacklist
func (s *AuthServiceImpl) IsTokenBlacklisted(token string) bool {
	// Check if redisRepo is nil using reflection
	if s.redisRepo == nil {
		log.Printf("Redis repository not initialized, token blacklist check skipped")
		return false
	}

	exists, err := s.redisRepo.IsBlacklisted(context.Background(), token)
	if err != nil {
		log.Printf("Error checking token blacklist: %v", err)
		return false // If there's an error, assume token is not blacklisted
	}
	return exists
}

// AddToBlacklist adds a token to the blacklist
func (s *AuthServiceImpl) AddToBlacklist(token string, expiry time.Duration) error {
	// Check if redisRepo is nil
	if s.redisRepo == nil {
		log.Printf("Redis repository not initialized, token blacklist not available")
		return nil
	}

	// Add token to blacklist with expiry
	if err := s.redisRepo.AddToBlacklist(context.Background(), token, expiry); err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}
	return nil
}

// GetJWTSecret returns the JWT secret key
func (s *AuthServiceImpl) GetJWTSecret() string {
	return s.config.JWTSecret
}

// GetJWTExpiration returns the JWT expiration duration
func (s *AuthServiceImpl) GetJWTExpiration() time.Duration {
	return s.config.JWTExpiration
}

// GetRefreshExpiration returns the refresh token expiration duration
func (s *AuthServiceImpl) GetRefreshExpiration() time.Duration {
	return s.config.RefreshExpiration
}

// GetUserByID retrieves a user by ID
func (s *AuthServiceImpl) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	// Convert int to string as FindByID expects string ID
	userID := fmt.Sprintf("%d", id)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}
	return user, nil
}

// HashPassword hashes a password using bcrypt
func (s *AuthServiceImpl) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a password with its hash
func (s *AuthServiceImpl) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
