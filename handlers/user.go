package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/roaris/go_sns_api/models"
	"gopkg.in/go-playground/validator.v9"
)

type UserRequest struct {
	Name     string
	Email    string
	Password string
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	// application/jsonのみ受け付ける
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// リクエストボディをUserRequestに変換する
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var userRequest UserRequest
	json.Unmarshal(body, &userRequest)

	err := models.CreateUser(userRequest.Name, userRequest.Email, userRequest.Password)
	if _, ok := err.(validator.ValidationErrors); ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if _, ok := err.(*mysql.MySQLError); ok {
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil && err.Error() == "too short password" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
