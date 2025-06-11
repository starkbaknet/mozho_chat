package redisdb

import (
	"fmt"
	"time"
)

func (r *RedisClient) SetUserOnline(userID string, ttl time.Duration) error {
	return r.Client.Set(r.Ctx, fmt.Sprintf("user:online:%s", userID), "1", ttl).Err()
}

func (r *RedisClient) IsUserOnline(userID string) (bool, error) {
	val, err := r.Client.Exists(r.Ctx, fmt.Sprintf("user:online:%s", userID)).Result()
	return val == 1, err
}

func (r *RedisClient) SetSessionToken(token string, userID string, ttl time.Duration) error {
	return r.Client.Set(r.Ctx, fmt.Sprintf("session:%s", token), userID, ttl).Err()
}

func (r *RedisClient) GetUserIDFromToken(token string) (string, error) {
	return r.Client.Get(r.Ctx, fmt.Sprintf("session:%s", token)).Result()
}
