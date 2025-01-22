package app

import (
	"database/sql"
	"log"
	controller "main/Controller"
	"main/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes(db *sql.DB) {

	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/users", controller.GetUsers(db)).Methods("GET")
	router.HandleFunc("/users/{id}", controller.GetUser(db)).Methods("GET")
	router.HandleFunc("/users", controller.CreateUser(db)).Methods("POST")
	router.HandleFunc("/users/{id}", controller.UpdateUser(db)).Methods("PUT")
	router.HandleFunc("/users/{id}", controller.DeleteUser(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", utils.JsonContentTypeMiddleware(router)))
}
