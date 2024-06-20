package controllers

import (
	"blogbackend/database"
	"blogbackend/models"
	"blogbackend/util"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateChatRoom(c *fiber.Ctx) error {

	var data struct {
		Post1ID uint `json:"post1_id"`
		Post2ID uint `json:"post2_id"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Не удалось получить тело запроса",
		})
	}

	var post1, post2 models.Post
	if err := database.DB.First(&post1, data.Post1ID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Объявление не найдено",
		})
	}
	if err := database.DB.First(&post2, data.Post2ID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Объявление не найдено",
		})
	}

	chat := models.ChatRoom{
		Posts: []*models.Post{&post1, &post2},
	}

	if err := database.DB.Create(&chat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Ошибка создания чата",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Чат успешно создан",
		"data":    chat,
	})
}

func GetUserChats(c *fiber.Ctx) error {
	userID := c.Params("id")
	chats, err := util.GetChatsByUserIDWithPosts(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка получения чатов пользователя",
		})
	}
	return c.JSON(fiber.Map{
		"chats": chats,
	})
}

func GetChatByIDWithDetails(c *fiber.Ctx) error {
	chatID := c.Params("id")

	var chat models.ChatRoom
	result := database.DB.Preload("Posts").Preload("Posts.User").Preload("Messages").Preload("Messages.Sender").Where("id = ?", chatID).First(&chat)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Ошибка получения данных чата",
			"error":   result.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"chat": chat,
	})
}

func SendMessage(c *fiber.Ctx) error {
	var message models.Message

	if err := c.BodyParser(&message); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ошибка при обработке запроса",
		})
	}

	if err := database.DB.Create(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при сохранении сообщения",
		})
	}

	return c.JSON(message)
}

func GetNewMessages(c *fiber.Ctx) error {
	chatRoomID := c.Params("id")

	var messages []models.Message
	if err := database.DB.Preload("Sender").Where("chat_room_id = ?", chatRoomID).Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при получении сообщений",
		})
	}

	return c.JSON(fiber.Map{
		"messages": messages,
	})
}

func CheckPostChat(c *fiber.Ctx) error {
	postIdStr := c.Params("id")
	postId, err := strconv.ParseUint(postIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ошибка ID объявления"})
	}

	belongs, err := util.PostBelongsToChat(uint(postId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"belongs_to_chat": belongs})
}
