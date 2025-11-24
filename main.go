package main

import (
	"log"
	"smart-queue/config"
	"smart-queue/controller"
	"smart-queue/routes"
	"smart-queue/sse"

	db "smart-queue/database"

	"github.com/gin-gonic/gin"
)

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	db.ConnectDB(cfg.DatabaseURL)
}

func main() {

	r := gin.Default()
	sse := sse.NewBroadcaster()
	bc := &controller.BusinessController{Broadcaster: sse, DB: db.DB}
	auth := &controller.AuthController{DB: db.DB, Config: cfg}

	routes.SetupRoutes(r, bc, auth)
	r.Run(":8080")
}
