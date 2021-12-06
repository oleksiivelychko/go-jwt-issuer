package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"log"
	"net/http"
)

func ClearTokenHandler(tokenService *service.Service) func(w http.ResponseWriter, r *http.Request) {
	if tokenService.Redis == nil {
		log.Fatal("cannot established redis connection")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		w.Header().Set("Content-Type", "application/json")

		claims, _, err := middleware.ValidateRequest(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, err)
			log.Print(err)
			return
		}

		err = tokenService.ValidateCachedToken(claims, false)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprint(w, err.Error())
			log.Print(err.Error())
			return
		}

		err = tokenService.ClearCachedToken(claims)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = fmt.Fprintf(w, "unable to clear user token pair: %s", err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("token was successfully deleted")
	}
}
