package database

import (
	"final-project-dts-go/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

var (
	db *gorm.DB
	err error
)

func StartDB()  {
	err = godotenv.Load("config/.env")
	if err != nil {
		fmt.Println("Failed load file environment")
	} else {
		fmt.Println("Success read file environment")
	}

	// config := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	config := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE"))

	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to database:", err)
	}

	fmt.Println("success connecting to database")

	db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.SocialMedia{}, models.Comment{})
}

func GetDB() *gorm.DB  {
	return db
} 