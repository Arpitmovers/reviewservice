package mq

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
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
		exchange, // name
		kind,     // type: "fanout", "direct", "topic"
		true,     // durable
		false,    // auto-delete ( exchanges get deleted 	when consumer unsubscribes)
		false,    // internal
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		fmt.Println("NewPublisher error is", err)
		return nil, err
	}

	return &Publisher{
		ch:       conn.Channel(),
		exchange: exchange,
	}, nil
}

// Publish with exponential backoff in minutes
func (p *Publisher) PublishSafe(routingKey string, body []byte) error {
	var err error

	for attempt := 1; attempt <= maxPublishRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = p.ch.PublishWithContext(
			ctx,
			p.exchange, // exchange
			routingKey, // routing key
			false,      // mandatory
			false,      // immediate
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		cancel()

		if err == nil {
			log.Printf("[mq] published message to %s with key %s", p.exchange, routingKey)
			return nil
		}

		wait := basePublishDelay * (1 << (attempt - 1))
		log.Printf("[mq] publish failed (attempt %d/%d), retrying in %v: %v", attempt, maxPublishRetries, wait, err)

		time.Sleep(wait)
	}

	// all retries failed
	return fmt.Errorf("failed to publish after %d retries: %w", maxPublishRetries, err)
}
