package routes

import (
	"blogbackend/controllers"
	"blogbackend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup(app *fiber.App) {

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000/", // Allow all origins
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}))

	app.Get("/api/posts", controllers.AllPost)
	app.Get("/api/post/:id", controllers.DetailPost)
	app.Post("/api/registration", controllers.Registration)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/categories", controllers.GetCategories)
	app.Get("/api/category/:id", controllers.GetByCategory)
	app.Get("/api/user/:id", controllers.DetailUser)
	app.Get("/api/user/:id/posts", controllers.GetPostsByUser)
	app.Get("/api/search-posts", controllers.SearchPosts)
	app.Static("/api/uploads", "./uploads")

	app.Use(middleware.IsAuthenticate)

	app.Get("/api/logout", controllers.Logout)
	app.Post("/api/post", controllers.CreatePost)
	app.Post("/api/upload-image", controllers.UploadPhoto)
	app.Post("/api/create-chat", controllers.CreateChatRoom)
	app.Get("/api/user/:id/chats", controllers.GetUserChats)
	app.Get("/api/chat/:id", controllers.GetChatByIDWithDetails)
	app.Post("/api/message", controllers.SendMessage)
	app.Get("/api/chat/:id/messages", controllers.GetNewMessages)
	app.Post("/api/exchange/:id/confirm", controllers.ConfirmExchange)
	app.Get("/api/exchange/:chat_room_id/:user_id/status", controllers.GetExchangeStatus)
	app.Post("/api/archive-posts/:user_post_id/:other_post_id", controllers.ArchivePosts)

	app.Get("/api/personalposts", controllers.PersonalPosts)
	app.Get("/api/user/:id/posts-similar-price", controllers.GetPostsByUserAndPriceRange)
	app.Put("/api/update-post/:id", controllers.UpdatePost)
	app.Get("/api/check-post-chat/:id", controllers.CheckPostChat)
	app.Put("/api/soft-delete/:id", controllers.SoftDeletePost)
	app.Delete("/api/delete-post/:id", controllers.DeletePost)
	app.Delete("/api/delete-account", controllers.DeleteProfile)

	app.Use(middleware.IsAdmin)
	app.Get("api/admin/non-submit-posts", controllers.GetNonSubmitPosts)
	app.Put("api/admin/post/:id", controllers.SubmitPost)
	app.Get("api/admin/logout", controllers.AdminLogout)
	app.Post("api/admin/create-admin", controllers.CreateAdmin)
	app.Post("api/admin/create-category", controllers.CreateCategory)
	app.Put("/api/admin/update-category/:id", controllers.UpdateCategory)
	app.Delete("/api/admin/delete-category/:id", controllers.DeleteCategory)

}
