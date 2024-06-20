package controllers

import (
	"blogbackend/database"
	"blogbackend/models"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetCategories(c *fiber.Ctx) error {
	var getcategories []models.Category
	database.DB.Find(&getcategories)

	return c.JSON(fiber.Map{
		"data": getcategories,
	})
}

func CreateCategory(c *fiber.Ctx) error {
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		log.Println("Не удалось получить тело запроса")
	}
	if err := database.DB.Create(&category).Error; err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Не удалось создать категорию",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Категория создана",
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	post := models.Post{
		Id: uint(id),
	}
	if err := c.BodyParser(&post); err != nil {
		fmt.Println("Не удалось получить тело запроса")
	}

	database.DB.Model(&post).Updates(post)
	return c.JSON(fiber.Map{
		"message": "Объявление успешно изменено",
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	category := models.Category{
		Id: uint(id),
	}
	deleteQuery := database.DB.Delete(&category)
	if errors.Is(deleteQuery.Error, gorm.ErrRecordNotFound) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Категория не найдена",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Категория успешно удалена!",
	})
}

func GetByCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	sortParam := c.Query("sort", "created_at_desc")

	var posts []models.Post
	db := database.DB.Where("category_id = ? AND is_submit = ? AND is_archive = ?", id, true, false).Preload("Category").Preload("User")

	switch sortParam {
	case "created_at_asc":
		db = db.Order("created_at ASC")
	case "created_at_desc":
		db = db.Order("created_at DESC")
	case "estimated_price_asc":
		db = db.Order("estimated_price ASC")
	case "estimated_price_desc":
		db = db.Order("estimated_price DESC")
	default:
		db = db.Order("created_at DESC")
	}

	if err := db.Find(&posts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Ошибка при выполнении поиска",
			"error":   err.Error(),
		})
	}

	log.Println(len(posts))
	return c.JSON(posts)
}
