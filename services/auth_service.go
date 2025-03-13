package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"rota-api/config"
	"rota-api/models"
	"rota-api/repositories"
	"rota-api/utils"
)

// AuthService interface defines methods for authentication service
type AuthService interface {
	Register(ctx context.Context, username, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GoogleLogin() string
	GoogleCallback(ctx context.Context, code string) (string, error)
	ValidateToken(token string) (uint, error)
	RefreshToken(refreshToken string) (string, error)
}

// authService implements AuthService
type authService struct {
	userRepo       repositories.UserRepository
	oauthTokenRepo repositories.OAuthTokenRepository
	config         *config.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repositories.UserRepository, oauthTokenRepo repositories.OAuthTokenRepository, config *config.Config) AuthService {
	return &authService{
		userRepo:       userRepo,
		oauthTokenRepo: oauthTokenRepo,
		config:         config,
	}
}

// Register creates a new user
func (s *authService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Check if username already exists
	existingUser, err = s.userRepo.FindByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with username %s already exists", username)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Username:   username,
		Email:      email,
		Password:   string(hashedPassword),
		Provider:   "local",
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	// Check if it's a social login account without a password
	if user.Provider != "local" && user.Password == "" {
		return "", fmt.Errorf("please login with %s", user.Provider)
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, s.config.JWTSecret, s.config.TokenExpiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// getOAuthConfig returns the OAuth2 config for Google
func (s *authService) getOAuthConfig() *oauth2.Config {
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
func (s *authService) GoogleLogin() string {
	oauth2Config := s.getOAuthConfig()
	return oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

// GoogleCallback handles the Google OAuth callback
func (s *authService) GoogleCallback(ctx context.Context, code string) (string, error) {
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
	user, err := s.userRepo.FindByProviderID(ctx, "google", userInfo.ID)
	if err != nil {
		// User doesn't exist, check if email exists
		user, err = s.userRepo.FindByEmail(ctx, userInfo.Email)
		if err == nil {
			// Email exists, update provider info
			user.Provider = "google"
			user.ProviderID = userInfo.ID
			user.IsVerified = true
			user.UpdatedAt = time.Now()
			if err := s.userRepo.Update(ctx, user); err != nil {
				return "", fmt.Errorf("failed to update user: %w", err)
			}
		} else {
			// Create new user
			user = &models.User{
				Username:       userInfo.Name,
				Email:          userInfo.Email,
				Provider:       "google",
				ProviderID:     userInfo.ID,
				ProfilePicture: userInfo.Picture,
				IsVerified:     true,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			if err := s.userRepo.Create(ctx, user); err != nil {
				return "", fmt.Errorf("failed to create user: %w", err)
			}
		}
	}

	// Store OAuth token
	oauthToken := &models.OAuthToken{
		UserID:       user.ID,
		Provider:     "google",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}

	// Check if token exists
	existingToken, err := s.oauthTokenRepo.FindByUserAndProvider(ctx, user.ID, "google")
	if err == nil {
		// Update existing token
		existingToken.AccessToken = token.AccessToken
		existingToken.RefreshToken = token.RefreshToken
		existingToken.ExpiresAt = token.Expiry
		if err := s.oauthTokenRepo.Update(ctx, existingToken); err != nil {
			return "", fmt.Errorf("failed to update OAuth token: %w", err)
		}
	} else {
		// Create new token
		if err := s.oauthTokenRepo.Create(ctx, oauthToken); err != nil {
			return "", fmt.Errorf("failed to create OAuth token: %w", err)
		}
	}

	// Generate JWT token
	jwtToken, err := utils.GenerateJWT(user.ID, s.config.JWTSecret, s.config.TokenExpiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return jwtToken, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *authService) ValidateToken(tokenString string) (uint, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if token is expired
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return 0, fmt.Errorf("token expired")
			}
		}

		// Extract user ID
		if userID, ok := claims["user_id"].(float64); ok {
			return uint(userID), nil
		}
		return 0, fmt.Errorf("invalid user ID in token")
	}

	return 0, fmt.Errorf("invalid token")
}

// RefreshToken refreshes a JWT token
func (s *authService) RefreshToken(refreshToken string) (string, error) {
	// For simplicity, we'll just validate the token and issue a new one
	// In a real application, you'd use a dedicated refresh token mechanism
	userID, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate new token
	newToken, err := utils.GenerateJWT(userID, s.config.JWTSecret, s.config.TokenExpiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return newToken, nil
}
