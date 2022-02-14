package models

import "gorm.io/gorm"

type Like struct {
	ID     int64
	UserID int64 `gorm:"index:idx_like,unique"`
	User   User  `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
	PostID int64 `gorm:"index:idx_like,unique"`
	Post   Post  `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
}

func CreateLike(db *gorm.DB, userID int64, postID int64) error {
	err := db.Create(&Like{
		UserID: userID,
		PostID: postID,
	}).Error
	return err
}
