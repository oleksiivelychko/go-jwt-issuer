package main

import (
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/handlers"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"log"
	"net/http"
)

func main() {
	cfg := env.InitConfig()
	tokenService := service.Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	http.HandleFunc("/", handlers.AccessTokenHandler(&tokenService))
	http.HandleFunc("/access-token", handlers.AccessTokenHandler(&tokenService))
	log.Fatal(http.ListenAndServe(env.GetPort(), nil))
}
