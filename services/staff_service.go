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
	GetStaffByID(ctx context.Context, id uint) (*models.Staff, error)
	GetAllStaff(ctx context.Context) ([]*models.Staff, error)
	CreateStaff(ctx context.Context, staff *models.Staff) error
	UpdateStaff(ctx context.Context, staff *models.Staff) error
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
func (s *staffService) CreateStaff(ctx context.Context, staff *models.Staff) error {
	// Check if email already exists
	existingStaff, err := s.staffRepo.FindByEmail(ctx, staff.Email)
	if err == nil && existingStaff != nil {
		return fmt.Errorf("staff with email %s already exists", staff.Email)
	}

	// Check if username already exists
	existingStaff, err = s.staffRepo.FindByUsername(ctx, staff.Username)
	if err == nil && existingStaff != nil {
		return fmt.Errorf("staff with username %s already exists", staff.Username)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(staff.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	staff.Password = string(hashedPassword)
	staff.CreatedAt = time.Now()
	staff.UpdatedAt = time.Now()

	if err := s.staffRepo.Create(ctx, staff); err != nil {
		return fmt.Errorf("failed to create staff: %w", err)
	}

	return nil
}

// GetStaffByID retrieves a staff by ID
func (s *staffService) GetStaffByID(ctx context.Context, id uint) (*models.Staff, error) {
	return s.staffRepo.FindByID(ctx, id)
}

// GetAllStaff retrieves all staff
func (s *staffService) GetAllStaff(ctx context.Context) ([]*models.Staff, error) {
	return s.staffRepo.FindAll(ctx)
}

// UpdateStaff updates a staff
func (s *staffService) UpdateStaff(ctx context.Context, staff *models.Staff) error {
	existingStaff, err := s.staffRepo.FindByID(ctx, staff.ID)
	if err != nil {
		return fmt.Errorf("failed to find staff: %w", err)
	}

	// Update fields if provided
	if staff.Username != "" {
		existingStaff.Username = staff.Username
	}
	if staff.Email != "" {
		existingStaff.Email = staff.Email
	}
	if staff.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(staff.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		existingStaff.Password = string(hashedPassword)
	}
	if staff.StationID != 0 {
		existingStaff.StationID = staff.StationID
	}
	existingStaff.UpdatedAt = time.Now()

	if err := s.staffRepo.Update(ctx, existingStaff); err != nil {
		return fmt.Errorf("failed to update staff: %w", err)
	}

	return nil
}

// DeleteStaff deletes a staff
func (s *staffService) DeleteStaff(ctx context.Context, id uint) error {
	return s.staffRepo.Delete(ctx, id)
}
