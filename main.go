package main

import (
	"log"
	"rota-api/config"
	"rota-api/controllers"
	"rota-api/repositories"
	"rota-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	stationRepo := repositories.NewStationRepository(db)
	routeRepo := repositories.NewRouteRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	favoriteRepo := repositories.NewFavoriteRepository(db)
	oauthTokenRepo := repositories.NewOAuthTokenRepository(db)
	staffRepo := repositories.NewStaffRepository(db)
	vehicleRepo := repositories.NewVehicleRepository(db)
	scheduleLogRepo := repositories.NewScheduleLogRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, oauthTokenRepo, cfg)
	userService := services.NewUserService(userRepo)
	stationService := services.NewStationService(stationRepo)
	routeService := services.NewRouteService(routeRepo)
	scheduleService := services.NewScheduleService(scheduleRepo)
	favoriteService := services.NewFavoriteService(favoriteRepo)
	staffService := services.NewStaffService(staffRepo)
	vehicleService := services.NewVehicleService(vehicleRepo)
	scheduleLogService := services.NewScheduleLogService(scheduleLogRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	stationController := controllers.NewStationController(stationService)
	routeController := controllers.NewRouteController(routeService)
	scheduleController := controllers.NewScheduleController(scheduleService)
	favoriteController := controllers.NewFavoriteController(favoriteService)
	staffController := controllers.NewStaffController(staffService)
	vehicleController := controllers.NewVehicleController(vehicleService)
	scheduleLogController := controllers.NewScheduleLogController(scheduleLogService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: config.ErrorHandler,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Setup routes
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Get("/google/login", authController.GoogleLogin)
	auth.Get("/google/callback", authController.GoogleCallback)
	auth.Get("/refresh", authController.RefreshToken)

	// User routes
	users := api.Group("/users")
	users.Get("/", authMiddleware(authService), userController.GetAllUsers)
	users.Get("/:id", authMiddleware(authService), userController.GetUserByID)
	users.Put("/:id", authMiddleware(authService), userController.UpdateUser)
	users.Delete("/:id", authMiddleware(authService), userController.DeleteUser)

	// Station routes
	stations := api.Group("/stations")
	stations.Get("/", stationController.GetAllStations)
	stations.Get("/:id", stationController.GetStationByID)
	stations.Post("/", authMiddleware(authService), stationController.CreateStation)
	stations.Put("/:id", authMiddleware(authService), stationController.UpdateStation)
	stations.Delete("/:id", authMiddleware(authService), stationController.DeleteStation)

	// Route routes
	routes := api.Group("/routes")
	routes.Get("/", routeController.GetAllRoutes)
	routes.Get("/:id", routeController.GetRouteByID)
	routes.Post("/", authMiddleware(authService), routeController.CreateRoute)
	routes.Put("/:id", authMiddleware(authService), routeController.UpdateRoute)
	routes.Delete("/:id", authMiddleware(authService), routeController.DeleteRoute)

	// Schedule routes
	schedules := api.Group("/schedules")
	schedules.Get("/", scheduleController.GetAllSchedules)
	schedules.Get("/:id", scheduleController.GetScheduleByID)
	schedules.Post("/", authMiddleware(authService), scheduleController.CreateSchedule)
	schedules.Put("/:id", authMiddleware(authService), scheduleController.UpdateSchedule)
	schedules.Delete("/:id", authMiddleware(authService), scheduleController.DeleteSchedule)

	// Favorite routes
	favorites := api.Group("/favorites")
	favorites.Get("/", authMiddleware(authService), favoriteController.GetUserFavorites)
	favorites.Post("/", authMiddleware(authService), favoriteController.AddFavorite)
	favorites.Delete("/:id", authMiddleware(authService), favoriteController.RemoveFavorite)

	// Staff routes
	staff := api.Group("/staff")
	staff.Get("/", authMiddleware(authService), staffController.GetAllStaff)
	staff.Get("/:id", authMiddleware(authService), staffController.GetStaffByID)
	staff.Post("/", authMiddleware(authService), staffController.CreateStaff)
	staff.Put("/:id", authMiddleware(authService), staffController.UpdateStaff)
	staff.Delete("/:id", authMiddleware(authService), staffController.DeleteStaff)

	// Vehicle routes
	vehicles := api.Group("/vehicles")
	vehicles.Get("/", authMiddleware(authService), vehicleController.GetAllVehicles)
	vehicles.Get("/:id", authMiddleware(authService), vehicleController.GetVehicleByID)
	vehicles.Post("/", authMiddleware(authService), vehicleController.CreateVehicle)
	vehicles.Put("/:id", authMiddleware(authService), vehicleController.UpdateVehicle)
	vehicles.Delete("/:id", authMiddleware(authService), vehicleController.DeleteVehicle)

	// Schedule log routes
	scheduleLogs := api.Group("/schedule-logs")
	scheduleLogs.Get("/", authMiddleware(authService), scheduleLogController.GetAllScheduleLogs)
	scheduleLogs.Get("/:id", authMiddleware(authService), scheduleLogController.GetScheduleLogByID)
	scheduleLogs.Post("/", authMiddleware(authService), scheduleLogController.CreateScheduleLog)

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Fatal(app.Listen(":" + cfg.ServerPort))
}

func authMiddleware(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization token")
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		userID, err := authService.ValidateToken(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Set userID in context for later use
		c.Locals("userID", userID)
		return c.Next()
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
