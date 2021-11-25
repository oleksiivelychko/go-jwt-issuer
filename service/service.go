package service

import (
	"context"
	"encoding/json"
	"errors"
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

func (service *Service) ValidateCachedToken(claims *issuer.JwtClaims, isRefresh bool) error {
	var ctx = context.Background()
	cachedJSON, _ := service.Redis.Get(ctx, fmt.Sprintf("token-%d", claims.ID)).Result()
	cachedTokens := new(CachedTokens)
	err := json.Unmarshal([]byte(cachedJSON), cachedTokens)

	var tokenUID string
	if isRefresh {
		tokenUID = cachedTokens.RefreshUID
	} else {
		tokenUID = cachedTokens.AccessUID
	}

	if err != nil || tokenUID != claims.UID {
		return errors.New("unable to validate cached token")
	}

	return nil
}
