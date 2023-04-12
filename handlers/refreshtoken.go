package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/token"
	"log"
	"net/http"
	"strconv"
)

type RefreshToken struct {
	tokenService *token.Service
}

func NewRefreshToken(tokenService *token.Service) *RefreshToken {
	return &RefreshToken{tokenService}
}

func (handler *RefreshToken) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Print("handler RefreshToken is served")
	resp.Header().Set("Content-Type", "application/json")

	claims := req.Context().Value(middleware.ContextClaimsJWT{}).(*issuer.ClaimsJWT)

	err := handler.tokenService.ValidateCachedToken(claims, true)
	if err != nil {
		resp.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(resp).Encode(&issuer.ResponseJWT{ErrorMessage: err.Error()})
		return
	}

	accessToken, refreshToken, exp, err := handler.tokenService.GeneratePairTokens(claims.ID)
	if err != nil {
		resp.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(resp).Encode(&issuer.ResponseJWT{ErrorMessage: err.Error()})
		return
	}

	resp.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(resp).Encode(&issuer.ResponseJWT{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		ExpirationTime: strconv.FormatInt(exp, 10),
	})
}
