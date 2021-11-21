package issuer

import (
	"github.com/golang-jwt/jwt"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"os"
	"testing"
	"time"
)

func TestIssueAndValidateToken(t *testing.T) {
	_ = os.Setenv("SECRET_KEY", "secretkey")
	_ = os.Setenv("EXPIRES_MINUTES", "10")

	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

	token, _, _, err := IssueUserJWT(secretKey, aud, iss, expiresMinutes, 1)
	if err != nil {
		t.Errorf("failed to issue signed token: %s", err.Error())
	}

	exp := time.Now().Add(time.Minute * time.Duration(expiresMinutes)).Unix()
	jwtClaims, err := ValidateToken(token, secretKey, aud, iss, exp)
	if _, ok := jwtClaims.Claims.(jwt.MapClaims); !ok || !jwtClaims.Valid {
		t.Errorf("failed to validate signed token: %s", err.Error())
	}
}

func BenchmarkIssueAndValidateToken(b *testing.B) {
	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

	for i := 0; i < b.N; i++ {
		token, _, _, _ := IssueUserJWT(secretKey, aud, iss, expiresMinutes, 1)
		exp := time.Now().Add(time.Minute * time.Duration(expiresMinutes)).Unix()
		_, _ = ValidateToken(token, secretKey, aud, iss, exp)
	}
}
