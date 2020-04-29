package main

import (
	"database/sql"
	"log"
	"net/http"

	"service.auth/controllers"
	"service.auth/driver"

	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	db = driver.ConnectDB()
	router := mux.NewRouter()
	controller := controllers.NewController()

	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")
	router.HandleFunc("/protected", controller.TokenVerifyMiddleWare(controller.ProtectedEndpoint())).Methods("GET")

	log.Println("Server listening on :8080...")
	log.Println("[POSTMAN] Hit http://127.0.0.1:8080/signup and singup with a new user")
	log.Println("[BROWSER] Hit http://127.0.0.1:8081 to login to adminer db console")
	log.Fatal(http.ListenAndServe(":8080", router))
}
