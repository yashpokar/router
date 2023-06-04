package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/yashpokar/router"
)

type contextKey string

var (
	tokenKey = contextKey("token")
)

func RequestLoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		bytes, err := io.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("request body is %s", bytes)

		next.ServeHTTP(writer, request)
	}
}

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authorizationHeader := request.Header.Get("Authorization")
		// if the authorization token is not present
		// then send unauthorized status code
		if authorizationHeader == "" {
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You are not authorized to perform this action"))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		// extract the bearer token
		token := strings.TrimRight(authorizationHeader, "Bearer ")

		// get the existing context
		ctx := request.Context()
		// add token to the existing context
		ctx = context.WithValue(ctx, tokenKey, token)
		// update the request context so that it can be
		// used by the handler or next middleware if any
		request = request.WithContext(ctx)

		// call the next middleware / handler
		next.ServeHTTP(writer, request)
	}
}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// You have your token at handle now
	// do whatever you want to do with it
	token := ctx.Value(tokenKey)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Post created successfully",
		"token":   token,
	}); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := router.New()
	r.POST("/posts", RequestLoggerMiddleware(AuthenticationMiddleware(CreatePostsHandler)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
