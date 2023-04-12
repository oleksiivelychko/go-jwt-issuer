package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/token"
	"log"
	"net/http"
)

type AuthorizeToken struct {
	tokenService *token.Service
}

func NewAuthorizeToken(tokenService *token.Service) *AuthorizeToken {
	return &AuthorizeToken{tokenService}
}

func (handler *AuthorizeToken) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Print("handler AuthorizeToken is served")
	claims := req.Context().Value(middleware.ContextClaimsJWT{}).(*issuer.ClaimsJWT)

	err := handler.tokenService.ValidateCachedToken(claims, false)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(resp).Encode(&issuer.ResponseJWT{ErrorMessage: err.Error()})
		return
	}

	resp.WriteHeader(http.StatusOK)
}
