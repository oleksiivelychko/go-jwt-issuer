package issuer

import (
	"testing"
	"time"
)

func TestIssuer_ValidateIssuedToken(t *testing.T) {
	token, _, _, err := IssueJWT("secretkey", "jwt.account.local", "jwt.local", 1, 1)
	if err != nil {
		t.Fatal(err.Error())
	}

	// to validate expiration time
	time.Sleep(1 * time.Second)

	claimsJWT, err := ParseToken(token, "secretkey", "jwt.account.local", "jwt.local", 1)
	if _, ok := claimsJWT.Claims.(*ClaimsJWT); !ok || !claimsJWT.Valid {
		t.Fatal(err.Error())
	}
}

func TestIssuer_ValidateExpiredToken(t *testing.T) {
	token, _, _, err := IssueJWT("secretkey", "jwt.account.local", "jwt.local", 0, 1)
	if err != nil {
		t.Fatal(err.Error())
	}

	// to validate expiration time
	time.Sleep(1 * time.Second)

	_, err = ParseToken(token, "secretkey", "jwt.local", "jwt.local", 0)
	if err == nil {
		t.Error("unable to validate expired token")
	}

	t.Log(err)
}
