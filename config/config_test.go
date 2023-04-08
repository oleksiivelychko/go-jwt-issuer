package config

import (
	"context"
	"testing"
)

func TestConfig_InitRedis(t *testing.T) {
	redisClient := InitRedis()
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
