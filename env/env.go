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
	RedisAddr      string
	RedisPassword  string
	RedisDb        int
}

func InitConfig() *Config {
	return &Config{
		Port:           GetPort(),
		SecretKey:      GetSecretKey(),
		ISS:            GetISS(),
		AUD:            GetAUD(),
		ExpiresMinutes: GetExpiresMinutes(),
		RedisAddr:      GetRedisAddr(),
		RedisPassword:  GetRedisPassword(),
		RedisDb:        GetRedisDb(),
	}
}

func (cfg *Config) InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDb,
	})
}

func GetSecretKey() string {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	if len(secretKey) == 0 {
		return ""
	}

	return string(secretKey)
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
	var port = []byte(os.Getenv("PORT"))
	if len(port) == 0 {
		return ":8080"
	}

	return ":" + string(port)
}

func GetRedisAddr() string {
	var redisHost = os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	var redisPort = os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	return fmt.Sprintf("%s:%s", redisHost, redisPort)
}

func GetRedisPassword() string {
	return os.Getenv("REDIS_PASSWORD")
}

func GetRedisDb() int {
	var redisDb = os.Getenv("REDIS_DB")
	if redisDb != "" {
		db, err := strconv.ParseInt(redisDb, 10, 32)
		if err == nil {
			return int(db)
		}
	}
	return 0
}
