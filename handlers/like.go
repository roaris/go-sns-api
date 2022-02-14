package handlers

import (
	"net/http"
	"strconv"

	"github.com/roaris/go-sns-api/httputils"
	"github.com/roaris/go-sns-api/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type LikeHandler struct {
	db *gorm.DB
}

func NewLikeHandler(db *gorm.DB) *LikeHandler {
	return &LikeHandler{db}
}

func (l *LikeHandler) Create(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	userID := httputils.GetUserIDFromContext(r.Context())
	if err := models.CreateLike(l.db, userID, postID); err != nil {
		errCode := "Error XXXX"
		l := len(errCode)
		errMessage := err.Error()
		if len(errMessage) >= l && errMessage[:l] == "Error 1062" { // ユニーク制約に違反する
			return http.StatusConflict, nil, err
		} else if len(errMessage) >= l && errMessage[:l] == "Error 1452" { // 外部キー制約に違反する
			return http.StatusNotFound, nil, err
		}
	}
	return http.StatusCreated, nil, nil
}

func (l *LikeHandler) Destroy(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	userID := httputils.GetUserIDFromContext(r.Context())
	if err = models.DeleteLike(l.db, userID, postID); err != nil {
		return http.StatusNotFound, nil, err
	}
	return http.StatusOK, nil, nil
}
