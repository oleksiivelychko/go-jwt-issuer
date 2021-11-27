package main

import (
	"fmt"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"log"
	"net/http"
)

func issueAccessTokenHandler(tokenService *service.Service) func(w http.ResponseWriter, r *http.Request) {
	if tokenService.Redis == nil {
		log.Fatal("cannot established redis connection")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, refreshToken, exp, err := tokenService.GenerateUserTokenPair(1)
		if err != nil {
			_, _ = fmt.Fprintf(w, "failed to get access token: %s", err.Error())
		}
		_, _ = fmt.Fprintf(w, "access-token: %s \n", accessToken)
		_, _ = fmt.Fprintf(w, "refresh-token: %s \n", refreshToken)
		_, _ = fmt.Fprintf(w, "expiration-time: %d \n", exp)
	}
}

func main() {
	cfg := env.InitConfig()
	tokenService := service.Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	http.HandleFunc("/", issueAccessTokenHandler(&tokenService))
	http.HandleFunc("/access-token", issueAccessTokenHandler(&tokenService))
	log.Fatal(http.ListenAndServe(env.GetPort(), nil))
}
