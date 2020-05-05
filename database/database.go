//This is a mock database using a REST api
//http://localhost:8080 is used as the path to communicate with the mock database. An application like Postman can be used to give a body to requests sent.

package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type user struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"` //for this task, password will be an unsecure string.
	Email     string `json:"email"`
	Country   string `json:"country"`
}

var users []user

func InitUsers() {
	users = []user{
		{
			Id:        "1234",
			FirstName: "Tom",
			LastName:  "Beckingham",
			Nickname:  "qwertyuiop",
			Password:  "someTextHere12",
			Email:     "tebeckingham@gmail.com",
			Country:   "United Kingdom",
		},
		{
			Id:        "2345",
			FirstName: "Bob",
			LastName:  "Bobson",
			Nickname:  "asdghkl",
			Password:  "someTextHere34",
			Email:     "bobbob@gmail.com",
			Country:   "Italy",
		},
		{
			Id:        "3456",
			FirstName: "Doug",
			LastName:  "Dimmadome",
			Nickname:  "DimmaDude",
			Password:  "someTextHere56",
			Email:     "dimmadome@gmail.com",
			Country:   "Sealand",
		},
	}
}

func HomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new user")

	var newUser user
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
	}

	json.Unmarshal(reqBody, &newUser)
	users = append(users, newUser)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUser)

	log.Println("Successfully created new user", newUser.Id)

	//notify Search service for change
	_, err = http.Post("http://localhost:8081", "/user-create", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("Http create post to Search service has failed:", err)
	}
}

func GetOneUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	log.Println("Retrieving user", userId)

	for _, user := range users {
		if user.Id == userId {
			json.NewEncoder(w).Encode(user)
		}
	}

	log.Println("Successfully retrieved user", userId)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Retrieving all users")
	json.NewEncoder(w).Encode(users)

	log.Println("Successfully retrieved all users.")
}

func GetUserByCountry(w http.ResponseWriter, r *http.Request) {
	country := mux.Vars(r)["country"]
	log.Println("Retrieving all users with Country", country)
	var userList []user

	for _, user := range users {
		if user.Country == country {
			userList = append(userList, user)
		}
	}

	json.NewEncoder(w).Encode(userList)

	log.Println("Successfully retrieved all users with Country", country)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	log.Println("Updating user ", userId)
	var updatedUser user

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
	}
	json.Unmarshal(reqBody, &updatedUser)

	for i, user := range users {
		if user.Id == userId {
			user.FirstName = updatedUser.FirstName
			user.LastName = updatedUser.LastName
			user.Nickname = updatedUser.Nickname
			user.Password = updatedUser.Password
			user.Email = updatedUser.Email
			user.Country = updatedUser.Country
			users[i] = user
			json.NewEncoder(w).Encode(user)
		}
	}

	log.Printf("Successfully updated user %v with new data: %#v\n", userId, updatedUser)

	//notify Search and Competition services for change
	_, err = http.Post("http://localhost:8081", "/user-update", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("Http update post to Search service has failed:", err)
	}
	_, err = http.Post("http://localhost:8082", "/user-update", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("Http update post to Competition service has failed:", err)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	log.Println("Deleting user ", userId)

	for i, user := range users {
		if user.Id == userId {
			users = append(users[:i], users[i+1:]...)
		}
	}

	log.Println("Successfully deleted user ", userId)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
	}

	//notify Search and Competition services for change
	_, err = http.Post("http://localhost:8081", "/user-delete", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("Http delete post to Search service has failed:", err)
	}
	_, err = http.Post("http://localhost:8082", "/user-delete", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("Http delete post to Competition service has failed:", err)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(http.StatusOK)
	log.Println("Health check passed for", r.Host)
}
