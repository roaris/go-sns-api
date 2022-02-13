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

func TestShowFollowees(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		user2, _ := models.CreateUser(db, "bob", "bob@example.com", "password")
		models.CreateFollowee(db, user1.ID, user2.ID)

		friendshipHandler := NewFriendshipHandler(db)
		r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d/followees", user1.ID), nil)
		vars := map[string]string{
			"id": strconv.FormatInt(user1.ID, 10),
		}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()
		status, payload, err := friendshipHandler.ShowFollowees(w, r)

		assert.Equal(t, 200, status)
		assert.Equal(t, 1, len(payload.(gen.Followees).Followees))
		assert.Equal(t, user2.Name, payload.(gen.Followees).Followees[0].Name)
		assert.Equal(t, nil, err)
	})

	t.Run("not found", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		user2, _ := models.CreateUser(db, "bob", "bob@example.com", "password")
		models.CreateFollowee(db, user1.ID, user2.ID)

		friendshipHandler := NewFriendshipHandler(db)
		r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d/followees", 10), nil)
		vars := map[string]string{
			"id": strconv.FormatInt(10, 10),
		}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()
		status, payload, err := friendshipHandler.ShowFollowees(w, r)

		assert.Equal(t, 404, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)
	})
}

func TestShowFollowers(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		user2, _ := models.CreateUser(db, "bob", "bob@example.com", "password")
		models.CreateFollowee(db, user1.ID, user2.ID)

		friendshipHandler := NewFriendshipHandler(db)
		r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d/followers", user2.ID), nil)
		vars := map[string]string{
			"id": strconv.FormatInt(user2.ID, 10),
		}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()
		status, payload, err := friendshipHandler.ShowFollowers(w, r)

		assert.Equal(t, 200, status)
		assert.Equal(t, 1, len(payload.(gen.Followers).Followers))
		assert.Equal(t, user1.Name, payload.(gen.Followers).Followers[0].Name)
		assert.Equal(t, nil, err)
	})

	t.Run("not found", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		user2, _ := models.CreateUser(db, "bob", "bob@example.com", "password")
		models.CreateFollowee(db, user1.ID, user2.ID)

		friendshipHandler := NewFriendshipHandler(db)
		r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d/followers", 10), nil)
		vars := map[string]string{
			"id": strconv.FormatInt(10, 10),
		}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()
		status, payload, err := friendshipHandler.ShowFollowers(w, r)

		assert.Equal(t, 404, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)
	})
}

func TestDestroyFollowees(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		user2, _ := models.CreateUser(db, "bob", "bob@example.com", "password")
		models.CreateFollowee(db, user1.ID, user2.ID)

		// before
		beforeFollowees, _ := models.GetFollowees(db, user1.ID)
		assert.Equal(t, 1, len(beforeFollowees))

		// フォロー解除する
		friendshipHandler := NewFriendshipHandler(db)
		r := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/me/followees/%d", user2.ID), nil)
		vars := map[string]string{
			"id": strconv.FormatInt(user2.ID, 10),
		}
		r = mux.SetURLVars(r, vars)
		ctx := httputils.SetUserIDToContext(r.Context(), user1.ID)
		w := httptest.NewRecorder()
		status, payload, err := friendshipHandler.Destroy(w, r.WithContext(ctx))

		assert.Equal(t, 204, status)
		assert.Equal(t, nil, payload)
		assert.Equal(t, nil, err)

		// after
		afterFollowees, _ := models.GetFollowees(db, user1.ID)
		assert.Equal(t, 0, len(afterFollowees))
	})

	t.Run("not found", func(t *testing.T) {
		db := models.CreateTestDB()
		defer models.CleanUpTestDB(db)

		user1, _ := models.CreateUser(db, "alice", "alice@example.com", "password")
		user2, _ := models.CreateUser(db, "bob", "bob@example.com", "password")
		models.CreateFollowee(db, user1.ID, user2.ID)

		// before
		beforeFollowees, _ := models.GetFollowees(db, user1.ID)
		assert.Equal(t, 1, len(beforeFollowees))

		// フォロー解除する(存在しないユーザーを指定しているので失敗)
		friendshipHandler := NewFriendshipHandler(db)
		r := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/me/followees/%d", user2.ID+1), nil)
		vars := map[string]string{
			"id": strconv.FormatInt(user2.ID+1, 10),
		}
		r = mux.SetURLVars(r, vars)
		ctx := httputils.SetUserIDToContext(r.Context(), user1.ID)
		w := httptest.NewRecorder()
		status, payload, err := friendshipHandler.Destroy(w, r.WithContext(ctx))

		assert.Equal(t, 404, status)
		assert.Equal(t, nil, payload)
		assert.NotEqual(t, nil, err)

		// after
		afterFollowees, _ := models.GetFollowees(db, user1.ID)
		assert.Equal(t, 1, len(afterFollowees))
	})
}
