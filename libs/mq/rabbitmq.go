package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)


type RabbitMQ struct {
	conn *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Rabbitmq %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create channel: %v", err)
	}

	rmq := &RabbitMQ{
		conn: conn,
		Channel: ch,
	}

	return rmq, nil
}

func (r *RabbitMQ) publish(ctx context.Context, exchange, routingKey string, msg amqp.Publishing) {
	return r.Channel.PublishWithContext(ctx,
		exchange,
		routingKey,
		false,
		false,
		msg,
	)
}
