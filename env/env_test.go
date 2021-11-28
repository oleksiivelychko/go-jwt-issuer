package env

import (
	"context"
	"github.com/oleksiivelychko/go-helper/in"
	"testing"
)

func TestEnvConfig(t *testing.T) {
	InitEnv()

	cfg := InitConfig()

	var allowedRangePorts = []string{":80", ":8080", ":443"}
	_, result := in.StringIn(allowedRangePorts, cfg.Port)
	if !result {
		t.Errorf("given $PORT %s is not acceptable", cfg.Port)
	}

	if cfg.SecretKey == "" {
		t.Errorf("environment variable `SECRET_KEY` is not defined")
	}
}

func TestRedisConnection(t *testing.T) {
	InitEnv()

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
