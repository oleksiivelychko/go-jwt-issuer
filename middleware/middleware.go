package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"net/http"
	"time"
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
	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()
	var exp = time.Now().Add(time.Minute * time.Duration(expiresMinutes)).Unix()

	tokenHeader := r.Header.Get("Authorization")

	if len(secretKey) == 0 {
		_, _ = fmt.Fprintf(w, "environment variable `SECRET_KEY` is not defined")
	} else if tokenHeader == "" {
		_, _ = fmt.Fprintf(w, "failed to get token from header request")
	} else {
		token, err := issuer.ValidateToken(tokenHeader, secretKey, aud, iss, exp)
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
