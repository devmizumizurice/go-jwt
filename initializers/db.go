package initializers

import (
	"fmt"
	"os"

	"github.com/devmizumizurice/go-jwt/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetUpDB() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Tokyo",
		host, port, user, password, dbname,
	)

	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Database Connection Error")
	}
}

func SyncDB() {
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	var err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("Database Migration Error")
	}
}

func GetDB() *gorm.DB {
	return db
}
