package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/roaris/go-sns-api/handlers"
	"github.com/roaris/go-sns-api/middlewares"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	// CORSの設定
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodOptions,
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
	})

	v1r := r.PathPrefix("/api/v1").Subrouter()

	authMiddleware := middlewares.AuthMiddleware
	responseMiddleware := middlewares.ResponseMiddleware
	v1r.Methods(http.MethodPost).Path("/posts").HandlerFunc(authMiddleware(responseMiddleware(handlers.CreatePost)))
	v1r.Methods(http.MethodGet).Path("/posts/{id:[0-9]+}").HandlerFunc(responseMiddleware(handlers.GetPost))
	v1r.Methods(http.MethodPatch).Path("/posts/{id:[0-9]+}").HandlerFunc(authMiddleware(responseMiddleware(handlers.UpdatePost)))
	v1r.Methods(http.MethodDelete).Path("/posts/{id:[0-9]+}").HandlerFunc(authMiddleware(responseMiddleware(handlers.DeletePost)))
	v1r.Methods(http.MethodPost).Path("/users").HandlerFunc(responseMiddleware(handlers.CreateUser))
	v1r.Methods(http.MethodGet).Path("/users/me").HandlerFunc(authMiddleware(responseMiddleware(handlers.GetLoginUser)))
	v1r.Methods(http.MethodPatch).Path("/users/me").HandlerFunc(authMiddleware(responseMiddleware(handlers.UpdateLoginUser)))
	v1r.Methods(http.MethodPost).Path("/auth").HandlerFunc(responseMiddleware(handlers.Authenticate))
	http.ListenAndServe(":8080", c.Handler(r))
}
