package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/nemopss/subscription-service/internal/models"
	"github.com/nemopss/subscription-service/pkg/logger"
)

var DB *gorm.DB

func InitDB() error {
	logger.Log.Info("Initialising database")
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.WithError(err).Error("Error opening database")
		return err
	}

	DB = db
	DB.AutoMigrate(&models.Subscription{})
	logger.Log.Info("Database initialised")
	return nil
}
