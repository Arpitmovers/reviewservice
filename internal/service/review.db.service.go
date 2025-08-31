package services

import (
	"fmt"

	dto "github.com/Arpitmovers/reviewservice/internal/handlers/dto"
	"github.com/Arpitmovers/reviewservice/internal/repository/models"
	"gorm.io/gorm"
)

type ReviewService struct {
	Repo *models.ReviewRepository
}

func NewReviewService(repo *models.ReviewRepository) *ReviewService {
	return &ReviewService{Repo: repo}
}

func (reviewRepo *ReviewService) SaveReview(msg *dto.Review) error {

	hotel := mapToHotel(msg)
	reviewer := mapToReviewer(msg)
	review := mapToReview(msg)
	providerScore := mapToProviderScore(msg)

	return reviewRepo.Repo.Db.Transaction(func(tx *gorm.DB) error {
		if err := models.UpsertHotel(tx, hotel); err != nil {
			fmt.Println("error in  UpsertHotel ", err)
			return err
		}
		if err := models.UpsertReviewer(tx, reviewer); err != nil {
			fmt.Println("error in  UpsertReviewer ", err)
			return err
		}

		if err := models.UpsertProviderScore(tx, providerScore); err != nil {
			fmt.Println("error in  UpsertProviderScore ", err)
			return err
		}
		if err := models.UpsertReview(tx, review); err != nil {
			fmt.Println("error in  UpsertReview ", err)
			return err
		}
		return nil
	})
}
