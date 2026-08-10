package model

import (
	"time"

	"gorm.io/gorm"
)

type PostLike struct {
	gorm.Model
	PostID    int       `gorm:"primary_key;column:post_id"`
	UserID    int       `gorm:"primary_key;column:user_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (p *PostLike) TableName() string {
	return "post_likes"
}
