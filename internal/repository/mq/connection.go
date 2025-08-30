package mq

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
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
			log.Printf("connection failed: %v. Retrying...", err)
			time.Sleep(time.Duration(i+1) * time.Second)
		}
		if err != nil {
			err = fmt.Errorf("could not connect to RabbitMQ: %w", err)
			return
		}

		ch, err = conn.Channel()
		if err != nil {
			err = fmt.Errorf("could not open channel: %w", err)
			return
		}

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
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
