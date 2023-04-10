package config

import (
	"context"
	"testing"
)

func TestConfig_RedisClient(t *testing.T) {
	redisClient := NewRedisClient("localhost", "6379", "secret")
	var ctx = context.Background()

	err := redisClient.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		t.Error(err)
	}

	_, err = redisClient.Get(ctx, "key").Result()
	if err != nil {
		t.Error(err)
	}
}
