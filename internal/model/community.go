package model

import (
	"time"

	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	ID           int    `gorm:"primary_key;column:id"`
	Name         string `gorm:"column:name"`
	Slug         string `gorm:"column:slug"`
	Description  string `gorm:"column:desc"`
	Rules        string `gorm:"column:rules"`
	MembersCount int    `gorm:"column:members_count"`
	PostsCount   int    `gorm:"column:posts_count"`
	LogoPicture  string `gorm:"column:logo_picture"`
	CoverPicture string `gorm:"column:cover_picture"`

	CommunityMembers []CommunityMember `gorm:"foreignKey:CommunityID"`
	Moderators       []Moderator       `gorm:"foreignKey:CommunityID"`
	CommunityPosts   []CommunityPost   `gorm:"foreignKey:CommunityID"`
	CommunityRules   []CommunityRule   `gorm:"foreignKey:CommunityID"`
	Reports          []Report          `gorm:"foreignKey:CommunityID"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Community) TableName() string {
	return "communities"
}
