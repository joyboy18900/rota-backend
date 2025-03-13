package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"rota-api/models"
)

// InitDatabase initializes the database connection
func InitDatabase(cfg *Config) (*gorm.DB, error) {
	// Ensure directory exists
	dbDir := filepath.Dir(cfg.DBPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err := os.MkdirAll(dbDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Set up database connection
	gormConfig := &gorm.Config{}

	// Add logger for development environment
	if cfg.Environment == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.OAuthToken{},
		&models.Station{},
		&models.Route{},
		&models.Schedule{},
		&models.Favorite{},
		&models.Staff{},
		&models.Vehicle{},
		&models.ScheduleLog{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}
