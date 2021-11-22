package service

import (
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"os"
	"testing"
)

func TestServiceGenerateTokenPair(t *testing.T) {
	_ = os.Setenv("SECRET_KEY", "secretkey")
	_ = os.Setenv("EXPIRES_MINUTES", "10")

	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

	service := Service{}
	service.InitRedis()

	_, _, _, err := service.GenerateTokenPair(secretKey, aud, iss, expiresMinutes, 1)
	if err != nil {
		t.Errorf("unable to generate pair tokens: %s", err.Error())
	}
}
