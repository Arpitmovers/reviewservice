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

type HotelReviewRepository interface {
	InsertHotel(hotel *Hotel) error
	InsertProvider(provider *Provider) error
	InsertReviewer(reviewer *Reviewer) (int, error)
	InsertReview(review *Review) error
	UpsertProviderSummary(summary *ProviderSummary) error
}

func (r *ReviewRepository) InsertHotel(hotel *Hotel) error {
	return r.Db.FirstOrCreate(hotel, Hotel{HotelID: hotel.HotelID}).Error
}

func (r *ReviewRepository) InsertProvider(provider *Provider) error {
	return r.Db.FirstOrCreate(provider, Provider{ProviderID: provider.ProviderID}).Error
}

// return revierId,error
func (r *ReviewRepository) InsertReviewer(reviewer *Reviewer) (int, error) {
	err := r.Db.Create(reviewer).Error
	return reviewer.ReviewerID, err
}

// insert only if <reviewer_id + review_date > doenst exist
func (r *ReviewRepository) InsertReview(review *Review) error {

	return r.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "reviewer_id"}, {Name: "review_date"}},
		DoNothing: true, // skip insert if exists
	}).Create(review).Error
}

func (r *ReviewRepository) UpsertProviderSummary(summary *ProviderSummary) error {
	return r.Db.Clauses(
		clause.OnConflict{
			UpdateAll: true, // update all fields on conflict
		},
	).Create(summary).Error
}
