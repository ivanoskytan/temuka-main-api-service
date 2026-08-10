package model

import "time"

type PostSave struct {
	PostID    int       `gorm:"primaryKey;column:post_id"`
	UserID    int       `gorm:"primaryKey;column:user_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (PostSave) TableName() string {
	return "post_saves"
}
