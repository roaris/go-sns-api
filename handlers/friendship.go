package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/roaris/go-sns-api/httputils"
	"github.com/roaris/go-sns-api/models"

	"github.com/roaris/go-sns-api/swagger/gen"
)

type FriendshipHandler struct {
	db *gorm.DB
}

func NewFriendshipHandler(db *gorm.DB) *FriendshipHandler {
	return &FriendshipHandler{db}
}

func (f *FriendshipHandler) CreateFollowee(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusBadRequest, nil, nil
	}

	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var createFolloweeRequest gen.CreateFolloweeRequest
	if err := json.Unmarshal(body, &createFolloweeRequest); err != nil {
		return http.StatusBadRequest, nil, err
	}

	userID := httputils.GetUserIDFromContext(r.Context())
	if err := models.CreateFollowee(f.db, userID, createFolloweeRequest.FolloweeID); err != nil {
		errCode := "Error XXXX"
		l := len(errCode)
		errMessage := err.Error()
		if len(errMessage) >= l && errMessage[:l] == "Error 1062" { // ユニーク制約に違反する
			return http.StatusConflict, nil, err
		} else if len(errMessage) >= l && errMessage[:l] == "Error 1452" { // 外部キー制約に違反する
			return http.StatusNotFound, nil, err
		} else if errMessage == "forbidden follow" {
			return http.StatusBadRequest, nil, err
		}
	}

	return http.StatusNoContent, nil, nil
}

func (f *FriendshipHandler) GetFollowees(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	followees, err := models.GetFollowees(f.db, userID)
	if err != nil {
		return http.StatusNotFound, nil, err
	}
	var resFollowees []*gen.User
	for _, f := range followees {
		resFollowees = append(resFollowees, f.SwaggerModel())
	}
	return http.StatusOK, gen.Followees{Followees: resFollowees}, nil
}

func (f *FriendshipHandler) GetFollowers(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	followers, err := models.GetFollowers(f.db, userID)
	if err != nil {
		return http.StatusNotFound, nil, err
	}
	var resFollowers []*gen.User
	for _, f := range followers {
		resFollowers = append(resFollowers, f.SwaggerModel())
	}
	return http.StatusOK, gen.Followers{Followers: resFollowers}, nil
}

func (f *FriendshipHandler) DeleteFollowee(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	followerID := httputils.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	followeeID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if err = models.DeleteFollowee(f.db, followerID, followeeID); err != nil {
		return http.StatusNotFound, nil, err
	}

	return http.StatusNoContent, nil, nil
}
