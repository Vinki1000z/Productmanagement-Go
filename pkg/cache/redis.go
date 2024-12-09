package cache

import (
	"context"
	"encoding/json"
	"productManagmentBackend/models"
	"time"

	"strconv"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis(redisURL string) error {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}

	RedisClient = redis.NewClient(opt)
	return RedisClient.Ping(ctx).Err()
}

func SetProduct(product *models.Product) error {
	// Convert product ID to a string properly
	key := "product:" + strconv.Itoa(int(product.ID))

	// Marshal product data into JSON
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	// Store product data in Redis with a 24-hour expiration
	return RedisClient.Set(ctx, key, data, 24*time.Hour).Err()
}

// GetProduct retrieves a product from Redis by ID
func GetProduct(id uint) (*models.Product, error) {
	// Properly convert the uint ID to a string
	key := "product:" + strconv.Itoa(int(id))

	// Fetch product data from Redis
	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON into the Product model
	var product models.Product
	err = json.Unmarshal(data, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
func InvalidateProduct(id uint) error {
	// Properly convert the uint ID to a string
	key := "product:" + strconv.Itoa(int(id))

	// Delete the key from Redis
	return RedisClient.Del(ctx, key).Err()
}
