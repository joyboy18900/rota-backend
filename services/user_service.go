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
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateUser(ctx context.Context, id string, username, email, password string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
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
func (s *userService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// GetAllUsers retrieves all users
func (s *userService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.FindAll(ctx)
}

// UpdateUser updates a user
func (s *userService) UpdateUser(ctx context.Context, id string, username, email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Update fields if provided
	if username != "" {
		user.Username = &username
	}
	if email != "" {
		user.Email = email
	}
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		hashedPwdStr := string(hashedPassword)
		user.Password = &hashedPwdStr
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
