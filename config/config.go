package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Environment     string
	ServerPort      string
	JWTSecret       string
	TokenExpiration int // in hours

	// PostgreSQL configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	Redis struct {
		Host     string `env:"REDIS_HOST" envDefault:"localhost"`
		Port     string `env:"REDIS_PORT" envDefault:"6379"`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
		DB       int    `env:"REDIS_DB" envDefault:"0"`
	}

	JWT struct {
		Secret          string        `env:"JWT_SECRET" envDefault:"your-secret-key"`
		AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TOKEN_TTL" envDefault:"24h"`
		RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TOKEN_TTL" envDefault:"168h"` // 7 days
	}
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// Load environment specific .env file if exists
	env := getEnv("APP_ENV", "development")
	if env != "development" {
		if err := godotenv.Load(fmt.Sprintf(".env.%s", env)); err != nil {
			return nil, fmt.Errorf("error loading .env.%s file: %v", env, err)
		}
	}

	log.Printf("Running in %s environment", env)

	// After loading all environment files, create config
	cfg := &Config{
		Environment:     env,
		ServerPort:      getEnv("SERVER_PORT", "3000"),
		JWTSecret:       getEnv("JWT_SECRET", "your-jwt-secret-key"),
		TokenExpiration: getEnvAsInt("TOKEN_EXPIRATION", 24),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "rota_dev"),
		DBPassword: getEnv("DB_PASSWORD", "rota_dev"),
		DBName:     getEnv("DB_NAME", "rota_dev"),
	}

	// Load Redis configuration
	cfg.Redis.Host = getEnv("REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnv("REDIS_PORT", "6379")
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	cfg.Redis.DB = getEnvAsInt("REDIS_DB", 0)

	// Load JWT configuration
	cfg.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key")
	cfg.JWT.AccessTokenTTL, _ = time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_TTL", "24h"))
	cfg.JWT.RefreshTokenTTL, _ = time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_TTL", "168h"))

	return cfg, nil
}

// Helper function to get an environment variable or a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get an environment variable as int or a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
