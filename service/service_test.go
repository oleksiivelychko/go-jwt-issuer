package service

import (
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"os"
	"testing"
	"time"
)

func initEnv() {
	_ = os.Setenv("REDIS_HOST", "localhost")
	_ = os.Setenv("REDIS_PORT", "6379")
	_ = os.Setenv("REDIS_PASSWORD", "secret")
	_ = os.Setenv("REDIS_DB", "0")
	_ = os.Setenv("SECRET_KEY", "secretkey")
	_ = os.Setenv("EXPIRES_MINUTES", "10")
	_ = os.Setenv("AUDIENCE_AUD", "account.jwt.local")
	_ = os.Setenv("ISSUER_ISS", "jwt.local")
}

func TestServiceGenerateUserTokenPair(t *testing.T) {
	initEnv()

	cfg := env.InitConfig()
	tokenService := Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	accessToken, refreshToken, exp, err := tokenService.GenerateUserTokenPair(1)
	if err != nil {
		t.Errorf("unable to generate pair tokens: %s", err.Error())
	}

	// to validate expiration time
	time.Sleep(1 * time.Second)

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
}
