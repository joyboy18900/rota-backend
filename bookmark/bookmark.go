package bookmark

import (
	"rota-api/database"

	"github.com/gofiber/fiber/v2"
)

// GetAllBookmarks godoc
// @Summary Get all bookmarks
// @Description Get all bookmarks
// @Tags bookmarks
// @Accept json
// @Produce json
// @Success 200 {array} database.Bookmark
// @Failure 500 {object} fiber.Map
// @Router /api/bookmark [get]
func GetAllBookmarks(c *fiber.Ctx) error {
	result, err := database.GetAllBookmarks()
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
			"data":    nil,
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"success": true,
		"message": "",
		"data":    result,
	})
}

// SaveBookmark godoc
//
//	@Summary		Save a new bookmark
//	@Description	Save a new bookmark
//	@Tags			bookmarks
//	@Accept			json
//	@Produce		json
//	@Param			bookmark	body		database.Bookmark	true	"Bookmark"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Map
//	@Router			/api/bookmark [post]
func SaveBookmark(c *fiber.Ctx) error {
	newBookmark := new(database.Bookmark)

	err := c.BodyParser(newBookmark)
	if err != nil {
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return err
	}

	err = database.CreateBookmark(*newBookmark)
	if err != nil {
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return err
	}

	c.Status(200).JSON(&fiber.Map{
		"success": true,
		"message": "Bookmark saved successfully",
		"data":    newBookmark,
	})
	return nil
}
