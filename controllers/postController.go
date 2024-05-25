package controllers

import (
	"blogbackend/database"
	"blogbackend/models"
	"blogbackend/util"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreatePost(c *fiber.Ctx) error {
	var post models.Post
	if err := c.BodyParser(&post); err != nil {
		fmt.Println("Unable to parse body")
	}
	if err := database.DB.Create(&post).Error; err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid payload",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Your post created",
	})
}

func AllPost(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit := 10
	offset := (page - 1) * limit
	var total int64
	var getpost []models.Post
	database.DB.Preload("User").Offset(offset).Limit(limit).Find(&getpost)
	database.DB.Model(&models.Post{}).Count(&total)

	return c.JSON(fiber.Map{
		"data": getpost,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": math.Ceil(float64(int(total) / limit)),
		},
	})
}

func DetailPost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var post models.Post
	database.DB.Where("id=?", id).Preload("User").First(&post)
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
		fmt.Println("Unable to parse body")
	}

	database.DB.Model(&post).Updates(post)
	return c.JSON(fiber.Map{
		"message": "You successfully update",
	})
}

func UniquePost(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	id, _ := util.ParseJwt(cookie)
	var post []models.Post
	database.DB.Model(&post).Where("user_id=?", id).Preload("User").Find(&post)

	return c.JSON(post)
}

func DeletePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	post := models.Post{
		Id: uint(id),
	}
	deleteQuery := database.DB.Delete(&post)
	if errors.Is(deleteQuery.Error, gorm.ErrRecordNotFound) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Post not found!",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Post delete successfully!",
	})

}
