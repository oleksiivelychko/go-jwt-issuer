package middleware

import (
	"context"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"net/http"
	"strconv"
	"strings"
)

type JWTClaimsCTX struct{}

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var secretKey = env.GetSecretKey()
		var aud = env.GetAUD()
		var iss = env.GetISS()
		var status = 400

		tokenHeader := r.Header.Get("Authorization")
		exp, err := strconv.ParseInt(r.Header.Get("Expires"), 10, 64)
		if err != nil {
			exp = 0
		}

		if len(secretKey) == 0 {
			http.Error(rw, "environment variable `SECRET_KEY` is not defined", http.StatusBadRequest)
			return
		} else if tokenHeader == "" {
			http.Error(rw, "failed to get token from header request", http.StatusBadRequest)
			return
		} else {
			token, err := issuer.ValidateToken(tokenHeader, secretKey, aud, iss, exp)

			if err != nil {
				if strings.Contains(err.Error(), "unexpected signing method") {
					status = issuer.UnexpectedSigningMethod
				}
				if strings.Contains(err.Error(), "failed to verify `aud` claim") {
					status = issuer.FailedAudienceClaim
				}
				if strings.Contains(err.Error(), "failed to verify `iss` claim") {
					status = issuer.FailedIssuerClaim
				}
				if strings.Contains(err.Error(), "token is expired by") {
					status = issuer.FailedExpirationTimeClaim
				}
			}

			if token == nil {
				http.Error(rw, err.Error(), status)
				return
			} else if claims, ok := token.Claims.(*issuer.JwtClaims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), JWTClaimsCTX{}, claims)
				r = r.WithContext(ctx)
			} else {
				http.Error(rw, err.Error(), http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(rw, r)
	})
}
