package model

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	ID           int `gorm:"primary_key;column:id"`
	UserID       int `gorm:"column:user_id"`
	UniversityID int `gorm:"column:university_id"`

	Text      string    `gorm:column:text`
	Stars     int       `gorm:column:stars`
	User      *User     `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (r *Review) TableName() string {
	return "reviews"
}
