package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"log"
	"net/http"
)

func AuthorizeTokenHandler(tokenService *service.Service) func(w http.ResponseWriter, r *http.Request) {
	if tokenService.Redis == nil {
		log.Fatal("cannot established redis connection")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		claims, errorCode, err := middleware.ValidateRequest(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(&env.JsonJwt{ErrorMessage: err.Error(), ErrorCode: errorCode})
			return
		}

		err = tokenService.ValidateCachedToken(claims, false)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(&env.JsonJwt{ErrorMessage: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
