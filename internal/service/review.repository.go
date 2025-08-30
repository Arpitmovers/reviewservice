package services

import (
	"fmt"

	dto "github.com/Arpitmovers/reviewservice/internal/handlers/dto"

	"github.com/Arpitmovers/reviewservice/internal/repository/db"
	"gorm.io/gorm"
)

// home/arpit/code/personal/reviewService/internal/repository/review.repository.go

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// func NewReviewService(db *gorm.DB) *ReviewService {
// 	return &ReviewService{DB: db}
// }

func SaveReview(msg *dto.Review) error {
	// Map DTO → Models
	hotel := mapToHotel(msg)
	reviewer := mapToReviewer(msg)
	review := mapToReview(msg)
	providerScore := mapToProviderScore(msg)

	// Transaction
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := db.UpsertHotel(tx, hotel); err != nil {
			fmt.Println("error in  UpsertHotel ", err)
			return err
		}
		if err := db.UpsertReviewer(tx, reviewer); err != nil {
			fmt.Println("error in  UpsertReviewer ", err)
			return err
		}
		if err := db.UpsertReview(tx, review); err != nil {
			fmt.Println("error in  UpsertReview ", err)
			return err
		}
		if err := db.UpsertProviderScore(tx, providerScore); err != nil {
			fmt.Println("error in  UpsertProviderScore ", err)
			return err
		}
		return nil
	})
}
