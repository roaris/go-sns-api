package models

import (
	"errors"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content" validate:"required,max=140"`
	UserID    int       `json:"user_id" validate:"required"`
	User      User      `validate:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ShowPost(id int) (post Post, err error) {
	err = db.First(&post, "id=?", id).Error
	return post, err
}

func CreatePost(userID int, content string) (err error) {
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

func UpdatePost(id int, userID int, content string) (err error) {
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

func DeletePost(id int, userID int) (err error) {
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
