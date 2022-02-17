package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"net/http"
)

type AuthorizeTokenHandler struct {
	tokenService *service.Service
}

func NewAuthorizeTokenHandler(tokenService *service.Service) *AuthorizeTokenHandler {
	return &AuthorizeTokenHandler{tokenService}
}

func (h *AuthorizeTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.JWTClaimsCTX{}).(*issuer.JwtClaims)

	err := h.tokenService.ValidateCachedToken(claims, false)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(&env.JsonJwt{ErrorMessage: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
