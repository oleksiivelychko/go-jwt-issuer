package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/handlers"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
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

	cmd := tokenService.Redis.Echo(context.Background(), "check")
	if cmd.Err() != nil {
		log.Fatalf("cannot established redis connection: %s", cmd.Err())
	}

	serveMux := mux.NewRouter()

	accessTokenHandler := handlers.NewAccessTokenHandler(&tokenService)
	refreshTokenHandler := handlers.NewRefreshTokenHandler(&tokenService)
	clearTokenHandler := handlers.NewClearTokenHandler(&tokenService)
	authorizeTokenHandler := handlers.NewAuthorizeTokenHandler(&tokenService)

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.Use(middleware.JWT)

	getRouter.HandleFunc("/access-token/", accessTokenHandler.ServeHTTP)
	postRouter.HandleFunc("/access-token/", refreshTokenHandler.ServeHTTP)
	postRouter.HandleFunc("/clear-token/", clearTokenHandler.ServeHTTP)
	postRouter.HandleFunc("/authorize-token/", authorizeTokenHandler.ServeHTTP)

	server := &http.Server{
		Addr:         env.GetPort(),
		Handler:      serveMux,
		IdleTimeout:  2 * time.Minute,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}

	// server is run in a separate routine for each request
	go func() {
		log.Printf("Starting server on port %s", env.GetPort())
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
