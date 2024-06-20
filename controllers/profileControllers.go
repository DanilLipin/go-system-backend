package controllers

import (
	"blogbackend/database"
	"blogbackend/models"
	"blogbackend/util"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DetailUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var user models.User
	database.DB.Where("id = ?", id).Preload("Post").Preload("Posts.Category").First(&user)

	if user.Id == 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Пользователь не найден!",
		})

	}
	return c.JSON(fiber.Map{
		"data": user,
	})
}

func GetPostsByUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var posts []models.Post
	database.DB.Where("user_id=?", id).Preload("User").Preload("Category").Find(&posts)
	database.DB.Model(&models.Post{})
	return c.JSON(fiber.Map{
		"data": posts,
	})
}

func DeleteProfile(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	id, _, _ := util.ParseJWT(cookie)
	var user []models.User
	database.DB.Model(&user).Where("user_id=?", id).First(&user)

	deleteQuery := database.DB.Delete(&user)
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "jwt",
		Value:   "",
		Expires: expired,
	})

	if errors.Is(deleteQuery.Error, gorm.ErrRecordNotFound) {

		c.Status(400)

		return c.JSON(fiber.Map{

			"message": "Не удалось удалить профиль",
		})

	}

	return c.JSON(fiber.Map{

		"message": "Профиль удален успешно",
	})
}
