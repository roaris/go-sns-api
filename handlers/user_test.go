package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/roaris/go-sns-api/httputils"

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

		assert.Equal(t, 201, status)
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

func TestShowMe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		// ユーザーを作成しておく
		user, _ := models.CreateUser(db, "taro", "taro@example.com", "password")

		userHandler := NewUserHandler(db)
		r := httptest.NewRequest("GET", "/api/v1/users/me", nil)
		ctx := httputils.SetUserIDToContext(r.Context(), user.ID)
		w := httptest.NewRecorder()
		status, payload, err := userHandler.ShowMe(w, r.WithContext(ctx))

		assert.Equal(t, 200, status)
		assert.Equal(t, user.SwaggerModelWithEmail(), payload)
		assert.Equal(t, nil, err)
	})
}

func TestUpdateMe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		// ユーザーを作成しておく
		user, _ := models.CreateUser(db, "taro", "taro@example.com", "password")
		userHandler := NewUserHandler(db)

		// before
		r := httptest.NewRequest("GET", "/api/v1/users/me", nil)
		ctx := httputils.SetUserIDToContext(r.Context(), user.ID)
		w := httptest.NewRecorder()
		status, payload, err := userHandler.ShowMe(w, r.WithContext(ctx))

		assert.Equal(t, 200, status)
		assert.Equal(t, user.SwaggerModelWithEmail(), payload)
		assert.Equal(t, nil, err)

		// 更新
		rBody := strings.NewReader(`{"name":"taro2", "email":"taro2@example.com", "password":"password2"}`)
		r = httptest.NewRequest("PATCH", "/api/v1/users/me", rBody)
		r.Header.Add("Content-Type", "application/json")
		status, payload, err = userHandler.UpdateMe(w, r.WithContext(ctx))

		assert.Equal(t, 200, status)
		assert.Equal(t, "taro2", payload.(*gen.User).Name)
		assert.Equal(t, "taro2@example.com", payload.(*gen.User).Email)
		assert.Equal(t, nil, err)

		// after
		r = httptest.NewRequest("GET", "/api/v1/users/me", nil)
		status, payload, err = userHandler.ShowMe(w, r.WithContext(ctx))

		assert.Equal(t, 200, status)
		assert.Equal(t, "taro2", payload.(*gen.User).Name)
		assert.Equal(t, "taro2@example.com", payload.(*gen.User).Email)
		assert.Equal(t, nil, err)
	})

	t.Run("bad request", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		// ユーザーを作成しておく
		user, _ := models.CreateUser(db, "taro", "taro@example.com", "password")
		userHandler := NewUserHandler(db)

		// before
		r := httptest.NewRequest("GET", "/api/v1/users/me", nil)
		ctx := httputils.SetUserIDToContext(r.Context(), user.ID)
		w := httptest.NewRecorder()
		status, payload, err := userHandler.ShowMe(w, r.WithContext(ctx))

		assert.Equal(t, 200, status)
		assert.Equal(t, user.SwaggerModelWithEmail(), payload)
		assert.Equal(t, nil, err)

		// 更新(失敗)
		rBody := strings.NewReader(`{"name":"taro2", "email":"taro2@example..com", "password":"password2"}`) // メールアドレスの形式がおかしい
		r = httptest.NewRequest("PATCH", "/api/v1/users/me", rBody)
		r.Header.Add("Content-Type", "application/json")
		status, payload, err = userHandler.UpdateMe(w, r.WithContext(ctx))

		assert.Equal(t, 400, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)

		// after
		r = httptest.NewRequest("GET", "/api/v1/users/me", nil)
		status, payload, err = userHandler.ShowMe(w, r.WithContext(ctx))

		assert.Equal(t, 200, status)
		assert.Equal(t, user.SwaggerModelWithEmail(), payload)
		assert.Equal(t, nil, err)
	})
}
