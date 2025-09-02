package mq

import (
	"sync"
	"time"

	logger "github.com/Arpitmovers/reviewservice/internal/logging"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type AmqpConnection struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	url     string
}

var instance *AmqpConnection
var once sync.Once

func NewConnection(url string) (*AmqpConnection, error) {
	var err error

	once.Do(func() {
		var conn *amqp091.Connection
		var ch *amqp091.Channel

		// Retry loop with backoff
		for i := 0; i < 5; i++ {
			conn, err = amqp091.Dial(url)
			if err == nil {
				break
			}
			logger.Logger.Error("failed to connect to RabbitMQ",
				zap.Int("attempt", i+1),
				zap.String("url", url),
				zap.Error(err),
			)
			time.Sleep(time.Duration(i+1) * time.Second)
		}

		if err != nil {
			logger.Logger.Error("could not establish RabbitMQ connection", zap.String("url", url), zap.Error(err))
			return
		}

		ch, err = conn.Channel()
		if err != nil {
			logger.Logger.Error("could not open RabbitMQ channel", zap.String("url", url), zap.Error(err))
			return
		}

		logger.Logger.Info("successfully connected to RabbitMQ",
			zap.String("url", url),
		)

		instance = &AmqpConnection{
			conn:    conn,
			channel: ch,
			url:     url,
		}
	})

	return instance, err
}

func (c *AmqpConnection) Channel() *amqp091.Channel {
	return c.channel
}

func (c *AmqpConnection) Close() {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			logger.Logger.Warn("failed to close RabbitMQ channel", zap.Error(err))
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			logger.Logger.Warn("failed to close RabbitMQ connection", zap.Error(err))
		}
	}
}
