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
		t.Errorf("unable to set key-value into RedisClient storage: %s", err)
	}

	_, err = redisClient.Get(ctx, "key").Result()
	if err != nil {
		t.Errorf("unable to get key-value from RedisClient storage:. %s", err)
	}
}
