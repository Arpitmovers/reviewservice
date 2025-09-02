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

func (reviewRepo *ReviewService) SaveReview(msg *dto.HotelReviewDTO) error {

	hotel, provider, reviewer, review, providedSummary := MapHotelReviewDTOToModels(msg)

	return reviewRepo.Repo.Db.Transaction(func(tx *gorm.DB) error {
		if err := reviewRepo.Repo.InsertHotel(hotel); err != nil {
			logger.Logger.Error("error in  UpsertHotel ", zap.Error(err))
			return err
		}

		if err := reviewRepo.Repo.InsertProvider(provider); err != nil {
			logger.Logger.Error("error in  InsertProvider ", zap.Error(err))
			return err
		}
		var reviwerId int
		var err error

		if reviwerId, err = reviewRepo.Repo.InsertReviewer(reviewer); err != nil {

			logger.Logger.Error("error in  InsertReviewer ", zap.Error(err))
			return err
		}
		if reviwerId > 0 {
			logger.Logger.Info("reviewer with id exists", zap.Int("reviwerId", reviwerId))
		}

		review.ReviewerID = reviwerId
		if err := reviewRepo.Repo.InsertReview(review); err != nil {
			logger.Logger.Error("error in  InsertReview ", zap.Error(err))
			return err
		}

		if err := reviewRepo.Repo.UpsertProviderSummary(providedSummary); err != nil {
			logger.Logger.Error("error in  UpsertProviderSummary ", zap.Error(err))
			return err
		}
		return nil
	})
}
