package main

import (
	"log"
	"productManagmentBackend/config"
	"productManagmentBackend/database"
	"productManagmentBackend/pkg/cache"
	"productManagmentBackend/pkg/queue"
	"productManagmentBackend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize the database
	if err := database.ConnectDatabase(); err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Initialize Redis
	if err := cache.InitRedis(cfg.RedisURL); err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}

	// Initialize RabbitMQ
	if err := queue.InitRabbitMQ(cfg.RabbitMQURL); err != nil {
		log.Fatal("Failed to connect to RabbitMQ: ", err)
	}

	// Initialize the Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}