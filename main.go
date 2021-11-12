package main

import (
	"fmt"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"log"
	"net/http"
)

func jwtIssuer(w http.ResponseWriter, r *http.Request) {
	var secretKey = env.GetSecretKey()
	if len(secretKey) > 0 {
		var aud = env.GetAUD()
		var iss = env.GetISS()

		validToken, err := issuer.IssueJWT(secretKey, aud, iss)
		if err != nil {
			_, _ = fmt.Fprintf(w, "failed to get the complete signed token: %s", err.Error())
		}

		_, _ = fmt.Fprintf(w, validToken)
	} else {
		_, _ = fmt.Fprintf(w, "environment variable `SECRET_KEY` is not defined")
	}
}

func main() {
	http.HandleFunc("/", jwtIssuer)
	http.HandleFunc("/issue", jwtIssuer)
	log.Fatal(http.ListenAndServe(env.GetPort(), nil))
}
