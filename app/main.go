package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostRequest struct {
	Content string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&Post{}) // postsテーブルの作成

	router := mux.NewRouter()

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	router.HandleFunc("/api/v1/posts/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		// GETリクエストのみ受け付ける
		if r.Method != "GET" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]
		post := Post{}
		if err := db.First(&post, "id=?", id).Error; gorm.IsRecordNotFoundError(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// header → status code → response body の順番にしないと無効になる
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(post)
		w.Write(res)
	})

	router.HandleFunc("/api/v1/posts", func(w http.ResponseWriter, r *http.Request) {
		// POSTリクエストのみ受け付ける
		if r.Method != "POST" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// application/jsonのみ受け付ける
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// リクエストボディをpostRequestに変換する
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		var postRequest PostRequest
		json.Unmarshal(body, &postRequest)

		post := Post{}
		post.Content = postRequest.Content
		db.Create(&post)
		w.WriteHeader(http.StatusNoContent)
	})

	http.ListenAndServe(":8080", router)
}
