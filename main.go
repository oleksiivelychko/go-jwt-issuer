package main

import (
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/handlers"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"log"
	"net/http"
)

func main() {
	// env.InitEnv() // uncomment for local testing
	cfg := env.InitConfig()
	tokenService := service.Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	http.HandleFunc("/access-token/", handlers.AccessTokenHandler(&tokenService))
	http.HandleFunc("/refresh-token/", handlers.RefreshTokenHandler(&tokenService))
	http.HandleFunc("/clear-token/", handlers.ClearTokenHandler(&tokenService))
	log.Fatal(http.ListenAndServe(env.GetPort(), nil))
}
