package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/roaris/go_sns_api/models"
)

type PostRequest struct {
	Content string
}

func PostShow(w http.ResponseWriter, r *http.Request) {
	// GETリクエストのみ受け付ける
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// パスパラメータの取得
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	post, err := models.ShowPost(id)
	if gorm.IsRecordNotFoundError(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// header → status code → response body の順番にしないと無効になる
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(post)
	w.Write(res)
}

func PostCreate(w http.ResponseWriter, r *http.Request) {
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

	// リクエストボディをPostRequestに変換する
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var postRequest PostRequest
	json.Unmarshal(body, &postRequest)

	models.CreatePost(postRequest.Content)
	w.WriteHeader(http.StatusNoContent)
}
