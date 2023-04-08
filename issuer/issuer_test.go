package issuer

import (
	"github.com/oleksiivelychko/go-jwt-issuer/config"
	"testing"
	"time"
)

func TestIssuer_IssueValidateToken(t *testing.T) {
	var secretKey = config.GetSecretKey()
	var aud = config.GetAudience()
	var iss = config.GetIssuer()
	var expiresMinutes = config.GetExpirationTimeMinutes()

	token, _, exp, err := IssueJWT(secretKey, aud, iss, expiresMinutes, 1)
	if err != nil {
		t.Errorf("unable to issue signed token: %s", err.Error())
	}

	// to validate expiration time
	time.Sleep(1 * time.Second)

	claimsJWT, err := ParseToken(token, secretKey, aud, iss, exp)
	if _, ok := claimsJWT.Claims.(*ClaimsJWT); !ok || !claimsJWT.Valid {
		t.Errorf("unable to validate signed token: %s", err.Error())
	}
}

func BenchmarkIssuer_IssueAndValidateToken(b *testing.B) {
	var secretKey = config.GetSecretKey()
	var aud = config.GetAudience()
	var iss = config.GetIssuer()
	var expiresMinutes = config.GetExpirationTimeMinutes()

	for i := 0; i < b.N; i++ {
		token, _, exp, _ := IssueJWT(secretKey, aud, iss, expiresMinutes, 1)
		_, _ = ParseToken(token, secretKey, aud, iss, exp)
	}
}
