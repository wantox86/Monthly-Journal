package database

import (
	"fmt"
	"log"
	"monthly-journal/internal/config"
	"monthly-journal/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	log.Printf("Connected to database: %s", cfg.DBName)

	// Auto migration
	if err := db.AutoMigrate(&models.Expense{}); err != nil {
		log.Printf("Failed to migrate tables: %v", err)
		return nil, err
	}

	log.Println("Database migration completed")

	return db, nil
}
