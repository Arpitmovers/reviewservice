package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReviewRepository struct {
	Db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{Db: db}
}

func UpsertHotel(tx *gorm.DB, hotel Hotel) error {
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "hotel_id"}, {Name: "platform"}},
		DoUpdates: clause.AssignmentColumns([]string{"hotel_name"}),
	}).Create(&hotel).Error
}

func UpsertReviewer(tx *gorm.DB, reviewer Reviewer) error {
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "reviewer_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"display_name", "country_id"}),
	}).Create(&reviewer).Error
}

func UpsertReview(tx *gorm.DB, review Review) error {
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "review_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"rating", "review_title", "review_comments"}),
	}).Create(&review).Error
}

func UpsertProviderScore(tx *gorm.DB, ps ProviderScore) error {
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "hotel_id"}, {Name: "provider_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"overall_score", "review_count", "cleanliness",
			"facilities", "location", "service", "value_for_money",
		}),
	}).Create(&ps).Error
}
