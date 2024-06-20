package controllers

import (
	"blogbackend/database"
	"blogbackend/models"
	"blogbackend/util"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AdminLogout(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  expired,
		HTTPOnly: true,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Вы вышли из системы"})
}

func GetNonSubmitPosts(c *fiber.Ctx) error {
	var getpost []models.Post
	database.DB.Preload("User").Preload("Category").Where("is_submit = ?", false).Find(&getpost)

	return c.JSON(fiber.Map{
		"data": getpost,
	})
}

func SubmitPost(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB

	if err := db.Model(&models.Post{}).Where("id = ?", id).Update("is_submit", true).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось подтвердить публикацию объявления",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Объявление подтверждено успешно",
	})
}

func CreateAdmin(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Не удалось получить тело запроса")
	}

	if len(data["password"].(string)) <= 6 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Пароль должен быть больше 6 символов",
		})
	}

	if !validateEmail(strings.TrimSpace(data["email"].(string))) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Неправильный формат Email",
		})
	}

	//Check email in db
	database.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.Id != 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Аккаунт с этим Email уже зарегистрирован",
		})
	}

	user := models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Phone:     data["phone"].(string),
		Email:     strings.TrimSpace(data["email"].(string)),
		IsAdmin:   true,
	}

	avatarPath := fmt.Sprintf("./uploads/profile-pictures/%s.png", strings.TrimSpace(data["email"].(string)))
	if err := util.CreateAvatar(user.FirstName, user.LastName, avatarPath); err != nil {
		log.Println("Ошибка создания аватара:", err)
	}

	user.ProfilePhoto = fmt.Sprintf("http://localhost:8080/api/uploads/profile-pictures/%s.png", strings.TrimSpace(data["email"].(string)))

	user.SetPassword(data["password"].(string))
	err := database.DB.Create(&user)
	if err != nil {
		log.Println(err)
	}
	c.Status(200)
	return c.JSON(fiber.Map{
		"user":    user,
		"message": "Аккаунт администратора успешно зарегистрирован",
	})
}
