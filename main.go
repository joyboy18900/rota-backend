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
	var redisRepo *repositories.RedisRepository
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

	authService := services.NewAuthService(
		userRepo,
		redisRepo,
		authConfig,
	)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Rota API",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	routes.SetupAuthRoutes(app, authHandler, authService)

	// Start server
	log.Printf("Server starting on :%s", cfg.ServerPort)
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// import (
// 	"rota-api/bookmark"
// 	"rota-api/database"

// 	"github.com/gofiber/fiber/v2"
// )

// func main() {
// 	app := fiber.New()
// 	dbErr := database.InitDatabase()

// 	if dbErr != nil {
// 		panic(dbErr)
// 	}

// 	setupRoutes(app)
// 	app.Listen(":3000")
// }

// func status(c *fiber.Ctx) error {
// 	return c.SendString("Server is running! Send your request")
// }

// func setupRoutes(app *fiber.App) {
// 	app.Get("/", status)
// 	app.Get("/bookmark", bookmark.GetAllBookmarks)
// 	app.Post("/bookmark", bookmark.SaveBookmark)
// }
