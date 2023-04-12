package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"github.com/oleksiivelychko/go-jwt-issuer/token"
	"log"
	"net/http"
	"strconv"
)

type AccessToken struct {
	tokenService *token.Service
}

func NewAccessToken(tokenService *token.Service) *AccessToken {
	return &AccessToken{tokenService}
}

func (handler *AccessToken) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Print("handler AccessToken is served")
	resp.Header().Set("Content-Type", "application/json")

	queryValues := req.URL.Query()
	userID, err := strconv.ParseInt(queryValues.Get("userID"), 10, 64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(resp, "unable to get userID from URL: %s", err.Error())
		return
	}

	accessToken, refreshToken, exp, err := handler.tokenService.GeneratePairTokens(uint(userID))
	if err != nil {
		resp.WriteHeader(http.StatusServiceUnavailable)
		_, _ = fmt.Fprintf(resp, "unable to generate pair of tokens: %s", err.Error())
		return
	}

	resp.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(resp).Encode(&issuer.ResponseJWT{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		ExpirationTime: strconv.FormatInt(exp, 10),
	})
}
