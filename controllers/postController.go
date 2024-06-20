package controllers

import (
	"blogbackend/database"
	"blogbackend/models"
	"blogbackend/util"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreatePost(c *fiber.Ctx) error {
	var post models.Post
	if err := c.BodyParser(&post); err != nil {
		fmt.Println("Не удалось получить тело запроса")
	}

	post.IsSubmit = false
	post.IsArchive = false

	imagesJSON, err := json.Marshal(post.Images)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Ошибка сериализации изображений",
		})
	}
	post.Images = string(imagesJSON)

	if err := database.DB.Create(&post).Error; err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Не получилось создать объявление",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Ваше объявление создано",
	})
}

func AllPost(c *fiber.Ctx) error {
	sortParam := c.Query("sort", "created_at_desc")

	var getpost []models.Post
	db := database.DB.Preload("User").Preload("Category").Where("is_submit = ?", true)

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

	if err := db.Find(&getpost).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Ошибка при выполнении поиска",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": getpost,
	})
}

func DetailPost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var post models.Post
	database.DB.Where("id=? AND is_submit = ?", id, true).Preload("User").Preload("Category").First(&post)

	if post.Id == 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Объявление не найдено",
		})
	}

	return c.JSON(fiber.Map{
		"data": post,
	})
}

func UpdatePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	post := models.Post{
		Id: uint(id),
	}
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Не удалось получить тело запроса",
		})
	}

	post.IsSubmit = false
	fmt.Println(post.IsSubmit)

	if err := database.DB.Model(&post).Updates(post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Ошибка обновления объявления",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Объявление успешно изменено",
	})
}

func ArchivePosts(c *fiber.Ctx) error {
	userPostID := c.Params("user_post_id")
	otherPostID := c.Params("other_post_id")

	db := database.DB

	if userPostID == "" || otherPostID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ошибка получения ID объявлений",
		})
	}

	if err := db.Model(&models.Post{}).Where("id = ?", userPostID).Update("is_archive", true).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка архивации объявления",
		})
	}

	if err := db.Model(&models.Post{}).Where("id = ?", otherPostID).Update("is_archive", true).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка архивации объявления",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Объявления успешно архивированы",
	})
}

func PersonalPosts(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	id, _, _ := util.ParseJWT(cookie)
	var post []models.Post
	database.DB.Model(&post).Where("user_id=?", id).Preload("User").Find(&post)

	return c.JSON(post)
}

func DeletePost(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Params("id"))

	post := models.Post{}

	database.DB.Where("id=?", id).Preload("User").First(&post)

	if post.Id == 0 {

		c.Status(400)

		return c.JSON(fiber.Map{

			"message": "Объявление не найдено",
		})

	}

	deleteQuery := database.DB.Delete(&post)

	if errors.Is(deleteQuery.Error, gorm.ErrRecordNotFound) {

		c.Status(400)

		return c.JSON(fiber.Map{

			"message": "Объявление не найдено",
		})

	}

	return c.JSON(fiber.Map{

		"message": "Объявление успешно удалено",
	})

}

func SoftDeletePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	post := models.Post{}

	db := database.DB

	return db.Model(&post).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
