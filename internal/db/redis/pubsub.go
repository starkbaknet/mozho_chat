package redisdb

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func (r *RedisClient) PublishMessage(roomID string, message string) error {
	channel := fmt.Sprintf("chat:room:%s", roomID)
	return r.Client.Publish(r.Ctx, channel, message).Err()
}

func (r *RedisClient) SubscribeRoom(ctx context.Context, roomID string) *redis.PubSub {
	channel := fmt.Sprintf("chat:room:%s", roomID)
	return r.Client.Subscribe(ctx, channel)
}
