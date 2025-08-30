package mq

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
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

// Consume starts consuming messages and invokes handler per message.
func (c *Consumer) Consume(handler func([]byte) error) error {
	msgs, err := c.ch.Consume(
		c.queue,
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			fmt.Println("got msg in consume")
			if err := handler(msg.Body); err != nil {
				log.Printf("[mq] handler error: %v", err)
			}

		}
	}()

	log.Printf("[mq] consumer started for queue %s", c.queue)
	return nil
}
