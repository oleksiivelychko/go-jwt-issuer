package issuer

import (
	"github.com/golang-jwt/jwt"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"os"
	"testing"
)

func TestIssueAndValidateToken(t *testing.T) {
	_ = os.Setenv("SECRET_KEY", "secretkey")
	_ = os.Setenv("EXPIRES_MINUTES", "10")

	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

	token, _, exp, err := IssueUserJWT(secretKey, aud, iss, expiresMinutes, 1)
	if err != nil {
		t.Errorf("failed to get signed token: %s", err.Error())
	}

	jwtClaims, err := ValidateToken(token, secretKey, aud, iss, exp)
	if _, ok := jwtClaims.Claims.(jwt.MapClaims); !ok || !jwtClaims.Valid {
		t.Errorf("failed to validate signed token: %s", err.Error())
	}
}

func BenchmarkIssueJWT(b *testing.B) {
	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

	for i := 0; i < b.N; i++ {
		token, _, exp, _ := IssueUserJWT(secretKey, aud, iss, expiresMinutes, 1)
		_, _ = ValidateToken(token, secretKey, aud, iss, exp)
	}
}
