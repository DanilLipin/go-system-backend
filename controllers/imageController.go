package controllers

import (
	"image"
	"log"
	"math/rand"
	"unicode/utf8"

	"github.com/disintegration/letteravatar"
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
		fileName = randomLetter(5) + "-" + file.Filename
		log.Println(fileName)
		if err := c.SaveFile(file, "./uploads/"+fileName); err != nil {
			return nil
		}
	}

	return c.JSON(fiber.Map{
		"url": "http://localhost:3000/api/uploads/" + fileName,
	})
}

// Генерация аватара пользователя
func ProfilePhotoGenerator(first_name string) (image.Image, error) {
	firstLetter, _ := utf8.DecodeLastRuneInString(first_name)

	img, err := letteravatar.Draw(75, firstLetter, nil)
	if err != nil {
		log.Fatal(err)
	}

	return img, nil
}
