package main

import (
	"context"
	"encoding/json"
	"log"
	"productManagmentBackend/config"
	"productManagmentBackend/database"
	"productManagmentBackend/models"
	"productManagmentBackend/pkg/imageprocessor"
	"productManagmentBackend/pkg/logger"

	"github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type ImageProcessingJob struct {
	ProductID uint     `json:"product_id"`
	Images    []string `json:"images"`
}

func main() {
	cfg := config.LoadConfig()

	// Initialize database connection
	if err := database.ConnectDatabase(); err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Initialize Cloudinary
	cloudinary, err := imageprocessor.NewCloudinaryService(cfg.CloudinaryAPIKey, cfg.CloudinarySecret, cfg.CloudinaryCloud)
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}

	// Connect to RabbitMQ
	conn, err := amqp091.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ: ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel: ", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"image_processing",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare queue: ", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to register consumer: ", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var job ImageProcessingJob
			if err := json.Unmarshal(d.Body, &job); err != nil {
				logger.Log.WithError(err).Error("Failed to unmarshal job")
				continue
			}

			processImages(context.Background(), cloudinary, job)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func processImages(ctx context.Context, cloudinary *imageprocessor.CloudinaryService, job ImageProcessingJob) {
	var processedImages []string

	for _, imageURL := range job.Images {
		processedURL, err := cloudinary.ProcessImage(ctx, imageURL)
		if err != nil {
			logger.Log.WithError(err).Error("Failed to process image")
			continue
		}
		processedImages = append(processedImages, processedURL)

		// Log each processed URL to the console
		log.Printf("Processed image URL: %s", processedURL)
	}

	if len(processedImages) > 0 {
		// Update the product with processed images
		if err := database.DB.Model(&models.Product{}).
			Where("id = ?", job.ProductID).
			Update("compressed_product_images", gorm.Expr("?::text[]", pq.Array(processedImages))).Error; err != nil {
			logger.Log.WithError(err).Error("Failed to update product with processed images")
		}
	}
}
