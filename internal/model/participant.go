package model

import (
	"time"

	"gorm.io/gorm"
)

type Participant struct {
	gorm.Model
	ID             int       `gorm:"primary_key;column:id"`
	ConversationID int       `gorm:"column:conversation_id"`
	UserID         int       `gorm:"column:user_id"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Participant) TableName() string {
	return "participants"
}
