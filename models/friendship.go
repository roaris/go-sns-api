package models

import "errors"

type Friendship struct {
	ID         int64
	FollowerID int64
	FolloweeID int64
}

func CreateFollowee(followerID int64, followeeID int64) error {
	if followerID == followeeID {
		return errors.New("forbidden follow")
	}
	err := db.Create(&Friendship{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}).Error
	return err
}
