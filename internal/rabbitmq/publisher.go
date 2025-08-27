package mq

import (
	"context"
	"log"
	"time"
	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch       *amqp091.Channel
	exchange string
}

// NewPublisher creates a publisher bound to an exchange.
func NewPublisher(conn *Connection, exchange string, kind string) (*Publisher, error) {
	err := conn.Channel().ExchangeDeclare(
		exchange, // name
		kind,     // type: "fanout", "direct", "topic"
		true,     // durable
		false,    // auto-delete
		false,    // internal
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		ch:       conn.Channel(),
		exchange: exchange,
	}, nil
}


func (p *Publisher) Publish(routingKey string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.ch.PublishWithContext(
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
	if err != nil {
		return err
	}
	log.Printf("[mq] published message to %s with key %s", p.exchange, routingKey)
	return nil
}
