package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"rota-api/models"
	"rota-api/repositories"
)

// UserService interface defines methods for user service
type UserService interface {
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
}

// userService implements UserService
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo}
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	// Convert int to string as repository expects string ID
	strID := fmt.Sprintf("%d", id)
	return s.userRepo.FindByID(ctx, strID)
}

// GetAllUsers retrieves all users
func (s *userService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.FindAll(ctx)
}

// UpdateUser updates a user with all provided fields
func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	// Make sure password is hashed if it's provided
	if user.Password != nil && *user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		hashedPwdStr := string(hashedPassword)
		user.Password = &hashedPwdStr
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, id int) error {
	// Convert int to string as repository expects string ID
	strID := fmt.Sprintf("%d", id)
	return s.userRepo.Delete(ctx, strID)
}
