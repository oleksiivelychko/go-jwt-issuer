package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/token"
	"log"
	"net/http"
)

type ClearToken struct {
	tokenService *token.Service
}

func NewClearToken(tokenService *token.Service) *ClearToken {
	return &ClearToken{tokenService}
}

func (handler *ClearToken) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Print("handler ClearToken is served")
	resp.Header().Set("Content-Type", "application/json")

	claims := req.Context().Value(middleware.ContextClaimsJWT{}).(*issuer.ClaimsJWT)

	err := handler.tokenService.ValidateCachedToken(claims, false)
	if err != nil {
		resp.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(resp, err.Error())
		log.Print(err.Error())
		return
	}

	err = handler.tokenService.ClearCachedTokens(claims)
	if err != nil {
		resp.WriteHeader(http.StatusServiceUnavailable)
		_, _ = fmt.Fprintf(resp, "unable to clear pair of tokens: %s", err.Error())
		return
	}

	resp.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(resp).Encode("pair of tokens were successfully cleared")
}
