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
	if err := db.AutoMigrate(&models.User{}, &models.File{}, &models.Pipeline{}, &models.Job{}); err != nil {
		return err
	}

	// Add unique index for (UserID, Name) on Pipeline
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_pipelines_user_name ON pipelines(user_id, name)").Error; err != nil {
		log.Printf("Warning: Failed to create unique index on pipelines: %v", err)
	}

	return nil
}
