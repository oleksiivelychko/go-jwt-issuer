package token

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oleksiivelychko/go-jwt-issuer/config"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
)

type Service struct {
	Config      *config.Config
	RedisClient *redis.Client
}

type CachedTokens struct {
	AccessUID  string `json:"accessUID"`
	RefreshUID string `json:"refreshUID"`
}

func NewService(config *config.Config, redis *redis.Client) *Service {
	return &Service{Config: config, RedisClient: redis}
}

func (service *Service) GeneratePairTokens(userID uint) (
	accessToken string,
	refreshToken string,
	exp int64,
	err error,
) {
	var accessUID, refreshUID string

	if accessToken, accessUID, exp, err = issuer.IssueJWT(
		service.Config.SecretKey, service.Config.AUD, service.Config.ISS, service.Config.EXP, userID,
	); err != nil {
		return
	}

	if refreshToken, refreshUID, _, err = issuer.IssueJWT(
		service.Config.SecretKey, service.Config.AUD, service.Config.ISS, service.Config.EXP, userID,
	); err != nil {
		return
	}

	cachedJSON, err := json.Marshal(CachedTokens{
		AccessUID:  accessUID,
		RefreshUID: refreshUID,
	})

	cmd := service.RedisClient.Set(
		context.Background(),
		fmt.Sprintf("token-%d", userID), string(cachedJSON),
		service.Config.GetAutoLogOffDuration(),
	)

	if cmd.Err() != nil {
		err = cmd.Err()
	}

	return
}

func (service *Service) ValidateToken(token string, exp int64) (claims *issuer.ClaimsJWT, err error) {
	parsedToken, err := issuer.ParseToken(token, service.Config.SecretKey, service.Config.AUD, service.Config.ISS, exp)
	if err != nil {
		return
	}

	var isOk = false
	if claims, isOk = parsedToken.Claims.(*issuer.ClaimsJWT); isOk && parsedToken.Valid {
		return
	}

	return claims, errors.New("unable to get claims of token")
}

func (service *Service) ValidateCachedToken(claims *issuer.ClaimsJWT, isRefresh bool) error {
	var ctx = context.Background()

	cachedJSON, _ := service.RedisClient.Get(ctx, fmt.Sprintf("token-%d", claims.ID)).Result()
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

	cmd := service.RedisClient.Expire(ctx, fmt.Sprintf("token-%d", claims.ID), service.Config.GetAutoLogOffDuration())
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (service *Service) ClearCachedTokens(claims *issuer.ClaimsJWT) error {
	cmd := service.RedisClient.Del(context.Background(), fmt.Sprintf("token-%d", claims.ID))
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}
