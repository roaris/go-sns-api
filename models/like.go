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

func DeleteLike(db *gorm.DB, userID int64, postID int64) error {
	var like Like
	if err := db.First(&like, "user_id = ? & post_id = ?", userID, postID).Error; err != nil {
		return err
	}
	db.Delete(&like)
	return nil
}

func IsLiked(db *gorm.DB, userID int64, postID int64) bool {
	var like Like
	if err := db.First(&like, "user_id = ? & post_id = ?", userID, postID).Error; err != nil {
		return false
	} else {
		return true
	}
}

func BulkIsLiked(db *gorm.DB, userID int64, postIDs []int64) (flags []bool) {
	var likes []Like
	db.Model(&Like{}).Where("user_id = ? & post_id IN ?", userID, postIDs).Find(&likes)
	m := map[int64]struct{}{} // setの代用
	for _, like := range likes {
		m[like.PostID] = struct{}{}
	}
	for _, postID := range postIDs {
		_, ok := m[postID]
		if ok {
			flags = append(flags, true)
		} else {
			flags = append(flags, false)
		}
	}
	return flags
}

func GetLikeNum(db *gorm.DB, postID int64) int64 {
	var count int64
	db.Model(&Like{}).Where("post_id = ?", postID).Count(&count)
	return count
}

func BulkGetLikeNum(db *gorm.DB, postIDs []int64) (nums []int64) {
	var postIDsAndNums []struct {
		PostID int64
		Count  int64
	}
	db.Model(&Like{}).Select("post_id, COUNT(*) AS count").Where("post_id IN ?", postIDs).Group("post_id").Find(&postIDsAndNums)
	m := map[int64]int64{}
	for _, postIDAndNum := range postIDsAndNums {
		m[postIDAndNum.PostID] = postIDAndNum.Count
	}
	for _, postID := range postIDs {
		nums = append(nums, m[postID])
	}
	return nums
}
