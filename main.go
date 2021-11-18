package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"log"
	"net/http"
	"os"
)

func jwtIssuer(w http.ResponseWriter, r *http.Request) {
	var secretKey = env.GetSecretKey()
	if len(secretKey) > 0 {
		var aud = env.GetAUD()
		var iss = env.GetISS()
		var expiresMinutes = env.GetExpiresMinutes()

		validToken, _, _, err := issuer.IssueJWT(secretKey, 0, expiresMinutes, aud, iss)
		if err != nil {
			_, _ = fmt.Fprintf(w, "failed to get the complete signed token: %s", err.Error())
		}

		_, _ = fmt.Fprintf(w, validToken)
	} else {
		_, _ = fmt.Fprintf(w, "environment variable `SECRET_KEY` is not defined")
	}
}

func initRedis() *redis.Client {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	_ = initRedis()

	http.HandleFunc("/", jwtIssuer)
	http.HandleFunc("/issue", jwtIssuer)
	log.Fatal(http.ListenAndServe(env.GetPort(), nil))
}
