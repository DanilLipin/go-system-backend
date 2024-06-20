package controllers

import (
	"blogbackend/database"
	"blogbackend/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetPostsByUserAndPriceRange(c *fiber.Ctx) error {
	userId := c.Params("id")

	priceParam := c.Query("estimatedPrice")

	price, err := strconv.Atoi(priceParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный параметр оценочной стоимости",
		})
	}

	var posts []models.Post
	db := database.DB

	lowPrice := float32(price) * 0.8
	highPrice := float32(price) * 1.2

	if err := db.Where("user_id = ? AND is_submit = ? AND is_archive = ? AND estimated_price BETWEEN ? AND ?", userId, true, false, lowPrice, highPrice).Find(&posts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка получения объявлений",
		})
	}

	return c.JSON(fiber.Map{
		"data": posts,
	})
}

type ConfirmExchangeRequest struct {
	ChatRoomId uint `json:"chat_room_id"`
}

func ConfirmExchange(c *fiber.Ctx) error {
	userIdParam := c.Params("id")
	userId, err := strconv.ParseUint(userIdParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный id параметр",
		})
	}

	var request ConfirmExchangeRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Не удалось получить тело запроса",
		})
	}

	db := database.DB
	var exchange models.Exchange

	if err := db.Where("chat_room_id = ?", request.ChatRoomId).First(&exchange).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			exchange = models.Exchange{
				ChatRoomID:    request.ChatRoomId,
				Confirmations: models.Uint64Array{userId},
				IsCompleted:   false,
			}
			db.Create(&exchange)
			return c.JSON(fiber.Map{
				"status": "Обмен подтвержден, ожидается подтверждение от другого пользователя",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Ошибка подтверждения обмена",
		})
	}

	alreadyConfirmed := false
	for _, id := range exchange.Confirmations {
		if id == userId {
			alreadyConfirmed = true
			break
		}
	}

	if !alreadyConfirmed {
		exchange.Confirmations = append(exchange.Confirmations, userId)
		if len(exchange.Confirmations) == 2 {
			exchange.IsCompleted = true
		}
		db.Save(&exchange)
	}

	if exchange.IsCompleted {
		return c.JSON(fiber.Map{
			"status": "Обмен завершен",
		})
	}

	return c.JSON(fiber.Map{
		"status": "Обмен подтвержден, ожидается подтверждение от другого пользователя",
	})
}

func GetExchangeStatus(c *fiber.Ctx) error {
	chatRoomIdParam := c.Params("chat_room_id")
	chatRoomId, err := strconv.ParseUint(chatRoomIdParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный id чата",
		})
	}

	userIdParam := c.Params("user_id")
	userId, err := strconv.ParseUint(userIdParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный id пользователя",
		})
	}

	db := database.DB
	var exchange models.Exchange

	if err := db.Where("chat_room_id = ?", chatRoomId).First(&exchange).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"status": "В процессе",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Ошибка получения статуса обмена",
		})
	}

	if exchange.IsCompleted {
		return c.JSON(fiber.Map{
			"status": "Обмен завершен",
		})
	}

	userConfirmed := false
	otherConfirmed := false
	for _, confirmId := range exchange.Confirmations {
		if confirmId == userId {
			userConfirmed = true
		} else {
			otherConfirmed = true
		}
	}

	if userConfirmed {
		return c.JSON(fiber.Map{
			"status": "Обмен подтвержден, ожидается подтверждение от другого пользователя",
		})
	}

	if otherConfirmed {
		return c.JSON(fiber.Map{
			"status": "Другой пользователь подтвердил обмен, ожидается подтверждение от вас",
		})
	}

	return c.JSON(fiber.Map{
		"status": "В процессе",
	})
}
