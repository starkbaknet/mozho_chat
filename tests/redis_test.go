package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"mozho_chat/internal/db/redis"
)

var rdb *redisdb.RedisClient

func setupRedis(t *testing.T) {
	err := godotenv.Load("../.env")
	assert.NoError(t, err)

	client, err := redisdb.NewRedisClient(redisdb.Config{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	assert.NoError(t, err)
	rdb = client
}

// ---- 1. Pub/Sub test ----

func TestPublishSubscribe(t *testing.T) {
	setupRedis(t)

	roomID := "test-room"
	testMessage := "hello chat!"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sub := rdb.SubscribeRoom(ctx, roomID)
	defer sub.Close()

	received := make(chan string)

	go func() {
		for msg := range sub.Channel() {
			received <- msg.Payload
			return
		}
	}()

	time.Sleep(500 * time.Millisecond)

	err := rdb.PublishMessage(roomID, testMessage)
	assert.NoError(t, err)

	select {
	case msg := <-received:
		assert.Equal(t, testMessage, msg)
	case <-time.After(3 * time.Second):
		t.Fatal("Did not receive published message")
	}
}

// ---- 2. User Online Cache ----

func TestUserOnlineCache(t *testing.T) {
	setupRedis(t)

	userID := "user-123"
	ttl := 2 * time.Second

	err := rdb.SetUserOnline(userID, ttl)
	assert.NoError(t, err)

	online, err := rdb.IsUserOnline(userID)
	assert.NoError(t, err)
	assert.True(t, online)

	time.Sleep(ttl + 1*time.Second)

	online, err = rdb.IsUserOnline(userID)
	assert.NoError(t, err)
	assert.False(t, online)
}

// ---- 3. Session Token Test ----

func TestSessionToken(t *testing.T) {
	setupRedis(t)

	token := "token-abc-xyz"
	userID := "user-456"
	ttl := 5 * time.Second

	err := rdb.SetSessionToken(token, userID, ttl)
	assert.NoError(t, err)

	result, err := rdb.GetUserIDFromToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, result)
}
