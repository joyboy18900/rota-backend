package repositories

import (
	"context"
	"errors"
	"fmt"
	"rota-api/models"

	"gorm.io/gorm"
)

// StaffRepository interface defines methods for staff database operations
type StaffRepository interface {
	Create(ctx context.Context, staff *models.Staff) error
	FindByID(ctx context.Context, id uint) (*models.Staff, error)
	FindByEmail(ctx context.Context, email string) (*models.Staff, error)
	FindByUsername(ctx context.Context, username string) (*models.Staff, error)
	FindAll(ctx context.Context) ([]models.Staff, error)
	FindByStation(ctx context.Context, stationID uint) ([]models.Staff, error)
	Update(ctx context.Context, staff *models.Staff) error
	Delete(ctx context.Context, id uint) error
}

// staffRepository implements StaffRepository
type staffRepository struct {
	db *gorm.DB
}

// NewStaffRepository creates a new staff repository
func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &staffRepository{db}
}

// Create stores a new staff in the database
func (r *staffRepository) Create(ctx context.Context, staff *models.Staff) error {
	if err := r.db.WithContext(ctx).Create(staff).Error; err != nil {
		return fmt.Errorf("failed to create staff: %w", err)
	}
	return nil
}

// FindByID retrieves a staff by ID
func (r *staffRepository) FindByID(ctx context.Context, id uint) (*models.Staff, error) {
	var staff models.Staff
	if err := r.db.WithContext(ctx).Preload("Station").First(&staff, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("staff not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find staff: %w", err)
	}
	return &staff, nil
}

// FindByEmail retrieves a staff by email
func (r *staffRepository) FindByEmail(ctx context.Context, email string) (*models.Staff, error) {
	var staff models.Staff
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("staff not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find staff: %w", err)
	}
	return &staff, nil
}

// FindByUsername retrieves a staff by username
func (r *staffRepository) FindByUsername(ctx context.Context, username string) (*models.Staff, error) {
	var staff models.Staff
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("staff not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find staff: %w", err)
	}
	return &staff, nil
}

// FindAll retrieves all staff
func (r *staffRepository) FindAll(ctx context.Context) ([]models.Staff, error) {
	var staff []models.Staff
	if err := r.db.WithContext(ctx).Preload("Station").Find(&staff).Error; err != nil {
		return nil, fmt.Errorf("failed to find staff: %w", err)
	}
	return staff, nil
}

// FindByStation retrieves all staff for a specific station
func (r *staffRepository) FindByStation(ctx context.Context, stationID uint) ([]models.Staff, error) {
	var staff []models.Staff
	if err := r.db.WithContext(ctx).Preload("Station").Where("station_id = ?", stationID).Find(&staff).Error; err != nil {
		return nil, fmt.Errorf("failed to find staff: %w", err)
	}
	return staff, nil
}

// Update updates a staff
func (r *staffRepository) Update(ctx context.Context, staff *models.Staff) error {
	if err := r.db.WithContext(ctx).Save(staff).Error; err != nil {
		return fmt.Errorf("failed to update staff: %w", err)
	}
	return nil
}

// Delete removes a staff
func (r *staffRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Staff{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete staff: %w", err)
	}
	return nil
}
