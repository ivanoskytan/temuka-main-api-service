package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID int `gorm:"primary_key;column:id"`

	Username       string    `gorm:"column:username"`
	Displayname    string    `gorm:"column:displayname"`
	Email          string    `gorm:"column:email"`
	Password       string    `gorm:"column:password"`
	ProfilePicture string    `gorm:"column:profile_picture"`
	CoverPicture   string    `gorm:"column:cover_picture"`
	SocialPoint    int       `gorm:"column:social_point"`
	Desc           string    `gorm:"column:description"`
	Country        string    `gorm:"column:country"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

	UniversityID *int `gorm:"column:university_id"`

	University       *University       `gorm:"foreignKey:UniversityID"`
	Followers        []UserFollow      `gorm:"foreignKey:FollowerID"`
	Followings       []UserFollow      `gorm:"foreignKey:FollowingID"`
	Messages         []Message         `gorm:"foreignKey:UserID"`
	Posts            []Post            `gorm:"foreignKey:UserID"`
	Comments         []Comment         `gorm:"foreignKey:UserID"`
	CommunityMembers []CommunityMember `gorm:"foreignKey:UserID"`
	Conversations    []Conversation    `gorm:"foreignKey:UserID"`
	Participants     []Participant     `gorm:"foreignKey:UserID"`
	Notifications    []Notification    `gorm:"foreignKey:UserID"`
	Reviews          []Review          `gorm:"foreignKey:UserID"`
}

func (u *User) TableName() string {
	return "users"
}
