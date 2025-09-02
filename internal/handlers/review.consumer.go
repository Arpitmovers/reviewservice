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
	logger.Logger.Info("Received review message", zap.ByteString("body", msg.Body))

	var review dto.Review
	if err := json.Unmarshal(msg.Body, &review); err != nil {
		logger.Logger.Error("failed to unmarshal review JSON", zap.Error(err))
		//message will be dropped by broker , as we have not configured Lead letter queue
		_ = msg.Nack(false, false)
		return err
	}

	if err := h.service.SaveReview(&review); err != nil {
		logger.Logger.Error("failed to save review", zap.Error(err))
		// nack with requeue = true so message is retried from broker
		_ = msg.Nack(false, true)
		return err
	}

	// only ACK after DB save succeeded
	if err := msg.Ack(false); err != nil {
		logger.Logger.Error("failed to ack message", zap.Error(err))
	}
	return nil
}
