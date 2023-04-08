package token

import (
	"github.com/oleksiivelychko/go-jwt-issuer/config"
	"testing"
)

func TestService_GeneratePairTokens(t *testing.T) {
	tokenService := Service{
		Config:      config.NewConfig(),
		RedisClient: config.InitRedis(),
	}

	accessToken, refreshToken, exp, err := tokenService.GeneratePairTokens(1)
	if err != nil {
		t.Fatal(err.Error())
	}

	claimsJWT, err := tokenService.ValidateToken(accessToken, exp)
	if err != nil {
		t.Error(err.Error())
	}

	err = tokenService.ValidateCachedToken(claimsJWT, false)
	if err != nil {
		t.Error(err.Error())
	}

	claimsJWT, err = tokenService.ValidateToken(refreshToken, exp)
	if err != nil {
		t.Error(err.Error())
	}

	err = tokenService.ValidateCachedToken(claimsJWT, true)
	if err != nil {
		t.Error(err.Error())
	}

	err = tokenService.ClearCachedTokens(claimsJWT)
	if err != nil {
		t.Error(err.Error())
	}
}
