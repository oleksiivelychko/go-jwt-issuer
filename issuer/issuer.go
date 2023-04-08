package issuer

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type ClaimsJWT struct {
	ID  uint   `json:"id"`
	UID string `json:"uid"`
	jwt.StandardClaims
}

type ResponseJWT struct {
	AccessToken    string `json:"accessToken,omitempty"`
	RefreshToken   string `json:"refreshToken,omitempty"`
	ExpirationTime string `json:"expirationTime,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
	ErrorCode      uint8  `json:"errorCode,omitempty"`
}

func IssueJWT(secretKey, aud, iss string, expMinutes, userID uint) (
	token string,
	uid string,
	exp int64,
	err error,
) {
	exp = time.Now().Add(time.Minute * time.Duration(expMinutes)).Unix()
	uid = uuid.New().String()

	claims := &ClaimsJWT{
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

func ParseToken(token, secretKey, aud, iss string, exp int64) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &ClaimsJWT{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if aud != "" {
			verifiedAudience := token.Claims.(*ClaimsJWT).VerifyAudience(aud, false)
			if !verifiedAudience {
				return nil, fmt.Errorf("unable to verify 'aud' claim")
			}
		}

		if iss != "" {
			verifiedIssuer := token.Claims.(*ClaimsJWT).VerifyIssuer(iss, false)
			if !verifiedIssuer {
				return nil, fmt.Errorf("unable to verify 'iss' claim")
			}
		}

		if exp > 0 {
			verifiedExpiration := token.Claims.(*ClaimsJWT).VerifyExpiresAt(exp, false)
			if !verifiedExpiration {
				return nil, fmt.Errorf("unable to verify 'exp' claim")
			}
		}

		return []byte(secretKey), nil
	})
}
