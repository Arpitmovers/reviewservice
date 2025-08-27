package mq

import (
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	url     string
}

func NewConnection(url string) (*Connection, error) {
	var conn *amqp091.Connection
	var ch *amqp091.Channel
	var err error

	// Retry loop with backoff
	for i := 0; i < 5; i++ {
		conn, err = amqp091.Dial(url)
		if err == nil {
			break
		}
		log.Printf("connection failed: %v. Retrying...", err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("could not connect to RabbitMQ: %w", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("could not open channel: %w", err)
	}

	return &Connection{
		conn:    conn,
		channel: ch,
		url:     url,
	}, nil
}

func (c *Connection) Channel() *amqp091.Channel {
	return c.channel
}

func (c *Connection) Close() {
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
