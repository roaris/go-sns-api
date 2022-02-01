package models

import (
	"errors"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content" validate:"required,max=140"`
	UserID    int64     `json:"user_id" validate:"required"`
	User      *User     `json:"user,omitempty" validate:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ShowPost(id int64) (post Post, err error) {
	err = db.First(&post, "id=?", id).Error
	return post, err
}

func CreatePost(userID int64, content string) (err error) {
	post := Post{}
	post.UserID = userID
	post.Content = content
	validate := validator.New()
	err = validate.Struct(post)
	if err != nil {
		return err
	}
	db.Create(&post)
	return nil
}

func UpdatePost(id int64, userID int64, content string) (err error) {
	post, err := ShowPost(id)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return errors.New("forbidden update")
	}
	postAfter := post
	postAfter.Content = content
	validate := validator.New()
	err = validate.Struct(postAfter)
	if err != nil {
		return err
	}
	db.Model(&post).Updates(postAfter)
	return nil
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
