package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/roaris/go-sns-api/httputils"
	"github.com/roaris/go-sns-api/swagger/gen"

	"github.com/roaris/go-sns-api/models"
)

func CreateUser(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// application/jsonのみ受け付ける
	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusBadRequest, nil, nil
	}

	// リクエストボディをUserRequestに変換する
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var createUserRequest gen.CreateUserRequest
	json.Unmarshal(body, &createUserRequest)

	if err := createUserRequest.Validate(strfmt.Default); err != nil {
		return http.StatusBadRequest, nil, err
	}

	user, err := models.CreateUser(createUserRequest.Name, string(createUserRequest.Email), createUserRequest.Password)
	if err != nil {
		return http.StatusConflict, nil, err
	}

	return http.StatusOK, user.SwaggerModel(), nil
}

func GetLoginUser(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	userID := httputils.GetUserIDFromContext(r.Context())
	user, _ := models.GetUserById(userID)
	return http.StatusOK, user.SwaggerModelWithEmail(), nil
}

func UpdateLoginUser(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusBadRequest, nil, nil
	}

	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var updateUserRequest gen.UpdateUserRequest
	json.Unmarshal(body, &updateUserRequest)

	if err := updateUserRequest.Validate(strfmt.Default); err != nil {
		return http.StatusBadRequest, nil, err
	}

	userID := httputils.GetUserIDFromContext(r.Context())
	user := models.UpdateUser(
		userID,
		updateUserRequest.Name,
		string(updateUserRequest.Email),
		updateUserRequest.Password,
	)

	return http.StatusOK, user.SwaggerModelWithEmail(), nil
}
