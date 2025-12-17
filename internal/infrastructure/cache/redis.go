package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	UseTLS   bool
}

func NewRedisClient(cfg RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}

// func NewRedisClient(addr, password string) (*redis.Client, error) {
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     addr,
// 		Password: password,
// 		DB:       0,
// 	})

// 	if err := client.Ping(context.Background()).Err(); err != nil {
// 		return nil, err
// 	}
// 	return client, nil
// }

// Example helpers
func Set(ctx context.Context, client *redis.Client, key string, value interface{}, ttl time.Duration) error {
	return client.Set(ctx, key, value, ttl).Err()
}

func Get(ctx context.Context, client *redis.Client, key string) (string, error) {
	return client.Get(ctx, key).Result()
}
