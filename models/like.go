package models

type Like struct {
	ID     int64
	UserID int64 `gorm:"index:idx_like,unique"`
	User   User  `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
	PostID int64 `gorm:"index:idx_like,unique"`
	Post   Post  `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
}
