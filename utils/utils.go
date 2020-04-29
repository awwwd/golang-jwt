package utils

import (
	"encoding/json"
	"log"

	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"service.auth/models"
)

// ResponseWithError ...
func ResponseWithError(w http.ResponseWriter, status int, err models.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

// ResponseWithJSON ...
func ResponseWithJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

// GenerateToken helps generating token
func GenerateToken(user models.User) (string, error) {
	var err error
	secret := os.Getenv("SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "course",
		"email": user.Email,
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString, nil
}
