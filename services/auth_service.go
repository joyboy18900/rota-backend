package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"rota-api/models"
	"rota-api/repositories"
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
}

// AuthConfig represents the configuration for AuthService
type AuthConfig struct {
	JWTSecret           string        `mapstructure:"jwt_secret"`
	JWTExpiration       time.Duration `mapstructure:"jwt_expiration"`
	RefreshExpiration   time.Duration `mapstructure:"refresh_expiration"`
	BlacklistExpiration time.Duration `mapstructure:"blacklist_expiration"`
}

// AuthServiceImpl implements AuthService
type AuthServiceImpl struct {
	userRepo repositories.UserRepository
	config   AuthConfig
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo repositories.UserRepository, config AuthConfig) AuthService {
	return &AuthServiceImpl{
		userRepo: userRepo,
		config:   config,
	}
}

// Register creates a new user with the provided information
func (s *AuthServiceImpl) Register(ctx *fiber.Ctx, user *models.User) (*models.User, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx.Context(), user.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Check if username is provided and already exists
	if user.Username != nil && *user.Username != "" {
		existingUser, err = s.userRepo.GetByUsername(ctx.Context(), *user.Username)
		if err == nil && existingUser != nil {
			return nil, errors.New("username already exists")
		}
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = models.RoleUser
	}

	// Hash password if provided
	if user.Password != nil && *user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		hashedPass := string(hashedPassword)
		user.Password = &hashedPass
	}

	// Create user in database
	if err := s.userRepo.Create(ctx.Context(), user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns the user and a JWT token
func (s *AuthServiceImpl) Login(ctx *fiber.Ctx, email, password string) (*models.User, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx.Context(), email)
	if err != nil || user == nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Check if user has a password (OAuth users might not have one)
	if user.Password == nil {
		return nil, "", errors.New("please use the appropriate login method")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.GenerateAccessToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
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
	// Parse the token with the claims
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthServiceImpl) RefreshToken(refreshToken string) (string, string, error) {
	// In a real implementation, you would validate the refresh token
	// and issue a new access token
	// For now, we'll just generate a new refresh token as well
	newRefreshToken, err := s.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// In a real implementation, you would:
	// 1. Validate the refresh token
	// 2. Get the user ID from the refresh token
	// 3. Get the user from the database
	// 4. Generate a new access token
	// 5. Return both tokens

	// For now, we'll just return empty strings for the access token
	// and the new refresh token
	return "", newRefreshToken, nil
}

// Logout invalidates the provided token
func (s *AuthServiceImpl) Logout(token string) error {
	// Add token to blacklist
	if err := s.AddToBlacklist(token, s.config.BlacklistExpiration); err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}
	return nil
}

// IsTokenBlacklisted checks if a token is in the blacklist
func (s *AuthServiceImpl) IsTokenBlacklisted(token string) bool {
	// In a real implementation, you would check if the token is in the blacklist
	// For simplicity, we'll just return false here
	return false
}

// AddToBlacklist adds a token to the blacklist
func (s *AuthServiceImpl) AddToBlacklist(token string, expiry time.Duration) error {
	// In a real implementation, you would add the token to a blacklist
	// with the specified expiration time
	// For example, using Redis:
	// return s.redisRepo.Set(ctx, fmt.Sprintf("blacklist:%s", token), true, expiry).Err()
	return nil
}

// GetJWTSecret returns the JWT secret key
func (s *AuthServiceImpl) GetJWTSecret() string {
	return s.config.JWTSecret
}

// GetUserByID retrieves a user by ID
func (s *AuthServiceImpl) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// HashPassword hashes a password using bcrypt
func (s *AuthServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword compares a password with its hash
func (s *AuthServiceImpl) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Register creates a new user with the provided information
func (s *AuthServiceImpl) Register(ctx *fiber.Ctx, user *models.User) (*models.User, error) {
	// Check if email exists
	existingUser, err := s.userRepo.FindByEmail(ctx.Context(), user.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Check if username exists
	if user.Username != nil {
		existingUser, err = s.userRepo.FindByUsername(ctx.Context(), *user.Username)
		if err == nil && existingUser != nil {
			return nil, fmt.Errorf("username already exists")
		}
	}

	// Hash password
	hashedPassword, err := s.HashPassword(*user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = &hashedPassword

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Create(ctx.Context(), user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns the user and a JWT token
func (s *AuthServiceImpl) Login(ctx *fiber.Ctx, email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(ctx.Context(), email)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	if user.Password == nil || !s.CheckPassword(password, *user.Password) {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Generate access token
	token, err := s.GenerateAccessToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx.Context(), user); err != nil {
		return nil, "", fmt.Errorf("failed to update user: %w", err)
	}

	return user, token, nil

	// Generate refresh token
	refreshToken, err := s.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save refresh token
	user.RefreshToken = refreshToken
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}
}

// GetUserByID retrieves a user by ID
func (s *AuthServiceImpl) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
}

// GenerateAccessToken generates a new JWT access token for the user
func (s *AuthServiceImpl) GenerateAccessToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(s.config.JWTExpiration)

	claims := &TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "rota-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}
	}

	exists, err := s.redisRepo.IsBlacklisted(context.Background(), token)
	if err != nil {
		logs.Error("error checking token blacklist", "error", err)
		return false // If Redis error, assume token is valid
	}
	return exists
}

// AddToBlacklist adds a token to the blacklist
func (s *AuthServiceImpl) AddToBlacklist(token string) {
	if s.redisRepo == nil {
		logs.Warn("Redis not available, token blacklist disabled")
		return // If Redis is not available, skip blacklisting
	}

	err := s.redisRepo.AddToBlacklist(context.Background(), token, s.config.TokenConfig.ExpiryTime)
	if err != nil {
		logs.Error("error adding token to blacklist", "error", err)
	}
}

/* TODO: Implement Google OAuth in the future
// getOAuthConfig returns the OAuth2 config for Google
func (s *AuthServiceImpl) getOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.config.GoogleClientID,
		ClientSecret: s.config.GoogleSecret,
		RedirectURL:  s.config.GoogleRedirect,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GoogleLogin returns the Google OAuth login URL
func (s *AuthServiceImpl) GoogleLogin() string {
	oauth2Config := s.getOAuthConfig()
	return oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

// GoogleCallback handles the Google OAuth callback
func (s *AuthServiceImpl) GoogleCallback(ctx context.Context, code string) (string, error) {
	oauth2Config := s.getOAuthConfig()

	// Exchange authorization code for token
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info
	userInfo, err := utils.GetGoogleUserInfo(ctx, token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if user exists
	user, err := s.userRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		// Create new user
		user = &models.User{
			Username:  userInfo.Name,
			Email:     userInfo.Email,
			FirstName: userInfo.GivenName,
			LastName:  userInfo.FamilyName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return "", fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Store OAuth token
	oauthToken := &models.OAuthToken{
		UserID:       user.ID,
		Provider:     "google",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.oauthRepo.Create(ctx, oauthToken); err != nil {
		return "", fmt.Errorf("failed to create OAuth token: %w", err)
	}

	// Generate JWT token
	return utils.GenerateJWT(user.ID, s.config.JWTSecret)
}
*/
