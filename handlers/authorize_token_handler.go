package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"log"
	"net/http"
)

type AuthorizeTokenHandler struct {
	tokenService *service.Service
}

func NewAuthorizeTokenHandler(tokenService *service.Service) *AuthorizeTokenHandler {
	return &AuthorizeTokenHandler{tokenService}
}

func (h *AuthorizeTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.tokenService.Redis == nil {
		log.Fatal("cannot established redis connection")
	}

	if r.Method != "POST" {
		return
	}

	claims, errorCode, err := middleware.ValidateRequest(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(&env.JsonJwt{ErrorMessage: err.Error(), ErrorCode: errorCode})
		return
	}

	err = h.tokenService.ValidateCachedToken(claims, false)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(&env.JsonJwt{ErrorMessage: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
