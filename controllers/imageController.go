package controllers

import (
	"log"
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func randomLetter(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func UploadPhoto(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["image"]
	fileName := ""

	for _, file := range files {
		fileName = randomLetter(15) + "-" + file.Filename
		log.Println(fileName)
		if err := c.SaveFile(file, "./uploads/items-photo/"+fileName); err != nil {
			return nil
		}
	}

	return c.JSON(fiber.Map{
		"url": "http://localhost:8080/api/uploads/items-photo/" + fileName,
	})
}
