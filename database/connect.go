package database

import (
	"blogbackend/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	if err := initConfig(); err != nil {
		log.Fatal("error initializing configs")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}

	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	username := viper.GetString("db.username")
	dbname := viper.GetString("db.dbname")
	sslmode := viper.GetString("db.sslmode")
	password := os.Getenv("DB_PASSWORD")

	fmt.Print(host)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, username, password, dbname, port, sslmode)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database")
	} else {
		log.Println("Connect successfully")
	}

	DB = database

	database.AutoMigrate(
		&models.User{},
		&models.Post{},
	)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
