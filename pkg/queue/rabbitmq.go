package queue

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
)

type ImageProcessingJob struct {
	ProductID uint     `json:"product_id"`
	Images    []string `json:"images"`
}

var (
	channel *amqp091.Channel
	queue   amqp091.Queue
)

func InitRabbitMQ(url string) error {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return err
	}

	channel, err = conn.Channel()
	if err != nil {
		return err
	}

	queue, err = channel.QueueDeclare(
		"image_processing",
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func PublishImageJob(ctx context.Context, job ImageProcessingJob) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return channel.PublishWithContext(ctx,
		"",
		queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}