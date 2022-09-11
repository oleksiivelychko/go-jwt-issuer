package env

import (
	"context"
	"testing"
)

func TestRedisConnection(t *testing.T) {
	SetDefaults()

	cfg := InitConfig()
	redis := cfg.InitRedis()

	var ctx = context.Background()

	err := redis.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		t.Errorf("unable to set key-value into redis storage: %s", err)
	}

	val, err := redis.Get(ctx, "key").Result()
	if err != nil {
		t.Errorf("unable to get value by key from redis storage, got: %s. %s", val, err)
	}
}
