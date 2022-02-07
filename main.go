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
	v1r.Methods(http.MethodPost).Path("/posts").Handler(authMiddleware(AppHandler{handlers.CreatePost}))
	v1r.Methods(http.MethodGet).Path("/posts/{id:[0-9]+}").Handler(AppHandler{handlers.GetPost})
	v1r.Methods(http.MethodGet).Path("/posts").Handler(authMiddleware(AppHandler{handlers.GetPosts}))
	v1r.Methods(http.MethodPatch).Path("/posts/{id:[0-9]+}").Handler(authMiddleware(AppHandler{handlers.UpdatePost}))
	v1r.Methods(http.MethodDelete).Path("/posts/{id:[0-9]+}").Handler(authMiddleware(AppHandler{handlers.DeletePost}))
	v1r.Methods(http.MethodPost).Path("/users").Handler(AppHandler{handlers.CreateUser})
	v1r.Methods(http.MethodGet).Path("/users/me").Handler(authMiddleware(AppHandler{handlers.GetLoginUser}))
	v1r.Methods(http.MethodPatch).Path("/users/me").Handler(authMiddleware(AppHandler{handlers.UpdateLoginUser}))
	v1r.Methods(http.MethodPost).Path("/auth").Handler(AppHandler{handlers.Authenticate})
	v1r.Methods(http.MethodPost).Path("/users/me/followees").Handler(authMiddleware(AppHandler{handlers.CreateFollowee}))
	v1r.Methods(http.MethodGet).Path("/users/{id:[0-9]+}/followees").Handler(AppHandler{handlers.GetFollowees})
	v1r.Methods(http.MethodGet).Path("/users/{id:[0-9]+}/followers").Handler(AppHandler{handlers.GetFollowers})
	v1r.Methods(http.MethodDelete).Path("/users/me/followees/{id:[0-9]+}").Handler(authMiddleware(AppHandler{handlers.DeleteFollowee}))
	http.ListenAndServe(":8080", c.Handler(r))
}
