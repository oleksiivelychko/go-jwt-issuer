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

type RefreshTokenHandler struct {
	tokenService *service.Service
}

func NewRefreshTokenHandler(tokenService *service.Service) *RefreshTokenHandler {
	return &RefreshTokenHandler{tokenService}
}

func (h *RefreshTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.tokenService.Redis == nil {
		log.Fatal("cannot established redis connection")
	}

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

	err = h.tokenService.ValidateCachedToken(claims, true)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(&env.JsonJwt{ErrorMessage: err.Error()})
		return
	}

	accessToken, refreshToken, exp, err := h.tokenService.GenerateUserTokenPair(claims.ID)
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
