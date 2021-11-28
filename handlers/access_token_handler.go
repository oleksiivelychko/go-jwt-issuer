package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"log"
	"net/http"
	"strconv"
)

func AccessTokenHandler(tokenService *service.Service) func(w http.ResponseWriter, r *http.Request) {
	if tokenService.Redis == nil {
		log.Fatal("cannot established redis connection")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, refreshToken, exp, err := tokenService.GenerateUserTokenPair(1)
		if err != nil {
			_, _ = fmt.Fprintf(w, "failed to get access token: %s", err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		data := &env.JsonJwt{
			AccessToken:    accessToken,
			RefreshToken:   refreshToken,
			ExpirationTime: strconv.FormatInt(exp, 10),
		}

		_ = json.NewEncoder(w).Encode(data)
	}
}
