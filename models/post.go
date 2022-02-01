package models

import (
	"errors"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/roaris/go-sns-api/swagger/gen"
	"gopkg.in/go-playground/validator.v9"
)

type Post struct {
	ID        int64
	Content   string `validate:"required,max=140"`
	UserID    int64  `validate:"required"`
	User      *User  `validate:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Post) SwaggerModel() *gen.Post {
	return &gen.Post{
		ID:        p.ID,
		Content:   p.Content,
		UserID:    p.UserID,
		CreatedAt: strfmt.DateTime(p.CreatedAt),
		UpdatedAt: strfmt.DateTime(p.UpdatedAt),
	}
}

func ShowPost(id int64) (post Post, err error) {
	err = db.First(&post, "id=?", id).Error
	return post, err
}

func CreatePost(userID int64, content string) (post Post, err error) {
	post.UserID = userID
	post.Content = content
	validate := validator.New()
	err = validate.Struct(post)
	if err != nil {
		return post, err
	}
	db.Create(&post)
	return post, nil
}

func UpdatePost(id int64, userID int64, content string) (post Post, err error) {
	post, err = ShowPost(id)
	if err != nil {
		return post, err
	}
	if post.UserID != userID {
		return post, errors.New("forbidden update")
	}
	postAfter := post
	postAfter.Content = content
	validate := validator.New()
	err = validate.Struct(postAfter)
	if err != nil {
		return post, err
	}
	db.Model(&post).Updates(postAfter)
	return post, nil
}

func DeletePost(id int64, userID int64) (err error) {
	post, err := ShowPost(id)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return errors.New("forbidden delete")
	}
	db.Delete(&post)
	return nil
}
