package handlers

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/roaris/go-sns-api/httputils"
	"github.com/roaris/go-sns-api/models"
	"github.com/roaris/go-sns-api/swagger/gen"
	"github.com/stretchr/testify/assert"
)

func TestCreateLike(t *testing.T) {
	t.Run("created", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		user2, _ := models.CreateUser(db, "bob", "bob@example.com", "password")
		post, _ := models.CreatePost(db, user2.ID, "I'm happy.")

		// before
		postHandler := NewPostHandler(db)
		r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/posts/%d", post.ID), nil)
		vars := map[string]string{
			"id": strconv.FormatInt(post.ID, 10),
		}
		r = mux.SetURLVars(r, vars)
		ctx := httputils.SetUserIDToContext(r.Context(), user1.ID)
		w := httptest.NewRecorder()
		_, payload, _ := postHandler.Show(w, r.WithContext(ctx))

		assert.Equal(t, int64(0), *payload.(gen.PostAndUser).Post.LikeNum)
		assert.Equal(t, false, *payload.(gen.PostAndUser).Post.IsLiked)

		// いいねする
		likeHandler := NewLikeHandler(db)
		r = httptest.NewRequest("POST", fmt.Sprintf("/api/v1/posts/%d/likes", post.ID), nil)
		status, payload, err := likeHandler.Create(w, r.WithContext(ctx))

		assert.Equal(t, 201, status)
		assert.Equal(t, nil, payload)
		assert.Equal(t, nil, err)

		// after
		r = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/posts/%d", post.ID), nil)
		_, payload, _ = postHandler.Show(w, r.WithContext(ctx))

		assert.Equal(t, int64(1), *payload.(gen.PostAndUser).Post.LikeNum)
		assert.Equal(t, true, *payload.(gen.PostAndUser).Post.IsLiked)
	})

	t.Run("conflict", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		user2, _ := models.CreateUser(db, "bob", "bob@example.com", "password")
		post, _ := models.CreatePost(db, user2.ID, "I'm happy.")
		models.CreateLike(db, user1.ID, post.ID)

		// before
		postHandler := NewPostHandler(db)
		r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/posts/%d", post.ID), nil)
		vars := map[string]string{
			"id": strconv.FormatInt(post.ID, 10),
		}
		r = mux.SetURLVars(r, vars)
		ctx := httputils.SetUserIDToContext(r.Context(), user1.ID)
		w := httptest.NewRecorder()
		_, payload, _ := postHandler.Show(w, r.WithContext(ctx))

		assert.Equal(t, int64(1), *payload.(gen.PostAndUser).Post.LikeNum)
		assert.Equal(t, true, *payload.(gen.PostAndUser).Post.IsLiked)

		// いいねする(既にしているので失敗)
		likeHandler := NewLikeHandler(db)
		r = httptest.NewRequest("POST", fmt.Sprintf("/api/v1/posts/%d/likes", post.ID), nil)
		status, payload, err := likeHandler.Create(w, r.WithContext(ctx))

		assert.Equal(t, 409, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)

		// after
		r = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/posts/%d", post.ID), nil)
		_, payload, _ = postHandler.Show(w, r.WithContext(ctx))

		assert.Equal(t, int64(1), *payload.(gen.PostAndUser).Post.LikeNum)
		assert.Equal(t, true, *payload.(gen.PostAndUser).Post.IsLiked)
	})
}
