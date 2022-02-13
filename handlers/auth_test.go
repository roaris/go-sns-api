package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/roaris/go-sns-api/models"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		models.CreateUser(db, "taro", "taro@example.com", "password")

		authHandler := NewAuthHandler(db)
		rBody := strings.NewReader(`{"email":"taro@example.com", "password":"password"}`)
		r := httptest.NewRequest("POST", "/api/v1/auth", rBody)
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		status, payload, err := authHandler.Authenticate(w, r)

		assert.Equal(t, 200, status)
		assert.NotEqual(t, nil, payload.(Token).Token)
		assert.Equal(t, nil, err)
	})

	t.Run("unauthorized", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		models.CreateUser(db, "taro", "taro@example.com", "password")

		authHandler := NewAuthHandler(db)
		rBody := strings.NewReader(`{"email":"taro@example.com", "password":"pass"}`) // パスワードが違う
		r := httptest.NewRequest("POST", "/api/v1/auth", rBody)
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		status, payload, err := authHandler.Authenticate(w, r)

		assert.Equal(t, 401, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)
	})
}
