package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Land0fChocolate/FaceIt/database"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting FaceIt Users.")

	//starting up fake microservices to run alongside the Users microservice
	go searchService()
	go competitionService()

	usersService()
}

func usersService() {
	database.InitUsers()
	log.Println("Initialised mock database.")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", database.HomeLink)
	router.HandleFunc("/user", database.CreateUser).Methods("POST")
	router.HandleFunc("/users", database.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", database.GetOneUser).Methods("GET")
	router.HandleFunc("/users/country/{country}", database.GetUserByCountry).Methods("GET")
	router.HandleFunc("/users/{id}", database.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", database.DeleteUser).Methods("DELETE")
	router.HandleFunc("/health", database.HealthCheck).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

//Fake Search service, uses http://localhost:8081
func searchService() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/user-create", getMessage).Methods("POST")
	router.HandleFunc("/user-update", getMessage).Methods("POST")
	router.HandleFunc("/user-delete", getMessage).Methods("POST")
	router.HandleFunc("/health", database.HealthCheck).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}

//Fake Competition service, uses http://localhost:8082
func competitionService() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/user-update", getMessage).Methods("POST")
	router.HandleFunc("/user-delete", getMessage).Methods("POST")
	router.HandleFunc("/health", database.HealthCheck).Methods("POST")
	log.Fatal(http.ListenAndServe(":8082", router))
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	//Service gets message here and does what it needs to do with it. For this task it will just print some text.
	fmt.Fprintf(w, "Message received.")
}
