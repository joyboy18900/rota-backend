package main

import (
	"rota-api/bookmark"
	"rota-api/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "rota-api/docs" // Import generated docs
)

// @title			Rota API
// @version		1.0
// @description	This is a sample server for a rota application.
// @host			localhost:3000
// @BasePath		/
func main() {
	app := fiber.New()
	dbErr := database.InitDatabase()

	if dbErr != nil {
		panic(dbErr)
	}

	setupRoutes(app)
	app.Listen(":3000")
}

func status(c *fiber.Ctx) error {
	return c.SendString("Server is running! Send your request")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", status)
	app.Get("/api/bookmark", bookmark.GetAllBookmarks)
	app.Post("/api/bookmark", bookmark.SaveBookmark)
	app.Get("/swagger/*", swagger.HandlerDefault)
}
