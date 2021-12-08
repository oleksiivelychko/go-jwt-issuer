package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"log"
	"net/http"
	"strconv"
)

func RefreshTokenHandler(tokenService *service.Service) func(w http.ResponseWriter, r *http.Request) {
	if tokenService.Redis == nil {
		log.Fatal("cannot established redis connection")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		w.Header().Set("Content-Type", "application/json")

		claims, errorCode, err := middleware.ValidateRequest(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(&env.JsonJwt{ErrorMessage: err.Error(), ErrorCode: errorCode})
			return
		}

		err = tokenService.ValidateCachedToken(claims, true)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(&env.JsonJwt{ErrorMessage: err.Error()})
			return
		}

		accessToken, refreshToken, exp, err := tokenService.GenerateUserTokenPair(claims.ID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(&env.JsonJwt{ErrorMessage: err.Error()})
			return
		}

		data := &env.JsonJwt{
			AccessToken:    accessToken,
			RefreshToken:   refreshToken,
			ExpirationTime: strconv.FormatInt(exp, 10),
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(data)
	}
}
