package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
)

type Service struct {
	Env   *env.Config
	Redis *redis.Client
}

type CachedTokens struct {
	AccessUID  string `json:"access"`
	RefreshUID string `json:"refresh"`
}

func (service *Service) GenerateUserTokenPair(userID uint) (
	accessToken string,
	refreshToken string,
	exp int64,
	err error,
) {
	var accessUID, refreshUID string

	if accessToken, accessUID, exp, err = issuer.IssueUserJWT(
		service.Env.SecretKey, service.Env.AUD, service.Env.ISS, service.Env.ExpiresMinutes, userID,
	); err != nil {
		return
	}

	if refreshToken, refreshUID, _, err = issuer.IssueUserJWT(
		service.Env.SecretKey, service.Env.AUD, service.Env.ISS, service.Env.ExpiresMinutes, userID,
	); err != nil {
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
