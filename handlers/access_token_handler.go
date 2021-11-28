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
		w.Header().Set("Content-Type", "application/json")

		v := r.URL.Query()
		userID, err := strconv.ParseInt(v.Get("userId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "unable to get user identifier as `userId` from URL query: %s", err.Error())
			return
		}

		accessToken, refreshToken, exp, err := tokenService.GenerateUserTokenPair(uint(userID))
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = fmt.Fprintf(w, "unable to generate user token pair: %s", err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)

		data := &env.JsonJwt{
			AccessToken:    accessToken,
			RefreshToken:   refreshToken,
			ExpirationTime: strconv.FormatInt(exp, 10),
		}

		_ = json.NewEncoder(w).Encode(data)
	}
}
