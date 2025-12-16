package bootstrap

import (
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/dhanarrizky/Golang-template/internal/config"
	"github.com/dhanarrizky/Golang-template/internal/infrastructure/cache"
)

func InitRedis(cfg *config.Config) *redis.Client {
	if !cfg.RedisEnable {
		log.Println("[BOOTSTRAP] Redis disabled")
		return nil
	}

	client, err := cache.NewRedisClient(cache.RedisConfig{
		Host:     cfg.RedisHost,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
		UseTLS:   cfg.RedisTLS,
	})
	if err != nil {
		log.Fatal("failed to connect redis:", err)
	}

	log.Println("[BOOTSTRAP] Redis connected")
	return client
}
