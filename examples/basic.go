package main

import (
	"fmt"
	"github.com/yashpokar/router"
	"log"
	"net/http"
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
	_, _ = fmt.Fprintf(w, "Welcome from basic handler")
}

func ListUsers(w http.ResponseWriter, _ *http.Request) {

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

}

func UserDetails(w http.ResponseWriter, r *http.Request) {

}

func BasicServer() {
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
