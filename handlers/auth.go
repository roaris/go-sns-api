package handlers

import (
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"net/http"

	"github.com/roaris/go_sns_api/models"
)

type AuthRequest struct {
	Email    string
	Password string
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	// application/jsonのみ受け付ける
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// リクエストボディをAuthRequestに変換する
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var authRequest AuthRequest
	json.Unmarshal(body, &authRequest)

	user, err := models.GetUserByEmail(authRequest.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// パスワードの検証
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authRequest.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
