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

func TestIndexPost(t *testing.T) {
	// フォローしているユーザーの投稿が見れるか
	t.Run("ok", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "taro", "taro@example.com", "password")
		user2, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		user3, _ := models.CreateUser(db, "bob", "bob@example.com", "password")
		post1, _ := models.CreatePost(db, user1.ID, "I'm happy.")
		post2, _ := models.CreatePost(db, user2.ID, "I'm sad.")
		models.CreatePost(db, user3.ID, "I'm angry.")
		models.CreateFollowee(db, user1.ID, user2.ID)

		postHandler := NewPostHandler(db)
		r := httptest.NewRequest("GET", "/api/v1/posts?limit=5&offset=0", nil)
		q := r.URL.Query()
		q.Add("limit", "5")
		q.Add("offset", "0")
		ctx := httputils.SetUserIDToContext(r.Context(), user1.ID)
		w := httptest.NewRecorder()
		status, payload, err := postHandler.Index(w, r.WithContext(ctx))

		assert.Equal(t, 200, status)
		assert.Equal(t, 2, len(payload.(gen.PostsAndUsers).PostsAndUsers))
		assert.Equal(t, post1.SwaggerModel().Content, payload.(gen.PostsAndUsers).PostsAndUsers[0].Post.Content)
		assert.Equal(t, post2.SwaggerModel().Content, payload.(gen.PostsAndUsers).PostsAndUsers[1].Post.Content)
		assert.Equal(t, nil, err)
	})

	// limitとoffsetが正しく動くか
	t.Run("ok", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "taro", "taro@example.com", "password")
		user2, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		var posts []models.Post
		for i := 0; i < 10; i++ {
			post, _ := models.CreatePost(db, user1.ID, fmt.Sprintf("I'm happy%d.", i))
			posts = append(posts, post)
			post, _ = models.CreatePost(db, user2.ID, fmt.Sprintf("I'm sad%d.", i))
			posts = append(posts, post)
		}
		models.CreateFollowee(db, user1.ID, user2.ID)

		postHandler := NewPostHandler(db)
		r := httptest.NewRequest("GET", "/api/v1/posts?limit=10&offset=5", nil)
		q := r.URL.Query()
		q.Add("limit", "10")
		q.Add("offset", "5")
		ctx := httputils.SetUserIDToContext(r.Context(), user1.ID)
		w := httptest.NewRecorder()
		status, payload, err := postHandler.Index(w, r.WithContext(ctx))

		assert.Equal(t, 200, status)
		assert.Equal(t, 10, len(payload.(gen.PostsAndUsers).PostsAndUsers))
		assert.Equal(t, posts[5].Content, payload.(gen.PostsAndUsers).PostsAndUsers[0].Post.Content)
		assert.Equal(t, posts[14].Content, payload.(gen.PostsAndUsers).PostsAndUsers[9].Post.Content)
		assert.Equal(t, nil, err)
	})
}

func TestUpdatePost(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user, _ := models.CreateUser(db, "taro", "taro@example.com", "password")
		post, _ := models.CreatePost(db, user.ID, "I'm happy.")

		postHandler := NewPostHandler(db)

		// before
		postBefore, _ := models.GetPost(db, post.ID)
		assert.Equal(t, post.Content, postBefore.Content)

		// 更新
		rBody := strings.NewReader(`{"content":"I'm sad."}`)
		r := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v1/posts/%d", post.ID), rBody)
		vars := map[string]string{
			"id": strconv.FormatInt(post.ID, 10),
		}
		r = mux.SetURLVars(r, vars)
		ctx := httputils.SetUserIDToContext(r.Context(), user.ID)
		w := httptest.NewRecorder()
		status, payload, err := postHandler.Update(w, r.WithContext(ctx))

		assert.Equal(t, 200, status)
		assert.Equal(t, "I'm sad.", payload.(*gen.Post).Content)
		assert.Equal(t, nil, err)

		// after
		postAfter, _ := models.GetPost(db, post.ID)
		assert.Equal(t, "I'm sad.", postAfter.Content)
	})

	t.Run("bad request", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user, _ := models.CreateUser(db, "taro", "taro@example.com", "password")
		post, _ := models.CreatePost(db, user.ID, "I'm happy.")

		postHandler := NewPostHandler(db)

		// before
		postBefore, _ := models.GetPost(db, post.ID)
		assert.Equal(t, post.Content, postBefore.Content)

		// 更新
		rBody := strings.NewReader(fmt.Sprintf(`{"content":"%s"}`, strings.Repeat("a", 141))) // 140文字より多い
		r := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v1/posts/%d", post.ID), rBody)
		vars := map[string]string{
			"id": strconv.FormatInt(post.ID, 10),
		}
		r = mux.SetURLVars(r, vars)
		ctx := httputils.SetUserIDToContext(r.Context(), user.ID)
		w := httptest.NewRecorder()
		status, payload, err := postHandler.Update(w, r.WithContext(ctx))

		assert.Equal(t, 400, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)

		// after
		postAfter, _ := models.GetPost(db, post.ID)
		assert.Equal(t, post.Content, postAfter.Content)
	})

	t.Run("forbidden", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "taro", "taro@example.com", "password")
		user2, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		post, _ := models.CreatePost(db, user1.ID, "I'm happy.")

		postHandler := NewPostHandler(db)

		// before
		postBefore, _ := models.GetPost(db, post.ID)
		assert.Equal(t, post.Content, postBefore.Content)

		// 更新
		rBody := strings.NewReader(`{"content":"I'm sad."}`)
		r := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v1/posts/%d", post.ID), rBody)
		vars := map[string]string{
			"id": strconv.FormatInt(post.ID, 10),
		}
		r = mux.SetURLVars(r, vars)
		ctx := httputils.SetUserIDToContext(r.Context(), user2.ID)
		w := httptest.NewRecorder()
		status, payload, err := postHandler.Update(w, r.WithContext(ctx))

		assert.Equal(t, 403, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)

		// after
		postAfter, _ := models.GetPost(db, post.ID)
		assert.Equal(t, post.Content, postAfter.Content)
	})
}

func TestDestroyPost(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user, _ := models.CreateUser(db, "taro", "taro@example.com", "password")
		post, _ := models.CreatePost(db, user.ID, "I'm happy.")

		postHandler := NewPostHandler(db)

		// before
		var count int
		db.Model(&models.Post{}).Count(&count)
		assert.Equal(t, 1, count)

		// 削除
		r := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/posts/%d", post.ID), nil)
		vars := map[string]string{
			"id": strconv.FormatInt(post.ID, 10),
		}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()
		ctx := httputils.SetUserIDToContext(r.Context(), user.ID)
		status, payload, err := postHandler.Destroy(w, r.WithContext(ctx))

		assert.Equal(t, 200, status)
		assert.Equal(t, nil, payload)
		assert.Equal(t, nil, err)

		// after
		db.Model(&models.Post{}).Count(&count)
		assert.Equal(t, 0, count)
	})

	t.Run("forbidden", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "taro", "taro@example.com", "password")
		user2, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		post, _ := models.CreatePost(db, user1.ID, "I'm happy.")

		postHandler := NewPostHandler(db)

		// before
		var count int
		db.Model(&models.Post{}).Count(&count)
		assert.Equal(t, 1, count)

		// 削除
		r := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/posts/%d", post.ID), nil)
		vars := map[string]string{
			"id": strconv.FormatInt(post.ID, 10),
		}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()
		ctx := httputils.SetUserIDToContext(r.Context(), user2.ID)
		status, payload, err := postHandler.Destroy(w, r.WithContext(ctx))

		assert.Equal(t, 403, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)

		// after
		db.Model(&models.Post{}).Count(&count)
		assert.Equal(t, 1, count)
	})
}
