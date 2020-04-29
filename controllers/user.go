package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	userrepository "service.auth/repository/user"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"service.auth/models"
	"service.auth/utils"
)

func NewController() Controller {
	return Controller{}
}

func (c Controller) Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error
		userRepo := userrepository.UserRepository{}

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" || user.Password == "" {
			error.Message = "Wrong input. Either email or password is missing."
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(error)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Fatal(err)
		}

		user.Password = string(hash)
		user, err = userRepo.SaveUser(db, user)

		if err != nil {
			error.Message = "server error"
			utils.ResponseWithError(w, http.StatusInternalServerError, error)
			return
		}
		user.Password = "hidden"
		w.Header().Set("Content-Type", "application/json")
		utils.ResponseWithJSON(w, user)
	}
}

func (c Controller) Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var jwt models.JWT
		var error models.Error
		userRepo := userrepository.UserRepository{}

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" || user.Password == "" {
			error.Message = "Wrong input. Either email or password is missing."
			utils.ResponseWithError(w, http.StatusBadRequest, error)
			return
		}

		password := user.Password

		user, err := userRepo.FetchUser(db, user)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "User doesn't exists"
				utils.ResponseWithError(w, http.StatusBadRequest, error)
				return
			} else {
				log.Fatal(err)
			}
		}
		hashedPassword := user.Password
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			user.Password = "hidden"
			error.Message = "Invalid password"
			utils.ResponseWithError(w, http.StatusUnauthorized, error)
			return
		}
		token, err := utils.GenerateToken(user)
		w.WriteHeader(http.StatusOK)
		jwt.Token = token
		utils.ResponseWithJSON(w, jwt)
	}
}

func (c Controller) TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObject models.Error
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) == 2 && bearerToken[0] == "Bearer" {
			authToken := bearerToken[1]
			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte(os.Getenv("SECRET")), nil
			})
			if error != nil {
				errorObject.Message = error.Error()
				utils.ResponseWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				errorObject.Message = error.Error()
				utils.ResponseWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Message = "Invalid token"
			utils.ResponseWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}
