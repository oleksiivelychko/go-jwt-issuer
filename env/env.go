package env

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

const AutoLogoffMinutes = 10

type Config struct {
	Port           string
	SecretKey      string
	AUD            string
	ISS            string
	ExpiresMinutes uint
	RedisUrl       string
}

func InitConfig() *Config {
	return &Config{
		Port:           GetPort(),
		SecretKey:      GetSecretKey(),
		ISS:            GetISS(),
		AUD:            GetAUD(),
		ExpiresMinutes: GetExpiresMinutes(),
		RedisUrl:       GetRedisUrl(),
	}
}

func (cfg *Config) InitRedis() *redis.Client {
	options, _ := redis.ParseURL(GetRedisUrl())
	return redis.NewClient(options)
}

func GetSecretKey() string {
	var secretKey = os.Getenv("SECRET_KEY")
	if len(secretKey) == 0 {
		return ""
	}

	return secretKey
}

func GetAUD() string {
	return os.Getenv("AUDIENCE_AUD")
}

func GetISS() string {
	return os.Getenv("ISSUER_ISS")
}

func GetExpiresMinutes() uint {
	expirationTimeExp := os.Getenv("EXPIRES_MINUTES")
	if expirationTimeExp != "" {
		minutes, err := strconv.ParseUint(expirationTimeExp, 10, 32)
		if err == nil {
			return uint(minutes)
		}
	}
	return 60
}

func GetPort() string {
	var port = os.Getenv("PORT")
	if len(port) == 0 {
		return ":8080"
	}

	return fmt.Sprintf(":%s", port)
}

func GetRedisUrl() string {
	var redisTlsUrl = os.Getenv("REDIS_TLS_URL")
	if redisTlsUrl != "" {
		return redisTlsUrl
	}

	var redisUrl = os.Getenv("REDIS_URL")
	if redisUrl != "" {
		return redisUrl
	}

	var redisHost = os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	var redisPort = os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	var redisPassword = os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = "secret"
	}

	return fmt.Sprintf("redis://:%s@%s:%s", redisPassword, redisHost, redisPort)
}
