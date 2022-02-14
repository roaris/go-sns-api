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
	likeHandler := handlers.NewLikeHandler(db)

	v1r := r.PathPrefix("/api/v1").Subrouter()

	authMiddleware := middlewares.AuthMiddleware
	v1r.Methods(http.MethodPost).Path("/posts").Handler(authMiddleware(AppHandler{postHandler.Create}))
	v1r.Methods(http.MethodGet).Path("/posts/{id:[0-9]+}").Handler(AppHandler{postHandler.Show})
	v1r.Methods(http.MethodGet).Path("/posts").Queries("limit", "{limit:[0-9]+}", "offset", "{offset:[0-9]+}").Handler(authMiddleware(AppHandler{postHandler.Index}))
	v1r.Methods(http.MethodPatch).Path("/posts/{id:[0-9]+}").Handler(authMiddleware(AppHandler{postHandler.Update}))
	v1r.Methods(http.MethodDelete).Path("/posts/{id:[0-9]+}").Handler(authMiddleware(AppHandler{postHandler.Destroy}))
	v1r.Methods(http.MethodPost).Path("/users").Handler(AppHandler{userHandler.Create})
	v1r.Methods(http.MethodGet).Path("/users/me").Handler(authMiddleware(AppHandler{userHandler.ShowMe}))
	v1r.Methods(http.MethodPatch).Path("/users/me").Handler(authMiddleware(AppHandler{userHandler.UpdateMe}))
	v1r.Methods(http.MethodPost).Path("/auth").Handler(AppHandler{authHandler.Authenticate})
	v1r.Methods(http.MethodPost).Path("/users/me/followees").Handler(authMiddleware(AppHandler{friendshipHandler.Create}))
	v1r.Methods(http.MethodGet).Path("/users/{id:[0-9]+}/followees").Handler(AppHandler{friendshipHandler.ShowFollowees})
	v1r.Methods(http.MethodGet).Path("/users/{id:[0-9]+}/followers").Handler(AppHandler{friendshipHandler.ShowFollowers})
	v1r.Methods(http.MethodDelete).Path("/users/me/followees/{id:[0-9]+}").Handler(authMiddleware(AppHandler{friendshipHandler.Destroy}))
	v1r.Methods(http.MethodPost).Path("/posts/{id:[0-9]+}/likes").Handler(authMiddleware(AppHandler{likeHandler.Create}))
	v1r.Methods(http.MethodDelete).Path("/posts/{id:[0-9]+}/likes").Handler(authMiddleware(AppHandler{likeHandler.Destroy}))
	http.ListenAndServe(":8080", c.Handler(r))
}
