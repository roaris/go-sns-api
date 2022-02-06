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

func GetFollowees(followerID int64) ([]User, error) {
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
	// GormのIN句が動作しないため、仕方なくfor文
	for _, i := range followeeIDs {
		var followee User
		db.First(&followee, "id=?", i)
		followees = append(followees, followee)
	}
	return followees, nil
}

func DeleteFollowee(followerID int64, followeeID int64) error {
	var friendship Friendship
	err := db.First(&friendship, "follower_id=? and followee_id=?", followerID, followeeID).Error
	if err != nil {
		return err
	}
	db.Delete(&friendship)
	return nil
}
