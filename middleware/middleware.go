package middleware

import (
	"fmt"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"net/http"
	"strconv"
	"strings"
)

func AllowToEndpoint(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := ValidateRequest(w, r)
		if err == nil {
			endpoint(w, r)
		} else {
			_, _ = w.Write([]byte(err.Error()))
		}
	})
}

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := ValidateRequest(w, r)
		if err == nil {
			next.ServeHTTP(w, r)
		} else {
			_, _ = w.Write([]byte(err.Error()))
		}
	})
}

func ValidateRequest(w http.ResponseWriter, r *http.Request) (*issuer.JwtClaims, uint8, error) {
	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()

	tokenHeader := r.Header.Get("Authorization")
	exp, err := strconv.ParseInt(r.Header.Get("Expires"), 10, 64)
	if err != nil {
		exp = 0
	}

	var status uint8 = 0

	if len(secretKey) == 0 {
		return nil, issuer.EnvironmentVariableSecretKeyIsNotDefined, fmt.Errorf("environment variable `SECRET_KEY` is not defined")
	} else if tokenHeader == "" {
		return nil, issuer.FailedToGetTokenFromHeaderRequest, fmt.Errorf("failed to get token from header request")
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
			return nil, status, err
		} else if claims, ok := token.Claims.(*issuer.JwtClaims); ok && token.Valid {
			return claims, status, nil
		} else {
			return nil, status, err
		}
	}
}
