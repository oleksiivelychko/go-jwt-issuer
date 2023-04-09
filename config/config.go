package config

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
	"time"
)

const autoLogOffMinutes = 10
const defaultIssuerISS = "jwt.local"
const defaultAudienceAUD = "account." + defaultIssuerISS
const defaultExpirationTimeMinutes = 1
const defaultRedisPort = "6379"
const defaultSecret = "secret"
const defaultServerName = "localhost"

type Config struct {
	SecretKey string
	AUD       string
	ISS       string
	EXP       uint
}

func NewConfig() *Config {
	return &Config{
		SecretKey: GetSecretKey(),
		ISS:       GetIssuer(),
		AUD:       GetAudience(),
		EXP:       GetExpirationTimeMinutes(),
	}
}

func (config *Config) GetAutoLogOffDuration() time.Duration {
	return time.Minute * autoLogOffMinutes
}

func InitRedis() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = defaultServerName
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = defaultRedisPort
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = defaultSecret
	}

	redisURL := fmt.Sprintf("redis://:%s@%s:%s", redisPassword, redisHost, redisPort)
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(options)
}

func GetSecretKey() string {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = defaultSecret
	}

	return secretKey
}

func GetAudience() string {
	audienceAUD := os.Getenv("AUDIENCE_AUD")
	if audienceAUD == "" {
		audienceAUD = defaultAudienceAUD
	}

	return audienceAUD
}

func GetIssuer() string {
	issuerISS := os.Getenv("ISSUER_ISS")
	if issuerISS == "" {
		issuerISS = defaultIssuerISS
	}

	return issuerISS
}

func GetExpirationTimeMinutes() uint {
	exp := os.Getenv("EXPIRATION_TIME_EXP")
	if exp != "" {
		minutes, err := strconv.ParseUint(exp, 10, 32)
		if err == nil {
			return uint(minutes)
		}
	}

	return defaultExpirationTimeMinutes
}

func GetServerAddr() string {
	if host, isOk := os.LookupEnv("HOST"); isOk == true {
		_ = os.Setenv("HOST", host)
	} else {
		_ = os.Setenv("HOST", defaultServerName)
	}

	return fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
}
