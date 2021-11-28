package issuer

import (
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"testing"
	"time"
)

func TestIssueAndValidateToken(t *testing.T) {
	env.InitEnv()

	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

	token, _, _, err := IssueUserJWT(secretKey, aud, iss, expiresMinutes, 1)
	if err != nil {
		t.Errorf("failed to issue signed token: %s", err.Error())
	}

	// to validate expiration time
	time.Sleep(1 * time.Second)

	exp := time.Now().Add(time.Minute * time.Duration(9)).Unix()
	jwtClaims, err := ValidateToken(token, secretKey, aud, iss, exp)
	if _, ok := jwtClaims.Claims.(*JwtClaims); !ok || !jwtClaims.Valid {
		t.Errorf("failed to validate signed token: %s", err.Error())
	}
}

func BenchmarkIssueAndValidateToken(b *testing.B) {
	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

	for i := 0; i < b.N; i++ {
		token, _, exp, _ := IssueUserJWT(secretKey, aud, iss, expiresMinutes, 1)
		_, _ = ValidateToken(token, secretKey, aud, iss, exp)
	}
}
