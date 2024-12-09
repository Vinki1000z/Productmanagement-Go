package cache

import (
	"context"
	"encoding/json"
	"productManagmentBackend/models"
	"time"

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
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	key := "product:" + string(product.ID)
	return RedisClient.Set(ctx, key, data, 24*time.Hour).Err()
}

func GetProduct(id uint) (*models.Product, error) {
	key := "product:" + string(id)
	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var product models.Product
	err = json.Unmarshal(data, &product)
	return &product, err
}

func InvalidateProduct(id uint) error {
	key := "product:" + string(id)
	return RedisClient.Del(ctx, key).Err()
}
