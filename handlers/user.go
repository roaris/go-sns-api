package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/roaris/go-sns-api/httputils"
	"github.com/roaris/go-sns-api/swagger/gen"

	"github.com/go-sql-driver/mysql"
	"github.com/roaris/go-sns-api/models"
	"gopkg.in/go-playground/validator.v9"
)

type UserRequest struct {
	Name     string
	Email    string
	Password string
}

func CreateUser(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// application/jsonのみ受け付ける
	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusBadRequest, nil, nil
	}

	// リクエストボディをUserRequestに変換する
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var userRequest UserRequest
	json.Unmarshal(body, &userRequest)

	user, err := models.CreateUser(userRequest.Name, userRequest.Email, userRequest.Password)
	if _, ok := err.(validator.ValidationErrors); ok {
		return http.StatusBadRequest, nil, err
	} else if _, ok := err.(*mysql.MySQLError); ok {
		return http.StatusConflict, nil, err
	} else if err != nil && err.Error() == "too short password" {
		return http.StatusBadRequest, nil, err
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
	var userUpdateRequest gen.UserUpdateRequest
	json.Unmarshal(body, &userUpdateRequest)

	if err := userUpdateRequest.Validate(strfmt.Default); err != nil {
		return http.StatusBadRequest, nil, err
	}

	userID := httputils.GetUserIDFromContext(r.Context())
	user := models.UpdateUser(
		userID,
		userUpdateRequest.Name,
		string(userUpdateRequest.Email),
		userUpdateRequest.Password,
	)

	return http.StatusOK, user.SwaggerModelWithEmail(), nil
}
