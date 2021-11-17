package issuer

import (
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"os"
	"testing"
)

func TestIssueJWT(t *testing.T) {
	var secretKey = env.GetSecretKey()
	_, _, _, err := IssueJWT(secretKey, 0, 60, "", "")
	if err != nil {
		t.Errorf("failed to get the complete signed token: %s", err.Error())
	}
}

func BenchmarkIssueJWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _, _ = IssueJWT([]byte(os.Getenv("SECRET_KEY")), 0, 60, "", "")
	}
}
