package model

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	ID        int        `gorm:"primary_key;column:id"`
	UserID    int        `gorm:"column:user_id"`
	ActorID   int        `gorm:"column:actor_id"`
	PostID    *int       `gorm:"column:post_id;default:null"`
	CommentID *int       `gorm:"column:comment_id;default:null"`
	Type      string     `gorm:"column:type"`
	Read      bool       `gorm:"column:read;default:false"`
	Message   string     `gorm:"column:message"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	ReadedAt  *time.Time `gorm:"column:readed_at;default:null"`
}

func (n *Notification) TableName() string {
	return "notifications"
}
