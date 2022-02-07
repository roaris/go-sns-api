package models

import (
	"errors"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/roaris/go-sns-api/swagger/gen"
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

func GetPost(id int64) (post Post, err error) {
	err = db.First(&post, "id=?", id).Error
	if err != nil {
		return post, err
	}
	var user User
	db.Model(&post).Association("User").Find(&user)
	post.User = &user
	return post, err
}

// タイムラインの取得
func GetPosts(userID int64, limit int64, offset int64) (posts []Post, users []User) {
	var friendships []Friendship
	db.Find(&friendships, "follower_id=?", userID)
	var followee_ids []int64
	for _, f := range friendships {
		followee_ids = append(followee_ids, f.FolloweeID)
	}
	db.Limit(limit).Offset(offset).Order("created_at").Find(&posts, "user_id IN (?)", append(followee_ids, userID))
	// N+1
	for _, post := range posts {
		var user User
		db.First(&user, "id=?", post.UserID)
		users = append(users, user)
	}
	return posts, users
}

func CreatePost(userID int64, content string) (post Post, err error) {
	post.UserID = userID
	post.Content = content
	db.Create(&post)
	return post, nil
}

func UpdatePost(id int64, userID int64, content string) (post Post, err error) {
	post, err = GetPost(id)
	if err != nil {
		return post, err
	}
	if post.UserID != userID {
		return post, errors.New("forbidden update")
	}
	postAfter := post
	postAfter.Content = content
	db.Model(&post).Updates(postAfter)
	return post, nil
}

func DeletePost(id int64, userID int64) (err error) {
	post, err := GetPost(id)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return errors.New("forbidden delete")
	}
	db.Delete(&post)
	return nil
}
