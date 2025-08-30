package db

import (
	"github.com/Arpitmovers/reviewservice/internal/repository/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpsertHotel(tx *gorm.DB, hotel models.Hotel) error {
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "hotel_id"}, {Name: "platform"}},
		DoUpdates: clause.AssignmentColumns([]string{"hotel_name"}),
	}).Create(&hotel).Error
}

func UpsertReviewer(tx *gorm.DB, reviewer models.Reviewer) error {
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "reviewer_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"display_name", "country_id"}),
	}).Create(&reviewer).Error
}

func UpsertReview(tx *gorm.DB, review models.Review) error {
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "review_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"rating", "review_title", "review_comment"}),
	}).Create(&review).Error
}

func UpsertProviderScore(tx *gorm.DB, ps models.ProviderScore) error {
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "hotel_id"}, {Name: "provider_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"overall_score", "review_count", "cleanliness",
			"facilities", "location", "service", "value_for_money",
		}),
	}).Create(&ps).Error
}
