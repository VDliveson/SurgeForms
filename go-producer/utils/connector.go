package utils

import (
	"context"
	"log"
	"time"

	"github.com/VDliveson/SurgeForms/go-producer/constants"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var ch *amqp.Channel

func ConnectQueue() error {
	var err error
	url := GetEnv("RABBITMQ", "amqp://localhost:5672")
	conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}
	log.Println("RabbitMQ Connection created successfully")

	ch, err = conn.Channel()
	if err != nil {
		return err
	}

	log.Printf("RabbitMQ Channel created successfully")

	err = ch.ExchangeDeclare(
		constants.Exchange, // name
		"direct",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	log.Printf("RabbitMQ Exchange declared successfully")

	// defer ch.Close()
	// defer conn.Close()
	return nil
}

func SendData(data []byte, service string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Publishing message...")
	err := ch.PublishWithContext(ctx,
		constants.Exchange, // name
		service,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
	if err != nil {
		return err
	}

	log.Println("Message sent successfully")
	return nil
}
