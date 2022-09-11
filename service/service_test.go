package service

import (
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"testing"
)

func TestServiceUserTokenPair(t *testing.T) {
	env.SetDefaults()

	cfg := env.InitConfig()
	tokenService := Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	accessToken, refreshToken, exp, err := tokenService.GenerateUserTokenPair(1)
	if err != nil {
		t.Errorf("unable to generate pair tokens: %s", err.Error())
	}

	claims, err := tokenService.ValidateParsedToken(accessToken, exp)
	if err != nil {
		t.Errorf("unable to validate parsed access token: %s", err.Error())
	}

	err = tokenService.ValidateCachedToken(claims, false)
	if err != nil {
		t.Errorf("unable to validate cached access token: %s", err.Error())
	}

	claims, err = tokenService.ValidateParsedToken(refreshToken, exp)
	if err != nil {
		t.Errorf("unable to validate parsed refresh token: %s", err.Error())
	}

	err = tokenService.ValidateCachedToken(claims, true)
	if err != nil {
		t.Errorf("unable to validate cached refresh token: %s", err.Error())
	}

	err = tokenService.ClearCachedToken(claims)
	if err != nil {
		t.Errorf("unable to clear cached user token pair: %s", err.Error())
	}
}
