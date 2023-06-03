package main

import (
	"encoding/json"
	"github.com/yashpokar/router"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
}

var users []User

func Index(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("Welcome from basic handler")); err != nil {
		log.Fatal(err)
	}
}

func ListUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := users
	if response == nil {
		response = []User{}
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var user User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(r.Body)

	if err = json.Unmarshal(body, &user); err != nil {
		log.Fatal(err)
	}

	user.ID = strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	users = append(users, user)

	if err = json.NewEncoder(w).Encode(&user); err != nil {
		log.Fatal(err)
	}
}

func UserDetails(w http.ResponseWriter, r *http.Request) {
	vars := router.Vars(r)
	userID := vars["user_id"]

	for _, user := range users {
		if user.ID == userID {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(user); err != nil {
				log.Fatal(err)
			}

			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	if _, err := w.Write([]byte("user not found.")); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := router.New()
	r.GET("/", Index)

	api := r.Group("/api/v1")

	api.GET("/users", ListUsers)
	api.POST("/users", CreateUser)
	api.GET("/users/:user_id", UserDetails)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
