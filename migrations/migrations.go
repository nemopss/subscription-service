package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/nemopss/subscription-service/internal/models"
)

func Migrate20250721(db *gorm.DB) *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20250721210000",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Subscription{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("subscriptions")
		},
	}
}
