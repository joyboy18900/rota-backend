package repositories

import (
	"context"
	"errors"
	"fmt"
	"rota-api/models"

	"gorm.io/gorm"
)

// OAuthTokenRepository interface defines methods for OAuth token database operations
type OAuthTokenRepository interface {
	Create(ctx context.Context, token *models.OAuthToken) error
	FindByUserAndProvider(ctx context.Context, userID uint, provider string) (*models.OAuthToken, error)
	Update(ctx context.Context, token *models.OAuthToken) error
	Delete(ctx context.Context, id uint) error
}

// oauthTokenRepository implements OAuthTokenRepository
type oauthTokenRepository struct {
	db *gorm.DB
}

// NewOAuthTokenRepository creates a new OAuth token repository
func NewOAuthTokenRepository(db *gorm.DB) OAuthTokenRepository {
	return &oauthTokenRepository{db}
}

// Create stores a new OAuth token in the database
func (r *oauthTokenRepository) Create(ctx context.Context, token *models.OAuthToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("failed to create OAuth token: %w", err)
	}
	return nil
}

// FindByUserAndProvider retrieves an OAuth token by user ID and provider
func (r *oauthTokenRepository) FindByUserAndProvider(ctx context.Context, userID uint, provider string) (*models.OAuthToken, error) {
	var token models.OAuthToken
	if err := r.db.WithContext(ctx).Where("user_id = ? AND provider = ?", userID, provider).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("OAuth token not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find OAuth token: %w", err)
	}
	return &token, nil
}

// Update updates an OAuth token
func (r *oauthTokenRepository) Update(ctx context.Context, token *models.OAuthToken) error {
	if err := r.db.WithContext(ctx).Save(token).Error; err != nil {
		return fmt.Errorf("failed to update OAuth token: %w", err)
	}
	return nil
}

// Delete removes an OAuth token
func (r *oauthTokenRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.OAuthToken{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete OAuth token: %w", err)
	}
	return nil
}
