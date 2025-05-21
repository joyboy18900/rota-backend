package repositories

import (
	"context"
	"rota-api/models"

	"gorm.io/gorm"
)

type StopRepository interface {
	Create(ctx context.Context, stop *models.Stop) error
	FindByID(ctx context.Context, id string) (*models.Stop, error)
	FindAll(ctx context.Context) ([]*models.Stop, error)
	Update(ctx context.Context, stop *models.Stop) error
	Delete(ctx context.Context, id string) error
}

type stopRepository struct {
	db *gorm.DB
}

func NewStopRepository(db *gorm.DB) StopRepository {
	return &stopRepository{db: db}
}

func (r *stopRepository) Create(ctx context.Context, stop *models.Stop) error {
	return r.db.WithContext(ctx).Create(stop).Error
}

func (r *stopRepository) FindByID(ctx context.Context, id string) (*models.Stop, error) {
	var stop models.Stop
	err := r.db.WithContext(ctx).
		Preload("Routes").
		First(&stop, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &stop, nil
}

func (r *stopRepository) FindAll(ctx context.Context) ([]*models.Stop, error) {
	var stops []*models.Stop
	err := r.db.WithContext(ctx).
		Preload("Routes").
		Find(&stops).Error
	return stops, err
}

func (r *stopRepository) Update(ctx context.Context, stop *models.Stop) error {
	return r.db.WithContext(ctx).Save(stop).Error
}

func (r *stopRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Stop{}, "id = ?", id).Error
}
