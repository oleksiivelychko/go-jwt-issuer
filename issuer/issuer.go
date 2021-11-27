package issuer

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type JwtClaims struct {
	ID  uint   `json:"id"`
	UID string `json:"uid"`
	jwt.StandardClaims
}

func IssueUserJWT(secretKey, aud, iss string, expiresMinutes, userID uint) (
	token string,
	uid string,
	exp int64,
	err error,
) {
	exp = time.Now().Add(time.Minute * time.Duration(expiresMinutes)).Unix()
	uid = uuid.New().String()

	claims := &JwtClaims{
		ID:  userID,
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
			Audience:  aud,
			Issuer:    iss,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(secretKey))

	return
}

func ValidateToken(token, secretKey, aud, iss string, exp int64) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if aud != "" {
			verifiedAudience := token.Claims.(*JwtClaims).VerifyAudience(aud, false)
			if !verifiedAudience {
				return nil, fmt.Errorf("failed to verify `aud` claim")
			}
		}

		if iss != "" {
			verifiedIssuer := token.Claims.(*JwtClaims).VerifyIssuer(iss, false)
			if !verifiedIssuer {
				return nil, fmt.Errorf("failed to verify `iss` claim")
			}
		}

		if exp > 0 {
			verifiedExpires := token.Claims.(*JwtClaims).VerifyExpiresAt(exp, false)
			if !verifiedExpires {
				return nil, fmt.Errorf("failed to verify `exp` claim")
			}
		}

		return []byte(secretKey), nil
	})
}
