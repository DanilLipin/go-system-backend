package controllers

import (
	"blogbackend/database"
	"blogbackend/models"

	"github.com/gofiber/fiber/v2"
)

func SearchPosts(c *fiber.Ctx) error {
	query := c.Query("query")

	var posts []models.Post
	db := database.DB

	if query != "" {
		db = db.Where("title LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if err := db.Preload("User").Preload("Category").Find(&posts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Ошибка при выполнении поиска",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": posts,
	})
}
