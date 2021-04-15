package auth

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	TokenMap  = make(map[string]string)
	mu        = &sync.Mutex{}
	Privkey   string
	Pubkey    string
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

func MakeUserLoginEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid user info.", http.StatusBadRequest)
		}

		err := s.ValidateUser(u.Username, u.Password)
		if err != nil {
			http.Error(w, err.Error()+" Sign-up or enter valid credentials.", http.StatusUnauthorized)
		}

		//create token
		token, err := createToken(u.Username)
		if err != nil {
			http.Error(w, "Error in auth flow: "+err.Error(), http.StatusInternalServerError)
		}

		//return token to user
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))

	}
}

func createToken(user string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RSA256"))

	t.Claims = &CustomClaimExample{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
		"test",
		UserInfo{Username: user},
	}

	return t.SignedString(SignKey)
}

func init() {
	signBytes, err := ioutil.ReadFile(Privkey)
	if err != nil {
		log.Fatal(err)
	}

	SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}

	verifyBytes, err := ioutil.ReadFile(Pubkey)
	if err != nil {
		log.Fatal(err)
	}

	VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err)
	}

}
