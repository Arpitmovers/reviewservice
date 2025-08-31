package services

import (
	dto "github.com/Arpitmovers/reviewservice/internal/handlers/dto"
	logger "github.com/Arpitmovers/reviewservice/internal/logging"
	"github.com/Arpitmovers/reviewservice/internal/repository/models"
	"go.uber.org/zap"
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
			logger.Logger.Error("error in  UpsertHotel ", zap.Error(err))
			return err
		}
		if err := models.UpsertReviewer(tx, reviewer); err != nil {

			logger.Logger.Error("error in  UpsertReviewer ", zap.Error(err))
			return err
		}

		if err := models.UpsertProviderScore(tx, providerScore); err != nil {
			logger.Logger.Error("error in  UpsertProviderScore ", zap.Error(err))
			return err
		}
		if err := models.UpsertReview(tx, review); err != nil {
			logger.Logger.Error("error in  UpsertReview ", zap.Error(err))
			return err
		}
		return nil
	})
}
