package models

import "time"

type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ShowPost(id int) (post Post, err error) {
	err = db.First(&post, "id=?", id).Error
	return post, err
}

func CreatePost(content string) {
	post := Post{}
	post.Content = content
	db.Create(&post)
}
