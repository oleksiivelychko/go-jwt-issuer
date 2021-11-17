package issuer

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type JwtClaims struct {
	ID  uint   `json:"id"`
	UID string `json:"uid"`
	jwt.StandardClaims
}

func IssueJWT(secretKey []byte, userID uint, expireMinutes int, aud, iss string) (
	token string,
	uid string,
	exp int64,
	err error,
) {
	exp = time.Now().Add(time.Minute * time.Duration(expireMinutes)).Unix()
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
