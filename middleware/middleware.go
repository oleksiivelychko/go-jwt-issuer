package middleware

import (
	"context"
	"github.com/oleksiivelychko/go-jwt-issuer/config"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"net/http"
	"strconv"
)

type ContextClaimsJWT struct{}

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		var secretKey = config.GetSecretKey()
		var audienceAUD = config.GetAudience()
		var issuerISS = config.GetIssuer()

		tokenHeader := request.Header.Get("Authorization")
		expiresIn, err := strconv.ParseInt(request.Header.Get("Expires"), 10, 64)
		if err != nil {
			expiresIn = 0
		}

		if secretKey == "" {
			http.Error(responseWriter, "environment variable 'SECRET_KEY' is not defined", http.StatusBadRequest)
			return
		}

		if tokenHeader == "" {
			http.Error(responseWriter, "unable to get token from 'Authorization' header", http.StatusBadRequest)
			return
		}

		token, tokenErr := issuer.ParseToken(tokenHeader, secretKey, audienceAUD, issuerISS, expiresIn)
		if tokenErr != nil {
			http.Error(responseWriter, tokenErr.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*issuer.ClaimsJWT); ok && token.Valid {
			contextWithClaimJWT := context.WithValue(request.Context(), ContextClaimsJWT{}, claims)
			request = request.WithContext(contextWithClaimJWT)
		} else {
			http.Error(responseWriter, "unable to validate token", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(responseWriter, request)
	})
}
