package issuer

import (
	"testing"
)

func TestIssueJWT(t *testing.T) {
	_, _, _, err := IssueUserJWT(0)
	if err != nil {
		t.Errorf("failed to get the complete signed token: %s", err.Error())
	}
}

func BenchmarkIssueJWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _, _ = IssueUserJWT(0)
	}
}
