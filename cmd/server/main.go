package main

import (
	"log"
	"monthly-journal/internal/config"
	"monthly-journal/internal/database"
	"monthly-journal/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	gin.SetMode(gin.DebugMode)
	r := routes.SetupRoutes(db, cfg)

	log.Printf("Starting server on port %s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
