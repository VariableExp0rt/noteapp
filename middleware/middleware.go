package middleware

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/VariableExp0rt/dddexample/auth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

var (
	ErrInvalidSigningAlg = errors.New("Invalid signing algorithm for JWT.")
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//do things
		log.Printf("[%v] - %v %v %v", time.Now().UTC(), r.Method, r.Header, r.Host)
		next(w, r)
	})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//auth logic, validate JWT created via the login endpoint, proceed if token already exists
		//check token map, if exists, forward to target (e.g. /notes) else
		//return error, must login to proceed
		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) { return auth.VerifyKey, nil })
		if err != nil {
			http.Error(w, request.ErrNoTokenInRequest.Error()+" Please authenticate to obtain token.", http.StatusUnauthorized)
			return
		}

		//TODO: Need to save the metadata of the JWT on login, and verify it here by retrieving the K/V pairs
		//this needs to be done in boltstorage and a helper function in auth service

		if token.Method.Alg() != "RS256" {
			http.Error(w, ErrInvalidSigningAlg.Error(), http.StatusBadRequest)
			return
		}

		if ok := token.Valid; !ok {
			http.Error(w, "Invalid bearer token.", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
