package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID     int `gorm:"primary_key;column:id"`
	UserID int `gorm:"column:user_id"`

	User          *User          `gorm:"foreignKey:UserID"`
	CommunityPost *CommunityPost `gorm:"foreignKey:PostID"`

	IsAnonymous      bool   `gorm:"column:is_anonymous;default:false"`
	UniversityOrigin string `gorm:"column:university_origin"`

	Title         string         `gorm:"column:title"`
	Description   string         `gorm:"column:desc"`
	Image         string         `gorm:"column:image"`
	Comments      []Comment      `gorm:"foreignKey:PostID"`
	Notifications []Notification `gorm:"foreignKey:PostID"`

	Likes     []*User   `gorm:"many2many:post_likes;"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (p *Post) TableName() string {
	return "posts"
}
