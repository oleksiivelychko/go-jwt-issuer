package issuer

import (
	"github.com/golang-jwt/jwt"
	"time"
)

/*IssueJWT
https://en.wikipedia.org/wiki/JSON_Web_Token
*/
func IssueJWT(secretKey []byte, aud, iss string) (string, error) {
	claims := jwt.MapClaims{
		/*
			Expiration Time as `exp`.
			Identifies the expiration time on and after which the JWT must not be accepted for processing.
			The value must be a NumericDate:[9] either an integer or decimal, representing seconds past 1970-01-01 00:00:00Z.
		*/
		"exp": time.Now().Add(time.Minute * 1).Unix(),
	}

	if aud != "" {
		/*
			Audience as `aud`.
			Identifies the recipients that the JWT is intended for.
			Each principal intended to process the JWT must identify itself with a value in the audience claim.
			If the principal processing the claim does not identify itself with a value in the aud claim when this claim is present, then the JWT must be rejected.
		*/
		claims["aud"] = aud
	}

	if iss != "" {
		/*
			Issuer as `iss`.
			Identifies principal that issued the JWT.
		*/
		claims["iss"] = iss
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
