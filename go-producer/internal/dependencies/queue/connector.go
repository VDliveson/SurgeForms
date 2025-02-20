package queue

import (
	"context"
	"log"
	"time"

	"github.com/VDliveson/SurgeForms/go-producer/constants"
	"github.com/VDliveson/SurgeForms/go-producer/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	Connection *amqp.Channel
	Ctx        context.Context
	Channel    *amqp.Channel
}

func ConnectQueue(ctx context.Context) (*Queue, error) {
	var err error
	url := utils.GetEnv("RABBITMQ", "amqp://localhost:5672")
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	log.Println("RabbitMQ Connection created successfully")

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	log.Printf("RabbitMQ Exchange declared successfully")
	queue := &Queue{Connection: ch, Channel: ch, Ctx: ctx}
	return queue, nil
}

func (queue *Queue) SendData(data []byte, service string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Publishing message to %s service ...", service)
	err := queue.Channel.PublishWithContext(ctx,
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
