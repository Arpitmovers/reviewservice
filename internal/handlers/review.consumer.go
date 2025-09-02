package handlers

import (
	"encoding/json"

	dto "github.com/Arpitmovers/reviewservice/internal/handlers/dto"
	logger "github.com/Arpitmovers/reviewservice/internal/logging"
	services "github.com/Arpitmovers/reviewservice/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type ReviewConsumer struct {
	service *services.ReviewService
}

func NewReviewConsumer(service *services.ReviewService) *ReviewConsumer {
	return &ReviewConsumer{service: service}
}

func (h *ReviewConsumer) ConsumeReview(msg amqp.Delivery) error {

	var review dto.HotelReviewDTO
	if err := json.Unmarshal(msg.Body, &review); err != nil {
		logger.Logger.Error("failed to unmarshal review JSON", zap.Error(err))
		//	BAD data , drop the data , no need to resend it
		_ = msg.Nack(false, false)
		return err
	}

	if err := h.service.SaveReview(&review); err != nil {
		logger.Logger.Error("failed to save review", zap.Error(err))
		//	retry from broker end , as db operation failed
		_ = msg.Nack(false, true)
		return err
	}

	//only ACK after DB save succeededs
	// If true, this delivery and all prior unacknowledged deliverieson the same channel will be acknowledged.
	if err := msg.Ack(false); err != nil {
		logger.Logger.Error("failed to ack message", zap.Error(err))
	}

	return nil
}
