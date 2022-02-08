package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/roaris/go-sns-api/handlers"
	"github.com/roaris/go-sns-api/middlewares"
	"github.com/roaris/go-sns-api/models"
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

	db := models.CreateDB()
	authHandler := handlers.NewAuthHandler(db)
	friendshipHandler := handlers.NewFriendshipHandler(db)
	postHandler := handlers.NewPostHandler(db)
	userHandler := handlers.NewUserHandler(db)

	v1r := r.PathPrefix("/api/v1").Subrouter()

	authMiddleware := middlewares.AuthMiddleware
	v1r.Methods(http.MethodPost).Path("/posts").Handler(authMiddleware(AppHandler{postHandler.CreatePost}))
	v1r.Methods(http.MethodGet).Path("/posts/{id:[0-9]+}").Handler(AppHandler{postHandler.GetPost})
	v1r.Methods(http.MethodGet).Path("/posts").Queries("limit", "{limit:[0-9]+}", "offset", "{offset:[0-9]+}").Handler(authMiddleware(AppHandler{postHandler.GetPosts}))
	v1r.Methods(http.MethodPatch).Path("/posts/{id:[0-9]+}").Handler(authMiddleware(AppHandler{postHandler.UpdatePost}))
	v1r.Methods(http.MethodDelete).Path("/posts/{id:[0-9]+}").Handler(authMiddleware(AppHandler{postHandler.DeletePost}))
	v1r.Methods(http.MethodPost).Path("/users").Handler(AppHandler{userHandler.CreateUser})
	v1r.Methods(http.MethodGet).Path("/users/me").Handler(authMiddleware(AppHandler{userHandler.GetLoginUser}))
	v1r.Methods(http.MethodPatch).Path("/users/me").Handler(authMiddleware(AppHandler{userHandler.UpdateLoginUser}))
	v1r.Methods(http.MethodPost).Path("/auth").Handler(AppHandler{authHandler.Authenticate})
	v1r.Methods(http.MethodPost).Path("/users/me/followees").Handler(authMiddleware(AppHandler{friendshipHandler.CreateFollowee}))
	v1r.Methods(http.MethodGet).Path("/users/{id:[0-9]+}/followees").Handler(AppHandler{friendshipHandler.GetFollowees})
	v1r.Methods(http.MethodGet).Path("/users/{id:[0-9]+}/followers").Handler(AppHandler{friendshipHandler.GetFollowers})
	v1r.Methods(http.MethodDelete).Path("/users/me/followees/{id:[0-9]+}").Handler(authMiddleware(AppHandler{friendshipHandler.DeleteFollowee}))
	http.ListenAndServe(":8080", c.Handler(r))
}
