package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID        int       `gorm:"primary_key;column:id"`
	UserID    int       `gorm:"column:user_id"`
	Text      string    `gorm:"column:text"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Message) TableName() string {
	return "messages"
}
