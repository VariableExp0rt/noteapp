package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"log"
	"time"

	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	SignKey   *rsa.PrivateKey
	VerifyKey *rsa.PublicKey
)

type UserInfo struct {
	Username string
}

type CustomClaimExample struct {
	*jwt.StandardClaims
	TokenType string
	UserInfo
}

func init() {
	err := MakeRSAPrivateKey()
	if err != nil {
		log.Fatal("Unable to generate signing key pair.")
	}
}

func MakeUserLoginEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u UserLoginReq
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid user info.", http.StatusBadRequest)
			return
		}

		err := s.ValidateUser(u.Username, u.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		//create token
		token, err := createToken(u.Username)
		if err != nil || token == "" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//return token to user
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))

	}
}

func MakeUserSignUpEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//Store credentials
		var usr UserSignUpReq
		if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
			http.Error(w, "Invalid sign-up request parameters.", http.StatusBadRequest)
			return
		}

		if err := s.StoreNewUser(usr); err != nil {
			http.Error(w, "Error signing up user."+err.Error(), http.StatusInternalServerError)
			return
		}

		//callback to an account page?
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully signed up new user."))

	}
}

func createToken(user string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	t.Claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		Id:        user,
	}

	return t.SignedString(SignKey)
}

func MakeRSAPrivateKey() error {
	var err error
	SignKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	VerifyKey = &SignKey.PublicKey
	return nil
}
