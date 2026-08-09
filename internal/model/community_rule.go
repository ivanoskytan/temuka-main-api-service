package model

import (
	"time"

	"gorm.io/gorm"
)

type CommunityRule struct {
	gorm.Model
	ID          int `gorm:"primary_key;column:id"`
	CommunityID int `gorm:"column:community_id"`

	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}
