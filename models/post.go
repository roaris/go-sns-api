package models

import (
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content" validate:"required,max=140"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ShowPost(id int) (post Post, err error) {
	err = db.First(&post, "id=?", id).Error
	return post, err
}

func CreatePost(content string) (err error) {
	post := Post{}
	post.Content = content
	validate := validator.New()
	err = validate.Struct(post)
	if err != nil {
		return err
	}
	db.Create(&post)
	return err
}

func UpdatePost(id int, content string) (err error) {
	post, err := ShowPost(id)
	if err != nil {
		return err
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

func DeletePost(id int) (err error) {
	post, err := ShowPost(id)
	if err != nil {
		return err
	}
	db.Delete(&post)
	return nil
}
