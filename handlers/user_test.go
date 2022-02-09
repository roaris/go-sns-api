package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/roaris/go-sns-api/models"
	"github.com/roaris/go-sns-api/swagger/gen"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		userHandler := NewUserHandler(db)
		rBody := strings.NewReader(`{"name":"test", "email":"test@example.com", "password":"password"}`)
		r := httptest.NewRequest("POST", "/api/v1/users", rBody)
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		status, payload, err := userHandler.Create(w, r)

		assert.Equal(t, 200, status)
		assert.Equal(t, "test", payload.(*gen.User).Name)
		assert.Equal(t, nil, err)
	})

	t.Run("bad request", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		userHandler := NewUserHandler(db)
		rBody := strings.NewReader(`{"name":"test", "email":"test@example.com", "password":"pass"}`) // passwordが短い
		r := httptest.NewRequest("POST", "/api/v1/users", rBody)
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		status, payload, err := userHandler.Create(w, r)

		assert.Equal(t, 400, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)
	})

	t.Run("conflict", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		userHandler := NewUserHandler(db)
		rBody := strings.NewReader(`{"name":"test", "email":"test@example.com", "password":"password"}`)
		r := httptest.NewRequest("POST", "/api/v1/users", rBody)
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		userHandler.Create(w, r)

		rBody = strings.NewReader(`{"name":"test", "email":"test@example.com", "password":"password"}`)
		r = httptest.NewRequest("POST", "/api/v1/users", rBody)
		r.Header.Add("Content-Type", "application/json")
		status, payload, err := userHandler.Create(w, r) // メールアドレスが重複する

		assert.Equal(t, 409, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)
	})
}
