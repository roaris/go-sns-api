package models

import (
	"errors"

	"gorm.io/gorm"
)

type Friendship struct {
	ID         int64
	FollowerID int64 `gorm:"index:idx_friendship,unique"`
	Follower   User  `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
	FolloweeID int64 `gorm:"index:idx_friendship,unique"`
	Followee   User  `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
}

func CreateFollowee(db *gorm.DB, followerID int64, followeeID int64) error {
	if followerID == followeeID {
		return errors.New("forbidden follow")
	}
	err := db.Create(&Friendship{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}).Error
	return err
}

func GetFollowees(db *gorm.DB, followerID int64) ([]User, error) {
	var user User
	err := db.First(&user, "id=?", followerID).Error
	if err != nil {
		return nil, err
	}

	var friendships []Friendship
	db.Find(&friendships, "follower_id=?", followerID)
	var followeeIDs []int64
	for _, f := range friendships {
		followeeIDs = append(followeeIDs, f.FolloweeID)
	}
	var followees []User
	db.Find(&followees, "id IN (?)", followeeIDs)
	return followees, nil
}

func GetFollowers(db *gorm.DB, followeeID int64) ([]User, error) {
	var user User
	err := db.First(&user, "id=?", followeeID).Error
	if err != nil {
		return nil, err
	}

	var friendships []Friendship
	db.Find(&friendships, "followee_id=?", followeeID)
	var followerIDs []int64
	for _, f := range friendships {
		followerIDs = append(followerIDs, f.FollowerID)
	}
	var followers []User
	db.Find(&followers, "id IN (?)", followerIDs)
	return followers, nil
}

func DeleteFollowee(db *gorm.DB, followerID int64, followeeID int64) error {
	var friendship Friendship
	err := db.First(&friendship, "follower_id=? and followee_id=?", followerID, followeeID).Error
	if err != nil {
		return err
	}
	db.Delete(&friendship)
	return nil
}
