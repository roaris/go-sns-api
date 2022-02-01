package middlewares

import (
	"errors"
	"net/http"
	"os"

	"github.com/roaris/go-sns-api/httputils"
	"github.com/roaris/go-sns-api/models"

	"github.com/dgrijalva/jwt-go"
)

// ヘッダからトークンを取得する
func getTokenFromHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("authorization header not found")
	}
	bearer := "Bearer"
	l := len(bearer)
	if len(header) > l && header[:l] == bearer {
		return header[l+1:], nil
	}
	return "", errors.New("invalid format token")
}

// JWTトークンの検証を行う
func parseToken(signedString string) (int64, error) {
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return -1, err
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(float64)
	return int64(userID), err
}

// JWTトークンの検証を行うミドルウェア
func AuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromHeader(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := parseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user, _ := models.ShowUser(userID)
		ctx := httputils.SetUserToContext(r.Context(), user)
		handler(w, r.WithContext(ctx))
	}
}
