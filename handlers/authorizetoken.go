package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"net/http"
)

type AuthorizeToken struct {
	tokenService *service.TokenService
}

func NewAuthorizeToken(tokenService *service.TokenService) *AuthorizeToken {
	return &AuthorizeToken{tokenService}
}

func (handler *AuthorizeToken) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	claims := req.Context().Value(middleware.ContextClaimsJWT{}).(*issuer.ClaimsJWT)

	err := handler.tokenService.ValidateCachedToken(claims, false)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(resp).Encode(&issuer.ResponseJWT{ErrorMessage: err.Error()})
		return
	}

	resp.WriteHeader(http.StatusOK)
}
