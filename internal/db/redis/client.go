package redisdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisClient(cfg Config) (*RedisClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Connected to Redis")
	return &RedisClient{Client: rdb, Ctx: context.Background()}, nil
}
