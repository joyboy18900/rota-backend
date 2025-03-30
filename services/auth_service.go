package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	authErrors "rota-api/errors"
	"rota-api/logs"
	"rota-api/models"
	"rota-api/repositories"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthService interface {
	// User operations
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, email, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)

	// Token operations
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateAccessToken(token string) (*TokenClaims, error)
	RefreshToken(ctx context.Context, refreshToken string) (accessToken string, newRefreshToken string, err error)
	Logout(ctx context.Context, userID string, token string) error

	// Token blacklist operations
	IsTokenBlacklisted(token string) bool
	AddToBlacklist(token string)

	// Password operations
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
}

// AuthConfig represents the configuration for AuthService
type AuthConfig struct {
	TokenConfig models.TokenConfig
	RedisConfig *repositories.RedisConfig
}

// AuthServiceImpl implements AuthService
type AuthServiceImpl struct {
	userRepo  repositories.UserRepository
	redisRepo *repositories.RedisRepository
	config    AuthConfig
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repositories.UserRepository,
	redisRepo *repositories.RedisRepository,
	config AuthConfig,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo:  userRepo,
		redisRepo: redisRepo,
		config:    config,
	}
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

// Register creates a new user
func (s *AuthServiceImpl) Register(ctx context.Context, user *models.User) error {
	// Check if email exists
	existingUser, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return authErrors.NewAuthError(authErrors.ErrEmailExists, "email already exists", nil)
	}

	// Check if username exists
	existingUser, err = s.userRepo.FindByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return authErrors.NewAuthError(authErrors.ErrUsernameExists, "username already exists", nil)
	}

	// Hash password
	hashedPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (*models.User, error) {
	logs.Info("attempting login", "email", email)

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		logs.Error("login failed", "error", err, "email", email)
		return nil, authErrors.NewAuthError(authErrors.ErrInvalidCredentials, "invalid credentials", err)
	}

	if !s.CheckPassword(password, user.Password) {
		logs.Error("invalid password", "email", email)
		return nil, authErrors.NewAuthError(authErrors.ErrInvalidCredentials, "invalid credentials", nil)
	}

	logs.Info("login successful", "user_id", user.ID, "email", email)

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

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

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthServiceImpl) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// GenerateAccessToken generates a new JWT access token
func (s *AuthServiceImpl) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.config.TokenConfig.ExpiryTime).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.TokenConfig.Secret))
}

// GenerateRefreshToken generates a new refresh token
func (s *AuthServiceImpl) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ValidateAccessToken validates the access token and returns the claims
func (s *AuthServiceImpl) ValidateAccessToken(token string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.TokenConfig.Secret), nil
	})

	if err != nil {
		return nil, authErrors.NewAuthError(authErrors.ErrInvalidToken, "invalid token", err)
	}

	if !parsedToken.Valid {
		return nil, authErrors.NewAuthError(authErrors.ErrTokenExpired, "token has expired", nil)
	}

	return claims, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (accessToken string, newRefreshToken string, err error) {
	// Find user by refresh token
	user, err := s.userRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	// Generate new access token
	accessToken, err = s.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err = s.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update refresh token in database
	user.RefreshToken = newRefreshToken
	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", "", fmt.Errorf("failed to update refresh token: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

// Logout invalidates the user's tokens
func (s *AuthServiceImpl) Logout(ctx context.Context, userID string, token string) error {
	// Add access token to blacklist first
	if s.redisRepo != nil {
		s.AddToBlacklist(token)
	}

	// Find and update user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Clear refresh token
	user.RefreshToken = ""
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// GetJWTSecret returns the JWT secret
func (s *AuthServiceImpl) GetJWTSecret() string {
	return s.config.TokenConfig.Secret
}

// IsTokenBlacklisted checks if a token is blacklisted
func (s *AuthServiceImpl) IsTokenBlacklisted(token string) bool {
	if s.redisRepo == nil {
		return false // If Redis is not available, assume token is valid
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
