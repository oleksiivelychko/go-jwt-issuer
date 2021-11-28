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
		w.Header().Set("Content-Type", "application/json")

		validated := middleware.ValidateRequest(w, r)
		if !validated {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintf(w, "\nunable to validate refresh token")
			log.Printf("\nunable to validate refresh token")
			return
		}

		token := r.Header.Get("Authorization")
		exp, err := strconv.ParseInt(r.Header.Get("ExpirationTime"), 10, 64)
		if err != nil {
			exp = 0
		}

		claims, err := tokenService.ValidateParsedToken(token, exp)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintf(w, "unable to validate parsed refresh token: %s", err.Error())
			log.Printf("unable to validate parsed refresh token: %s", err.Error())
			return
		}

		err = tokenService.ValidateCachedToken(claims, true)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprint(w, err.Error())
			log.Print(err.Error())
			return
		}

		v := r.URL.Query()
		userID, err := strconv.ParseInt(v.Get("userId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "unable to get user identifier as `userId` from URL query: %s", err.Error())
			return
		}

		accessToken, refreshToken, exp, err := tokenService.GenerateUserTokenPair(uint(userID))
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
