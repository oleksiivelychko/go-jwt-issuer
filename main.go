package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/oleksiivelychko/go-jwt-issuer/config"
	"github.com/oleksiivelychko/go-jwt-issuer/handlers"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/token"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	tokenService := token.NewService(
		config.NewConfig(
			os.Getenv("SECRET_KEY"),
			os.Getenv("AUDIENCE_AUD"),
			os.Getenv("ISSUER_ISS"),
			os.Getenv("EXPIRATION_TIME_EXP"),
		), config.NewRedisClient(
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
			os.Getenv("REDIS_PASSWORD"),
		))

	cmd := tokenService.RedisClient.Ping(context.Background())
	if cmd.Err() != nil {
		log.Fatal(cmd.Err())
	} else {
		log.Printf("established connection: %s", tokenService.RedisClient)
	}

	muxRouter := mux.NewRouter()

	accessTokenHandler := handlers.NewAccessToken(tokenService)
	refreshTokenHandler := handlers.NewRefreshToken(tokenService)
	clearTokenHandler := handlers.NewClearToken(tokenService)
	authorizeTokenHandler := handlers.NewAuthorizeToken(tokenService)

	getRouter := muxRouter.Methods(http.MethodGet).Subrouter()
	postRouter := muxRouter.Methods(http.MethodPost).Subrouter()
	postRouter.Use(middleware.JWT)

	getRouter.HandleFunc("/access-token", accessTokenHandler.ServeHTTP)
	postRouter.HandleFunc("/refresh-token", refreshTokenHandler.ServeHTTP)
	postRouter.HandleFunc("/clear-token", clearTokenHandler.ServeHTTP)
	postRouter.HandleFunc("/authorize-token", authorizeTokenHandler.ServeHTTP)

	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	server := &http.Server{
		Addr:         addr,
		Handler:      muxRouter,
		IdleTimeout:  2 * time.Minute,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}

	// server is being run in a separate goroutine for each request
	go func() {
		log.Printf("starting server on %s", addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt)
	signal.Notify(signalCh, os.Kill)

	sig := <-signalCh
	log.Println("received terminate, graceful shutdown", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_ = server.Shutdown(ctx)
}
