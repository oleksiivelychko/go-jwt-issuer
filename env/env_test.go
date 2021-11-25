package env

import (
	"context"
	"github.com/oleksiivelychko/go-helper/in"
	"os"
	"testing"
)

func TestEnvConfig(t *testing.T) {
	_ = os.Setenv("SECRET_KEY", "secretkey")

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
	_ = os.Setenv("REDIS_HOST", "localhost")
	_ = os.Setenv("REDIS_PORT", "6379")
	_ = os.Setenv("REDIS_PASSWORD", "secret")
	_ = os.Setenv("REDIS_DB", "0")

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
