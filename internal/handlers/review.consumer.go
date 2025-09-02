package handlers

import (
	"encoding/json"
	"fmt"

	dto "github.com/Arpitmovers/reviewservice/internal/handlers/dto"
	logger "github.com/Arpitmovers/reviewservice/internal/logging"
	services "github.com/Arpitmovers/reviewservice/internal/service"
	"go.uber.org/zap"
)

type ReviewConsumer struct {
	service *services.ReviewService
}

func NewReviewConsumer(service *services.ReviewService) *ReviewConsumer {
	return &ReviewConsumer{service: service}
}

func (h *ReviewConsumer) ConsumeReview() func([]byte) error {
	return func(body []byte) error {
		logger.Logger.Info("Received review message", zap.ByteString("body", body))

		var review dto.Review
		if err := json.Unmarshal(body, &review); err != nil {
			return fmt.Errorf("failed to unmarshal review JSON: %w", err)
		}

		if err := h.service.SaveReview(&review); err != nil {
			return fmt.Errorf("failed to save review: %w", err)
		}

		return nil
	}
}
