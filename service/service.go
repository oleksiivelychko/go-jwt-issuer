package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"os"
	"strconv"
)

type Service struct {
	Redis *redis.Client
}

type CachedTokens struct {
	AccessUID  string `json:"access"`
	RefreshUID string `json:"refresh"`
}

func (service *Service) InitRedis() {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	db, _ := strconv.ParseInt(os.Getenv("REDIS_HOST"), 10, 32)
	service.Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_HOST"),
		DB:       int(db),
	})
}

func (service *Service) GenerateTokenPair(secretKey, aud, iss string, expiresMinutes, userID uint) (
	accessToken string,
	refreshToken string,
	exp int64,
	err error,
) {
	var accessUID, refreshUID string

	if accessToken, accessUID, exp, err = issuer.IssueUserJWT(secretKey, aud, iss, expiresMinutes, userID); err != nil {
		return
	}

	if refreshToken, refreshUID, _, err = issuer.IssueUserJWT(secretKey, aud, iss, expiresMinutes, userID); err != nil {
		return
	}

	cachedJSON, err := json.Marshal(CachedTokens{
		AccessUID:  accessUID,
		RefreshUID: refreshUID,
	})

	var ctx = context.Background()
	service.Redis.Set(ctx, fmt.Sprintf("token-%d", userID), string(cachedJSON), 0)

	return
}
