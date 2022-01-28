package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/roaris/go_sns_api/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string
	Password string
}

type Token struct {
	Token string `json:"token"`
}

// JWTトークンを生成する
func GenerateToken(userID int, now time.Time) (string, error) {
	/*
		HMAC SHA-256を使用 sub(subject):識別子 iat(issued at):発行時刻 exp(expiration):有効期限
		iatとexpはUNIXタイムスタンプを使う
	*/
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(time.Hour * 24).Unix(),
	})
	// 秘密鍵で署名を作成
	return token.SignedString([]byte(os.Getenv("SECRET")))
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

	// トークンを返す
	token, _ := GenerateToken(user.ID, time.Now())
	res, _ := json.Marshal(Token{token})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
