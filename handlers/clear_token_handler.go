package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"log"
	"net/http"
)

type ClearTokenHandler struct {
	tokenService *service.Service
}

func NewClearTokenHandler(tokenService *service.Service) *ClearTokenHandler {
	return &ClearTokenHandler{tokenService}
}

func (h *ClearTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := r.Context().Value(middleware.JWTClaimsCTX{}).(*issuer.JwtClaims)

	err := h.tokenService.ValidateCachedToken(claims, false)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, err.Error())
		log.Print(err.Error())
		return
	}

	err = h.tokenService.ClearCachedToken(claims)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = fmt.Fprintf(w, "unable to clear user token pair: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("token was successfully deleted")
}
