package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"net/http"
)

func AllowToEndpoint(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if validate(w, r) {
			endpoint(w, r)
		}
	})
}

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if validate(w, r) {
			next.ServeHTTP(w, r)
		}
	})
}

func validate(w http.ResponseWriter, r *http.Request) bool {
	secretKey := env.GetSecretKey()
	aud := env.GetAUD()
	iss := env.GetISS()
	tokenHeader := r.Header.Get("Authorization")

	if len(secretKey) == 0 {
		_, _ = fmt.Fprintf(w, "environment variable `SECRET_KEY` is not defined")
	} else if tokenHeader == "" {
		_, _ = fmt.Fprintf(w, "failed to get token from header request")
	} else {
		token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			if env.GetAUD() != "" {
				verifiedAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
				if !verifiedAudience {
					return nil, fmt.Errorf("failed to verify `aud` claim")
				}
			}

			if iss != "" {
				verifiedIssuer := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !verifiedIssuer {
					return nil, fmt.Errorf("failed to verify `iss` claim")
				}
			}

			return secretKey, nil
		})

		if token == nil {
			_, _ = fmt.Fprintf(w, err.Error())
		} else if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return true
		} else {
			_, _ = fmt.Fprintf(w, err.Error())
		}
	}

	return false
}
