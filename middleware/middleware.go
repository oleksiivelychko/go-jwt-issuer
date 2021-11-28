package middleware

import (
	"fmt"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"net/http"
	"strconv"
)

func AllowToEndpoint(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ValidateRequest(w, r) {
			endpoint(w, r)
		}
	})
}

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ValidateRequest(w, r) {
			next.ServeHTTP(w, r)
		}
	})
}

func ValidateRequest(w http.ResponseWriter, r *http.Request) bool {
	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()

	tokenHeader := r.Header.Get("Authorization")
	exp, err := strconv.ParseInt(r.Header.Get("Expires"), 10, 64)
	if err != nil {
		exp = 0
	}

	if len(secretKey) == 0 {
		_, _ = fmt.Fprintf(w, "environment variable `SECRET_KEY` is not defined")
	} else if tokenHeader == "" {
		_, _ = fmt.Fprintf(w, "failed to get token from header request")
	} else {
		token, err := issuer.ValidateToken(tokenHeader, secretKey, aud, iss, exp)
		if token == nil {
			_, _ = fmt.Fprintf(w, err.Error())
		} else if _, ok := token.Claims.(*issuer.JwtClaims); ok && token.Valid {
			return true
		} else {
			_, _ = fmt.Fprintf(w, err.Error())
		}
	}

	return false
}
