package handlers

import (
	"encoding/json"
	"fmt"
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

		claims, _, err := middleware.ValidateRequest(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, err)
			log.Print(err)
			return
		}

		err = tokenService.ValidateCachedToken(claims, true)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprint(w, err.Error())
			log.Print(err.Error())
			return
		}

		accessToken, refreshToken, exp, err := tokenService.GenerateUserTokenPair(claims.ID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintf(w, "unable to regenerate user token pair: %s", err.Error())
			log.Printf("unable to regenerate user token pair: %s", err.Error())
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
