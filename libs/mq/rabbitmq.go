package mq

import (
    "context"
    "encoding/json"
    "fmt"

    amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
    conn    *amqp.Connection
    Channel *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
    conn, err := amqp.Dial(uri)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to open channel: %v", err)
    }

    return &RabbitMQ{
        conn:    conn,
        Channel: ch,
    }, nil
}

func (r *RabbitMQ) Close() {
    r.Channel.Close()
    r.conn.Close()
}

func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
    return r.Channel.QueueDeclare(
        name,
        true,
        false,
        false,
        false,
        nil,
    )
}

func (r *RabbitMQ) DeclareExchange(name, kind string) error {
    return r.Channel.ExchangeDeclare(
        name,
        kind,
        true,
        false,
        false,
        false,
        nil,
    )
}

func (r *RabbitMQ) BindQueue(queue, exchange, routingKey string) error {
    return r.Channel.QueueBind(
        queue,
        routingKey,
        exchange,
        false,
        nil,
    )
}

func (r *RabbitMQ) PublishJSON(ctx context.Context, exchange, routingKey string, data interface{}) error {
    body, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("json marshal failed: %v", err)
    }

    return r.Channel.PublishWithContext(
        ctx,
        exchange,
        routingKey,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}

func (r *RabbitMQ) Consume(queue string) (<-chan amqp.Delivery, error) {
    return r.Channel.Consume(
        queue,
        "",
        false, // manual ack
        false,
        false,
        false,
        nil,
    )
}
