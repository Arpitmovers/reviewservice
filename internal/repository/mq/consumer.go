package mq

import (
	"fmt"
	"log"

	logger "github.com/Arpitmovers/reviewservice/internal/logging"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	prefetchCount = 5
	preFetchSize  = 0 // not using size based flow control
)

type Consumer struct {
	ch       *amqp091.Channel
	queue    string
	exchange string
	key      string
}

// NewConsumer declares a queue and binds it to the exchange.
func NewConsumer(conn *AmqpConnection, queue, exchange, key string) (*Consumer, error) {
	ch := conn.Channel()

	// Declare queue
	_, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		queue,    // queue name
		key,      // routing key
		exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		ch:       ch,
		queue:    queue,
		exchange: exchange,
		key:      key,
	}, nil
}

func (c *Consumer) Consume(handler func(amqp.Delivery) error) error {

	msgs, err := c.ch.Consume(
		c.queue,
		"",    // consumer tag
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}
	// set prefetch (must be done on same channel)
	if err := c.ch.Qos(prefetchCount, preFetchSize, false); err != nil {
		log.Fatalf("failed to set QoS: %v", err)
	}

	// consume messages
	go func() {
		for msg := range msgs {
			logger.Logger.Info("got msg in consume")
			if err := handler(msg); err != nil {
				logger.Logger.Error(" handler error", zap.Error(err))

			}

		}
	}()

	logger.Logger.Info("[mq] consumer started for queue %s", zap.String("queue", c.queue))
	return nil
}
