package model

import (
	"time"

	"gorm.io/gorm"
)

type CommunityPost struct {
	gorm.Model
	ID          int `gorm:"primary_key;column:id"`
	PostID      int `gorm:"column:post_id"`
	CommunityID int `gorm:"column:community_id"`

	Mark      string     `gorm:"column:mark"`
	Topic     string     `gorm:"column:topic"`
	Community *Community `gorm:"foreignKey:CommunityID;references:ID"`
	Post      *Post      `gorm:"foreignKey:PostID;references:ID"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *CommunityPost) TableName() string {
	return "community_posts"
}
