package database

import (
	"log"
	"smart-queue/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(db_url string) *gorm.DB {
	if db_url == "" {
		log.Fatalf("Database URL is empty")
	}

	db, err := gorm.Open(postgres.Open(db_url), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to Connect  to database %v", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf(" Failed to get database object: %v", err)

	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	err = db.AutoMigrate(&model.User{}, &model.Business{}, &model.Queue{})
	if err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	DB = db
	return DB
}
