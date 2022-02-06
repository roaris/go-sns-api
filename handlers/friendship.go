package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/roaris/go-sns-api/httputils"
	"github.com/roaris/go-sns-api/models"

	"github.com/roaris/go-sns-api/swagger/gen"
)

func CreateFollowee(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
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
	if err := models.CreateFollowee(userID, createFolloweeRequest.FolloweeID); err != nil {
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
