package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	"github.com/roaris/go-sns-api/httputils"
	"github.com/roaris/go-sns-api/models"
	"github.com/roaris/go-sns-api/swagger/gen"
	"gorm.io/gorm"
)

type PostHandler struct {
	db *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{db}
}

func (p *PostHandler) Create(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
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
	post, err := models.CreatePost(p.db, userID, *createPostRequest.Content)

	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	return http.StatusCreated, post.SwaggerModel(false, 0), nil
}

func (p *PostHandler) Show(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// パスパラメータの取得
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	post, err := models.GetPost(p.db, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, nil, err
	}

	userID := httputils.GetUserIDFromContext(r.Context())

	return http.StatusOK, gen.PostAndUser{
			Post: post.SwaggerModel(models.IsLiked(p.db, userID, post.ID), models.GetLikeNum(p.db, post.ID)),
			User: post.User.SwaggerModel()},
		nil
}

func (p *PostHandler) Index(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	userID := httputils.GetUserIDFromContext(r.Context())
	q := r.URL.Query()
	limit, err := strconv.Atoi(q["limit"][0])
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	offset, err := strconv.Atoi(q["offset"][0])
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	posts := models.GetPosts(p.db, userID, limit, offset)
	var postIDs []int64
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}
	likeFlags := models.BulkIsLiked(p.db, userID, postIDs)
	likeNums := models.BulkGetLikeNum(p.db, postIDs)
	var resPostsAndUsers []*gen.PostAndUser
	for i, post := range posts {
		resPostsAndUsers = append(resPostsAndUsers, &gen.PostAndUser{
			Post: post.SwaggerModel(likeFlags[i], likeNums[i]),
			User: post.User.SwaggerModel(),
		})
	}
	return http.StatusOK, gen.PostsAndUsers{PostsAndUsers: resPostsAndUsers}, nil
}

func (p *PostHandler) Update(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// パスパラメータの取得
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusBadRequest, nil, nil
	}

	// リクエストボディをPostRequestに変換する
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var updatePostRequest gen.UpdatePostRequest
	json.Unmarshal(body, &updatePostRequest)

	if err := updatePostRequest.Validate(strfmt.Default); err != nil {
		return http.StatusBadRequest, nil, err
	}

	userID := httputils.GetUserIDFromContext(r.Context())
	post, err := models.UpdatePost(p.db, id, userID, *updatePostRequest.Content)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, nil, err
	} else if err != nil && err.Error() == "forbidden update" {
		return http.StatusForbidden, nil, err
	}

	return http.StatusOK, post.SwaggerModel(models.IsLiked(p.db, userID, post.ID), models.GetLikeNum(p.db, post.ID)), nil
}

func (p *PostHandler) Destroy(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// パスパラメータの取得
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	userID := httputils.GetUserIDFromContext(r.Context())
	err := models.DeletePost(p.db, id, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, nil, err
	} else if err != nil && err.Error() == "forbidden delete" {
		return http.StatusForbidden, nil, err
	}

	return http.StatusOK, nil, nil
}
