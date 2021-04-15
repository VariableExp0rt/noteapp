package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/VariableExp0rt/dddexample/auth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
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
			http.Error(w, request.ErrNoTokenInRequest.Error(), http.StatusBadRequest)
		}

		if ok := token.Valid; !ok {
			http.Error(w, "Invalid bearer token, please reauthenticate.", http.StatusBadRequest)
		}

		next(w, r)
	}
}
