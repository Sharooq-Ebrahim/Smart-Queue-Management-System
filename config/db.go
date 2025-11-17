package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnctDatabase() {

	// dsn := fmt.Sprintf(
	// 	"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	// 	getEnv("DB_HOST", "localhost"),
	// 	getEnv("DB_USER", "postgres"),
	// 	getEnv("DB_PASSWORD", "1234"),
	// 	getEnv("DB_NAME", "mydb"),
	// 	getEnv("DB_PORT", "5432"),
	// )

	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=smart_queue port=5432 sslmode=disable"), &gorm.Config{})

	if err != nil {
		log.Println("Connection Error in Database")
	}

	DB = db

	log.Println("Database connected successfully!")

}
