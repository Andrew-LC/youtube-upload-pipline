package main

import (
	"os"
	"log"
	"github.com/Andrew-LC/libs/mq"
)

const (
	mqURI = "MQ_URI"
)

func main() {
	os.Setenv(mqURI, "amqp://guest:guest@localhost:5672/")
	mqUri := os.Getenv(mqURI)

	rabbitMQ, err := mq.NewRabbitMQ(mqUri)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	_ = rabbitMQ.DeclareExchange("upload_events", "direct")

	queue, err := rabbitMQ.DeclareQueue("")
	if err != nil {
		log.Fatalf("Failed to create a queue")
	}

	rabbitMQ.BindQueue(queue.Name, "upload_events", "upload.created")

	msgs, err := rabbitMQ.Consume(queue.Name)
	if err != nil {
		log.Fatalf("Failed to register a consumer")
	}

	var forever chan struct{}

        go func() {
                for d := range msgs {
                        log.Printf(" [x] %s", d.Body)
                }
        }()

        log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
        <-forever
}
