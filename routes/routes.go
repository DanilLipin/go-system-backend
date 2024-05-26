package routes

import (
	"blogbackend/controllers"
	"blogbackend/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Get("/api", controllers.AllPost)
	app.Get("/api/:id", controllers.DetailPost)
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)

	app.Use(middleware.IsAuthenticate)
	app.Post("/api/post", controllers.CreatePost)
	app.Put("/api/updatepost/:id", controllers.UpdatePost)
	app.Get("/api/personalposts", controllers.PersonalPosts)
	app.Delete("/api/deletepost/:id", controllers.DeletePost)
	app.Post("/api/upload-image", controllers.UploadPhoto)
	app.Static("/api/uploads", "./uploads")

}
