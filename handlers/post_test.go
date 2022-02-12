package handlers

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"

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

func TestShowPost(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user, _ := models.CreateUser(db, "taro", "taro@example.com", "password")
		post, _ := models.CreatePost(db, user.ID, "I'm happy.")

		postHandler := NewPostHandler(db)
		r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/posts/%d", post.ID), nil)
		vars := map[string]string{
			"id": strconv.FormatInt(post.ID, 10),
		}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()
		status, payload, err := postHandler.Show(w, r)

		assert.Equal(t, 200, status)
		assert.Equal(t, post.Content, payload.(gen.PostAndUser).Post.Content)
		assert.Equal(t, user.Name, payload.(gen.PostAndUser).User.Name)
		assert.Equal(t, nil, err)
	})

	t.Run("not found", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user, _ := models.CreateUser(db, "taro", "taro@example.com", "password")
		post, _ := models.CreatePost(db, user.ID, "I'm happy.")

		postHandler := NewPostHandler(db)
		r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/posts/%d", post.ID+1), nil)
		vars := map[string]string{
			"id": strconv.FormatInt(post.ID+1, 10),
		}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()
		status, payload, err := postHandler.Show(w, r)

		assert.Equal(t, 404, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)
	})
}
