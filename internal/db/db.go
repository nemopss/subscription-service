package db

import (
	"fmt"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/nemopss/subscription-service/migrations"
	"github.com/nemopss/subscription-service/pkg/logger"
)

var DB *gorm.DB

func InitDB() error {
	logger.Log.Info("Initialising database")
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		return fmt.Errorf("переменная окружения POSTGRES_URL не установлена")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.WithError(err).Error("Error opening database")
		return err
	}

	DB = db

	m := gormigrate.New(
		DB,
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{migrations.Migrate20250721(DB)},
	)

	if err := m.Migrate(); err != nil {
		logger.Log.WithError(err).Error("Migration failed")
		return err
	}

	logger.Log.Info("Database initialised")
	return nil
}
