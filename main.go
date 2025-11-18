package main

import (
	"log"
	"smart-queue/config"
	"smart-queue/controller"
	"smart-queue/model"
	"smart-queue/routes"
	"smart-queue/sse"

	"github.com/gin-gonic/gin"
)

func init() {

	config.ConnctDatabase()

	if err := config.DB.AutoMigrate(&model.User{}, &model.Business{}, &model.Queue{}); err != nil {
		log.Println("Failed to migrate Db")
	}

	log.Println("Database connected and migrated successfully")
}

func main() {

	// config.ConnctDatabase()

	r := gin.Default()
	sse := sse.NewBroadcaster()
	bc := &controller.BusinessController{Broadcaster: sse, DB: config.DB}
	auth := &controller.AuthController{DB: config.DB}

	routes.SetupRoutes(r, bc, auth)
	r.Run(":8080")
}
