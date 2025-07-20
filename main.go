package main

import (
	"context"
	"fmt"
	"log"
	"rota-api/config"
	handler "rota-api/handlers"
	"rota-api/models"
	"rota-api/repositories"
	"rota-api/routes"
	"rota-api/services"
	"rota-api/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := utils.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis (optional)
	var redisRepo repositories.RedisRepository
	if cfg.Redis.Host != "" {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})

		// Test Redis connection with retry
		ctx := context.Background()
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			if err := redisClient.Ping(ctx).Err(); err != nil {
				log.Printf("Warning: Redis connection attempt %d failed: %v", i+1, err)
				if i == maxRetries-1 {
					log.Printf("Warning: Redis connection failed after %d attempts, proceeding without Redis", maxRetries)
					redisClient.Close()
					break
				}
				time.Sleep(time.Second * 2)
				continue
			}
			redisRepo = repositories.NewRedisRepository(redisClient)
			log.Println("Connected to Redis successfully")
			break
		}
	}

	// Run migrations
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	if err := utils.RunMigrations(dbURL); err != nil {
		log.Printf("Migration error: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	routeRepo := repositories.NewRouteRepository(db)
	stationRepo := repositories.NewStationRepository(db)
	favoriteRepo := repositories.NewFavoriteRepository(db)
	vehicleRepo := repositories.NewVehicleRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	scheduleLogRepo := repositories.NewScheduleLogRepository(db)
	staffRepo := repositories.NewStaffRepository(db)

	// Initialize services
	authConfig := services.AuthConfig{
		TokenConfig: models.TokenConfig{
			Secret:     cfg.JWT.Secret,
			ExpiryTime: cfg.JWT.AccessTokenTTL,
		},
		RedisConfig: &repositories.RedisConfig{
			Host:     cfg.Redis.Host,
			Port:     cfg.Redis.Port,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		},
	}

	googleOAuthService := services.NewGoogleOAuthService()
	
	authService := services.NewAuthService(
		userRepo,
		redisRepo,
		googleOAuthService,
		authConfig,
	)

	// Initialize services
	userService := services.NewUserService(userRepo)
	routeService := services.NewRouteService(routeRepo)
	stationService := services.NewStationService(stationRepo)
	favoriteService := services.NewFavoriteService(favoriteRepo)
	vehicleService := services.NewVehicleService(vehicleRepo)
	scheduleService := services.NewScheduleService(scheduleRepo)
	scheduleLogService := services.NewScheduleLogService(scheduleLogRepo)
	staffService := services.NewStaffService(staffRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	googleOAuthHandler := handler.NewGoogleOAuthHandler(authService, googleOAuthService)
	userHandler := handler.NewUserHandler(userService, authService)
	routeHandler := handler.NewRouteHandler(routeService)
	stationHandler := handler.NewStationHandler(stationService)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)
	vehicleHandler := handler.NewVehicleHandler(vehicleService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	scheduleLogHandler := handler.NewScheduleLogHandler(scheduleLogService)
	staffHandler := handler.NewStaffHandler(staffService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Rota API",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		ExposeHeaders: "Content-Length, Authorization",
		AllowCredentials: false, // Changed to false to work with wildcard origins
		MaxAge: 86400, // 24 hours
	}))

	// Health check endpoint
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "API is healthy",
		})
	})

	// Routes
	routes.SetupAuthRoutes(app, authHandler, googleOAuthHandler, authService)
	routes.SetupUserRoutes(app, userHandler, authService)
	routes.SetupRouteRoutes(app, routeHandler, authService)
	routes.SetupStationRoutes(app, stationHandler, scheduleHandler, authService)
	routes.SetupFavoriteRoutes(app, favoriteHandler, authService, favoriteService)
	routes.SetupVehicleRoutes(app, vehicleHandler, authService)
	routes.SetupScheduleRoutes(app, scheduleHandler, authService)
	routes.SetupScheduleLogRoutes(app, scheduleLogHandler, authService)
	// Add routes for staff
	routes.SetupStaffRoutes(app, staffHandler, authService)

	// Start server
	log.Printf("Server starting on :%s", cfg.ServerPort)
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
