package repositories

import (
	"context"

	"rota-api/models"

	"gorm.io/gorm"
)

// ScheduleLogRepository defines the interface for schedule log-related database operations
type ScheduleLogRepository interface {
	FindByID(ctx context.Context, id uint) (*models.ScheduleLog, error)
	FindAll(ctx context.Context) ([]*models.ScheduleLog, error)
	Create(ctx context.Context, scheduleLog *models.ScheduleLog) error
	Update(ctx context.Context, scheduleLog *models.ScheduleLog) error
	Delete(ctx context.Context, id uint) error
}

// scheduleLogRepository implements ScheduleLogRepository
type scheduleLogRepository struct {
	db *gorm.DB
}

// NewScheduleLogRepository creates a new schedule log repository
func NewScheduleLogRepository(db *gorm.DB) ScheduleLogRepository {
	return &scheduleLogRepository{db}
}

func (r *scheduleLogRepository) FindByID(ctx context.Context, id uint) (*models.ScheduleLog, error) {
	var scheduleLog models.ScheduleLog
	if err := r.db.WithContext(ctx).
		Preload("Schedule").
		Preload("Schedule.Route").
		Preload("Schedule.Route.StartStation").
		Preload("Schedule.Route.EndStation").
		Preload("Schedule.Vehicle").
		Preload("Schedule.Station").
		Preload("Staff").
		Preload("Staff.Station").
		First(&scheduleLog, id).Error; err != nil {
		return nil, err
	}
	return &scheduleLog, nil
}

func (r *scheduleLogRepository) FindAll(ctx context.Context) ([]*models.ScheduleLog, error) {
	var scheduleLogs []*models.ScheduleLog
	if err := r.db.WithContext(ctx).
		Preload("Schedule").
		Preload("Schedule.Route").
		Preload("Schedule.Route.StartStation").
		Preload("Schedule.Route.EndStation").
		Preload("Schedule.Vehicle").
		Preload("Schedule.Station").
		Preload("Staff").
		Preload("Staff.Station").
		Find(&scheduleLogs).Error; err != nil {
		return nil, err
	}
	return scheduleLogs, nil
}

func (r *scheduleLogRepository) Create(ctx context.Context, scheduleLog *models.ScheduleLog) error {
	// ใช้คำสั่ง SQL โดยตรงเพื่อหลีกเลี่ยงปัญหาคอลัมน์ updated_at
	sql := `INSERT INTO schedule_logs (schedule_id, staff_id, change_description) VALUES (?, ?, ?)`
	return r.db.WithContext(ctx).Exec(sql, scheduleLog.ScheduleID, scheduleLog.StaffID, scheduleLog.ChangeDescription).Error
}

func (r *scheduleLogRepository) Update(ctx context.Context, scheduleLog *models.ScheduleLog) error {
	// ใช้คำสั่ง SQL โดยตรงเพื่อหลีกเลี่ยงปัญหาคอลัมน์ updated_at
	sql := `UPDATE schedule_logs SET schedule_id = ?, staff_id = ?, change_description = ? WHERE id = ?`
	return r.db.WithContext(ctx).Exec(sql, scheduleLog.ScheduleID, scheduleLog.StaffID, scheduleLog.ChangeDescription, scheduleLog.ID).Error
}

func (r *scheduleLogRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ScheduleLog{}, id).Error
}
