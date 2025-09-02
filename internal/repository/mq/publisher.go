package mq

import (
	"context"
	"fmt"
	"time"

	logger "github.com/Arpitmovers/reviewservice/internal/logging"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Publisher struct {
	ch       *amqp091.Channel
	exchange string
}

const (
	maxPublishRetries = 5
	basePublishDelay  = 1 * time.Minute
)

func NewPublisher(conn *AmqpConnection, exchange string, kind string) (*Publisher, error) {
	err := conn.Channel().ExchangeDeclare(
		exchange,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Logger.Error("failed to declare exchange", zap.String("exchange", exchange), zap.String("kind", kind), zap.Error(err))
		return nil, err
	}

	logger.Logger.Info("exchange declared", zap.String("exchange", exchange), zap.String("kind", kind))
	return &Publisher{ch: conn.Channel(), exchange: exchange}, nil
}

// PublishSafe with exponential backoff
func (p *Publisher) PublishSafe(routingKey string, body []byte) error {
	var err error

	for attempt := 1; attempt <= maxPublishRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = p.ch.PublishWithContext(ctx, p.exchange, routingKey, false, false, amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		cancel()

		if err == nil {
			logger.Logger.Info("message published", zap.String("exchange", p.exchange), zap.String("routingKey", routingKey), zap.Int("attempt", attempt))
			return nil
		}

		wait := basePublishDelay * (1 << (attempt - 1))
		logger.Logger.Error("publish failed", zap.Int("attempt", attempt), zap.Int("maxRetries", maxPublishRetries), zap.Duration("retryIn", wait), zap.String("exchange", p.exchange), zap.String("routingKey", routingKey), zap.Error(err))
		time.Sleep(wait)
	}

	logger.Logger.Error("publish retries exhausted", zap.Int("maxRetries", maxPublishRetries), zap.String("exchange", p.exchange), zap.String("routingKey", routingKey), zap.Error(err))
	return fmt.Errorf("failed to publish after %d retries: %w", maxPublishRetries, err)
}
