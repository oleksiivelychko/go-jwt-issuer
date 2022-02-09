package main

import (
	"context"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/handlers"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	//env.InitEnv() // uncomment for local testing
	cfg := env.InitConfig()
	tokenService := service.Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	serveMux := http.NewServeMux()
	serveMux.Handle("/access-token/", handlers.NewAccessTokenHandler(&tokenService))
	serveMux.Handle("/refresh-token/", handlers.NewRefreshTokenHandler(&tokenService))
	serveMux.Handle("/clear-token/", handlers.NewClearTokenHandler(&tokenService))
	serveMux.Handle("/authorize-token/", handlers.NewAuthorizeTokenHandler(&tokenService))

	server := &http.Server{
		Addr:         env.GetPort(),
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	log.Println("Received terminate, graceful shutdown", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
