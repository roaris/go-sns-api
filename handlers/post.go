package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/roaris/go-sns-api/httputils"
	"github.com/roaris/go-sns-api/models"
	"github.com/roaris/go-sns-api/swagger/gen"
)

func CreatePost(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// application/jsonのみ受け付ける
	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusBadRequest, nil, nil
	}

	// リクエストボディをPostRequestに変換する
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var createPostRequest gen.CreatePostRequest
	json.Unmarshal(body, &createPostRequest)

	if err := createPostRequest.Validate(strfmt.Default); err != nil {
		return http.StatusBadRequest, nil, err
	}

	userID := httputils.GetUserIDFromContext(r.Context())
	post, err := models.CreatePost(userID, createPostRequest.Content)

	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	return http.StatusOK, post.SwaggerModel(), nil
}

func GetPost(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// パスパラメータの取得
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	post, err := models.GetPost(id)
	if gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, nil, err
	}

	return http.StatusOK, gen.PostAndUser{
			Post: post.SwaggerModel(),
			User: post.User.SwaggerModel()},
		nil
}

func GetPosts(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	userID := httputils.GetUserIDFromContext(r.Context())
	q := r.URL.Query()
	limit, _ := strconv.ParseInt(q["limit"][0], 10, 64)
	offset, _ := strconv.ParseInt(q["offset"][0], 10, 64)
	posts := models.GetPosts(userID, limit, offset)
	var resPostsAndUsers []*gen.PostAndUser
	for i := 0; i < len(posts); i++ {
		resPostsAndUsers = append(resPostsAndUsers, &gen.PostAndUser{Post: posts[i].SwaggerModel(), User: posts[i].User.SwaggerModel()})
	}
	return http.StatusOK, gen.PostsAndUsers{PostsAndUsers: resPostsAndUsers}, nil
}

func UpdatePost(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// パスパラメータの取得
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	// リクエストボディをPostRequestに変換する
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var updatePostRequest gen.UpdatePostRequest
	json.Unmarshal(body, &updatePostRequest)

	if err := updatePostRequest.Validate(strfmt.Default); err != nil {
		return http.StatusBadRequest, nil, err
	}

	userID := httputils.GetUserIDFromContext(r.Context())
	post, err := models.UpdatePost(id, userID, updatePostRequest.Content)

	if gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, nil, err
	} else if err != nil && err.Error() == "forbidden update" {
		return http.StatusForbidden, nil, err
	}

	return http.StatusOK, post.SwaggerModel(), nil
}

func DeletePost(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// パスパラメータの取得
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	userID := httputils.GetUserIDFromContext(r.Context())
	err := models.DeletePost(id, userID)
	if gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, nil, err
	} else if err != nil && err.Error() == "forbidden delete" {
		return http.StatusForbidden, nil, err
	}

	return http.StatusOK, nil, nil
}
