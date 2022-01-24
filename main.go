package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/roaris/go_sns_api/handlers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	router.HandleFunc("/api/v1/posts/{id:[0-9]+}", handlers.PostShow)
	router.HandleFunc("/api/v1/posts", handlers.PostCreate)
	http.ListenAndServe(":8080", router)
}
