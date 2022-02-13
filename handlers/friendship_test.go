package handlers

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/roaris/go-sns-api/httputils"
	"github.com/roaris/go-sns-api/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateFollowee(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		user2, _ := models.CreateUser(db, "bob", "bob@example.com", "password")

		// before
		beforeFollowees, _ := models.GetFollowees(db, user1.ID)
		assert.Equal(t, 0, len(beforeFollowees))

		// フォローする
		friendshipHandler := NewFriendshipHandler(db)
		rBody := strings.NewReader(fmt.Sprintf(`{"followee_id":%d}`, user2.ID))
		r := httptest.NewRequest("POST", "/api/v1/users/me/followees", rBody)
		r.Header.Add("Content-Type", "application/json")
		ctx := httputils.SetUserIDToContext(r.Context(), user1.ID)
		w := httptest.NewRecorder()
		status, payload, err := friendshipHandler.Create(w, r.WithContext(ctx))

		assert.Equal(t, 201, status)
		assert.Equal(t, nil, payload)
		assert.Equal(t, nil, err)

		// after
		afterFollowees, _ := models.GetFollowees(db, user1.ID)
		assert.Equal(t, 1, len(afterFollowees))
	})

	t.Run("conflict", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		user2, _ := models.CreateUser(db, "bob", "bob@example.com", "password")
		models.CreateFollowee(db, user1.ID, user2.ID)

		// before
		beforeFollowees, _ := models.GetFollowees(db, user1.ID)
		assert.Equal(t, 1, len(beforeFollowees))

		// フォローする(フォロー済みなので失敗)
		friendshipHandler := NewFriendshipHandler(db)
		rBody := strings.NewReader(fmt.Sprintf(`{"followee_id":%d}`, user2.ID))
		r := httptest.NewRequest("POST", "/api/v1/users/me/followees", rBody)
		r.Header.Add("Content-Type", "application/json")
		ctx := httputils.SetUserIDToContext(r.Context(), user1.ID)
		w := httptest.NewRecorder()
		status, payload, err := friendshipHandler.Create(w, r.WithContext(ctx))

		assert.Equal(t, 409, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)

		// after
		afterFollowees, _ := models.GetFollowees(db, user1.ID)
		assert.Equal(t, 1, len(afterFollowees))
	})

	t.Run("not found", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		user2, _ := models.CreateUser(db, "bob", "bob@example.com", "password")

		// before
		beforeFollowees, _ := models.GetFollowees(db, user1.ID)
		assert.Equal(t, 0, len(beforeFollowees))

		// フォローする(存在しないユーザーを指定しているので失敗)
		friendshipHandler := NewFriendshipHandler(db)
		rBody := strings.NewReader(fmt.Sprintf(`{"followee_id":%d}`, user2.ID+1))
		r := httptest.NewRequest("POST", "/api/v1/users/me/followees", rBody)
		r.Header.Add("Content-Type", "application/json")
		ctx := httputils.SetUserIDToContext(r.Context(), user1.ID)
		w := httptest.NewRecorder()
		status, payload, err := friendshipHandler.Create(w, r.WithContext(ctx))

		assert.Equal(t, 404, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)

		// after
		afterFollowees, _ := models.GetFollowees(db, user1.ID)
		assert.Equal(t, 0, len(afterFollowees))
	})

	t.Run("bad request", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user, _ := models.CreateUser(db, "alice", "alice@example.com", "password")

		// before
		beforeFollowees, _ := models.GetFollowees(db, user.ID)
		assert.Equal(t, 0, len(beforeFollowees))

		// フォローする(自分自身を指定しているので失敗)
		friendshipHandler := NewFriendshipHandler(db)
		rBody := strings.NewReader(fmt.Sprintf(`{"followee_id":%d}`, user.ID))
		r := httptest.NewRequest("POST", "/api/v1/users/me/followees", rBody)
		r.Header.Add("Content-Type", "application/json")
		ctx := httputils.SetUserIDToContext(r.Context(), user.ID)
		w := httptest.NewRecorder()
		status, payload, err := friendshipHandler.Create(w, r.WithContext(ctx))

		assert.Equal(t, 400, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)

		// after
		afterFollowees, _ := models.GetFollowees(db, user.ID)
		assert.Equal(t, 0, len(afterFollowees))
	})
}
