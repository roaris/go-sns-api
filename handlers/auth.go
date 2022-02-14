package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/roaris/go-sns-api/models"
	"github.com/roaris/go-sns-api/swagger/gen"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db}
}

type Token struct {
	Token string `json:"token"`
}

// JWTトークンを生成する
func generateToken(userID int64, now time.Time) (string, error) {
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

func (a *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// application/jsonのみ受け付ける
	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusBadRequest, nil, nil
	}

	// リクエストボディをAuthRequestに変換する
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var authRequest gen.AuthRequest
	json.Unmarshal(body, &authRequest)

	user, err := models.GetUserByEmail(a.db, authRequest.Email)
	if err != nil {
		return http.StatusUnauthorized, nil, err
	}

	// パスワードの検証
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authRequest.Password))
	if err != nil {
		return http.StatusUnauthorized, nil, err
	}

	// トークンを返す
	token, _ := generateToken(user.ID, time.Now())
	return http.StatusOK, Token{token}, nil
}
