package main

import (
	"blogbackend/database"
	"blogbackend/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {

	database.Connect()

	if err := initConfig(); err != nil {
		log.Fatal("error initializing configs")
	}
	port := viper.GetString("port")

	app := fiber.New()
	routes.Setup(app)
	app.Listen(":" + port)

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
