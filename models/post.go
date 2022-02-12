package models

import (
	"errors"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/jinzhu/gorm"
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

func GetPost(db *gorm.DB, id int64) (post Post, err error) {
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
func GetPosts(db *gorm.DB, userID int64, limit int64, offset int64) (posts []Post) {
	var friendships []Friendship
	db.Find(&friendships, "follower_id=?", userID)
	var followee_ids []int64
	for _, f := range friendships {
		followee_ids = append(followee_ids, f.FolloweeID)
	}
	// idにはインデックスがついてるので、created_atでソートするよりも高速
	db.Preload("User").Limit(limit).Offset(offset).Order("id").Find(&posts, "user_id IN (?)", append(followee_ids, userID))
	return posts
}

func CreatePost(db *gorm.DB, userID int64, content string) (post Post, err error) {
	post.UserID = userID
	post.Content = content
	db.Create(&post)
	return post, nil
}

func UpdatePost(db *gorm.DB, id int64, userID int64, content string) (post Post, err error) {
	post, err = GetPost(db, id)
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

func DeletePost(db *gorm.DB, id int64, userID int64) (err error) {
	post, err := GetPost(db, id)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return errors.New("forbidden delete")
	}
	db.Delete(&post)
	return nil
}
