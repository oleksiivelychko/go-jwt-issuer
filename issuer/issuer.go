package issuer

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"time"
)

type JwtClaims struct {
	ID  uint   `json:"id"`
	UID string `json:"uid"`
	jwt.StandardClaims
}

func IssueUserJWT(userID uint) (
	token string,
	uid string,
	exp int64,
	err error,
) {
	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

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
	token, err = jwtToken.SignedString(secretKey)

	return
}
