package db

import (
	"log"

	"github.com/Arpitmovers/reviewservice/internal/repository/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Hotel{},
		&models.Reviewer{},
		&models.Review{},
		&models.ProviderScore{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
