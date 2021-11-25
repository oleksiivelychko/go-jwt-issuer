package service

import (
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"os"
	"testing"
)

func TestServiceGenerateUserTokenPair(t *testing.T) {
	_ = os.Setenv("SECRET_KEY", "secretkey")
	_ = os.Setenv("EXPIRES_MINUTES", "10")

	cfg := env.InitConfig()
	tokenService := Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	_, _, _, err := tokenService.GenerateUserTokenPair(1)
	if err != nil {
		t.Errorf("unable to generate pair tokens: %s", err.Error())
	}
}
