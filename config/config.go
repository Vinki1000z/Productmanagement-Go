package config

import (
	"os"
)

type Config struct {
	DatabaseURL      string
	RedisURL         string
	RabbitMQURL      string
	CloudinaryURL    string
	CloudinaryCloud  string
	CloudinaryAPIKey string
	CloudinarySecret string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/products?sslmode=disable"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379"),
		RabbitMQURL:      getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		CloudinaryURL:    getEnv("CLOUDINARY_URL", "https://dmvs9syar.api.cloudinary.com/v1_1/dmvs9syar"),
		CloudinaryCloud:  getEnv("CLOUDINARY_CLOUD", "dmvs9syar"),
		CloudinaryAPIKey: getEnv("CLOUDINARY_API_KEY", "723711842418993"),
		CloudinarySecret: getEnv("CLOUDINARY_SECRET", "z887ZYOxHhlSfBzLofoa5Tl-UEs"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
