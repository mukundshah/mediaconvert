package db

import (
	"fmt"
	"log"

	"github.com/mukund/mediaconvert/internal/config"
	"github.com/mukund/mediaconvert/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to database")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	log.Println("Running migrations...")
	return db.AutoMigrate(&models.User{}, &models.File{}, &models.Job{})
}
