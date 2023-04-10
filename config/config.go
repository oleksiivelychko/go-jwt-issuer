package config

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
	"time"
)

const defaultAutoLogOutMinutes = time.Minute * 10

type Config struct {
	SecretKey string
	AUD       string
	ISS       string
	EXP       uint
}

func NewConfig(secretKey, aud, iss, exp string) *Config {
	if _, isOk := os.LookupEnv("SECRET_KEY"); !isOk {
		_ = os.Setenv("SECRET_KEY", secretKey)
	}
	if _, isOk := os.LookupEnv("ISSUER_ISS"); !isOk {
		_ = os.Setenv("ISSUER_ISS", iss)
	}
	if _, isOk := os.LookupEnv("AUDIENCE_AUD"); !isOk {
		_ = os.Setenv("AUDIENCE_AUD", aud)
	}
	if _, isOk := os.LookupEnv("EXPIRATION_TIME_EXP"); !isOk {
		_ = os.Setenv("EXPIRATION_TIME_EXP", exp)
	}

	return &Config{
		SecretKey: secretKey,
		ISS:       iss,
		AUD:       aud,
		EXP:       ParseExpirationTime(exp),
	}
}

func (config *Config) GetAutoLogOutMinutes() time.Duration {
	return defaultAutoLogOutMinutes
}

func NewRedisClient(host, port, password string) *redis.Client {
	redisURL := fmt.Sprintf("redis://:%s@%s:%s", password, host, port)
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(options)
}

func ParseExpirationTime(exp string) uint {
	parseUint, err := strconv.ParseUint(exp, 10, 32)
	if err != nil {
		panic(err)
	}

	return uint(parseUint)
}
