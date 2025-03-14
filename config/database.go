package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"rota-api/models"
)

// InitDatabase initializes the database connection
func InitDatabase(cfg *Config) (*gorm.DB, error) {
	// Set up database connection
	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// Add logger for development environment
	if cfg.Environment == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	// PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
	}
	log.Println("Connected to PostgreSQL database")

	// Auto-migrate models in development only
	if cfg.Environment == "development" {
		log.Println("Running auto-migrations for development environment...")
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
	} else {
		log.Println("Skipping auto-migrations in production environment")
	}

	log.Println("Database initialized successfully")
	return db, nil
}
