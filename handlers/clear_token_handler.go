package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"log"
	"net/http"
	"strconv"
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

		validated := middleware.ValidateRequest(w, r)
		if !validated {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintf(w, "\nunable to validate refresh token")
			log.Printf("\nunable to validate refresh token")
			return
		}

		token := r.Header.Get("Authorization")
		exp, err := strconv.ParseInt(r.Header.Get("Expires"), 10, 64)
		if err != nil {
			exp = 0
		}

		claims, err := tokenService.ValidateParsedToken(token, exp)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintf(w, "unable to validate parsed refresh token: %s", err.Error())
			log.Printf("unable to validate parsed refresh token: %s", err.Error())
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
