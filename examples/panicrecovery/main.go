package main

import (
	"log"
	"net/http"

	"github.com/yashpokar/router"
)

func PanicMaker(w http.ResponseWriter, r *http.Request) {
	panic("simulating panic in handler")
}

func main() {
	r := router.New()
	r.OnPanic(func(w http.ResponseWriter, _ *http.Request, r any) {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(r.(string))); err != nil {
			log.Fatal(err)
		}
	})

	r.GET("/", PanicMaker)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
