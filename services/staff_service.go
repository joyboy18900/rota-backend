package services

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"rota-api/models"
	"rota-api/repositories"
)

// StaffService interface defines methods for staff service
type StaffService interface {
	CreateStaff(ctx context.Context, username, email, password string, stationID uint) (*models.Staff, error)
	GetStaffByID(ctx context.Context, id uint) (*models.Staff, error)
	GetAllStaff(ctx context.Context) ([]models.Staff, error)
	GetStaffByStation(ctx context.Context, stationID uint) ([]models.Staff, error)
	UpdateStaff(ctx context.Context, id uint, username, email, password string, stationID uint) (*models.Staff, error)
	DeleteStaff(ctx context.Context, id uint) error
}

// staffService implements StaffService
type staffService struct {
	staffRepo repositories.StaffRepository
}

// NewStaffService creates a new staff service
func NewStaffService(staffRepo repositories.StaffRepository) StaffService {
	return &staffService{staffRepo}
}

// CreateStaff creates a new staff
func (s *staffService) CreateStaff(ctx context.Context, username, email, password string, stationID uint) (*models.Staff, error) {
	// Check if email already exists
	existingStaff, err := s.staffRepo.FindByEmail(ctx, email)
	if err == nil && existingStaff != nil {
		return nil, fmt.Errorf("staff with email %s already exists", email)
	}

	// Check if username already exists
	existingStaff, err = s.staffRepo.FindByUsername(ctx, username)
	if err == nil && existingStaff != nil {
		return nil, fmt.Errorf("staff with username %s already exists", username)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	staff := &models.Staff{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		StationID: stationID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.staffRepo.Create(ctx, staff); err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	return staff, nil
}

// GetStaffByID retrieves a staff by ID
func (s *staffService) GetStaffByID(ctx context.Context, id uint) (*models.Staff, error) {
	return s.staffRepo.FindByID(ctx, id)
}

// GetAllStaff retrieves all staff
func (s *staffService) GetAllStaff(ctx context.Context) ([]models.Staff, error) {
	return s.staffRepo.FindAll(ctx)
}

// GetStaffByStation retrieves all staff for a specific station
func (s *staffService) GetStaffByStation(ctx context.Context, stationID uint) ([]models.Staff, error) {
	return s.staffRepo.FindByStation(ctx, stationID)
}

// UpdateStaff updates a staff
func (s *staffService) UpdateStaff(ctx context.Context, id uint, username, email, password string, stationID uint) (*models.Staff, error) {
	staff, err := s.staffRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find staff: %w", err)
	}

	// Update fields if provided
	if username != "" {
		staff.Username = username
	}
	if email != "" {
		staff.Email = email
	}
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		staff.Password = string(hashedPassword)
	}
	if stationID != 0 {
		staff.StationID = stationID
	}
	staff.UpdatedAt = time.Now()

	if err := s.staffRepo.Update(ctx, staff); err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}

	return staff, nil
}

// DeleteStaff deletes a staff
func (s *staffService) DeleteStaff(ctx context.Context, id uint) error {
	return s.staffRepo.Delete(ctx, id)
}
