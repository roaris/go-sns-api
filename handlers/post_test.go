package handlers

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/roaris/go-sns-api/httputils"

	"github.com/roaris/go-sns-api/models"
	"github.com/roaris/go-sns-api/swagger/gen"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	t.Run("created", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user, _ := models.CreateUser(db, "taro", "taro@example.com", "password")

		// before
		var count int
		db.Model(&models.Post{}).Count(&count)
		assert.Equal(t, 0, count)

		// 投稿作成
		postHandler := NewPostHandler(db)
		rBody := strings.NewReader(`{"content":"I'm happy."}`)
		r := httptest.NewRequest("POST", "/api/v1/posts", rBody)
		r.Header.Add("Content-Type", "application/json")
		ctx := httputils.SetUserIDToContext(r.Context(), user.ID)
		w := httptest.NewRecorder()
		status, payload, err := postHandler.Create(w, r.WithContext(ctx))

		assert.Equal(t, 201, status)
		assert.Equal(t, "I'm happy.", payload.(*gen.Post).Content)
		assert.Equal(t, nil, err)

		// after
		db.Model(&models.Post{}).Count(&count)
		assert.Equal(t, 1, count)
	})

	t.Run("bad request", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user, _ := models.CreateUser(db, "taro", "taro@example.com", "password")

		// before
		var count int
		db.Model(&models.Post{}).Count(&count)
		assert.Equal(t, 0, count)

		// 投稿作成(失敗)
		postHandler := NewPostHandler(db)
		rBody := strings.NewReader(fmt.Sprintf(`{"content":"%s"}`, strings.Repeat("a", 141))) // 140文字より多い
		r := httptest.NewRequest("POST", "/api/v1/posts", rBody)
		r.Header.Add("Content-Type", "application/json")
		ctx := httputils.SetUserIDToContext(r.Context(), user.ID)
		w := httptest.NewRecorder()
		status, payload, err := postHandler.Create(w, r.WithContext(ctx))

		assert.Equal(t, 400, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)

		// after
		db.Model(&models.Post{}).Count(&count)
		assert.Equal(t, 0, count)
	})
}
