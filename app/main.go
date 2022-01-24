package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

type Post struct {
	ID        int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&Post{}) // postsテーブルの作成

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	http.HandleFunc("/api/v1/posts", func(w http.ResponseWriter, r *http.Request) {
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
		err := db.Create(&post)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	http.ListenAndServe(":8080", nil)
}
