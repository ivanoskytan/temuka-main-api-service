package model

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	ID              int  `gorm:"primary_key;column:id"`
	CommunityID     int  `gorm:column:community_id;not null`
	ReportedUserID  int  `gorm:column:reported_user_id;not null`
	CommunityRuleID *int `gorm:column:community_rule_id`

	TargetType string `gorm:"column:target_type;type:varchar(20);not null"`
	TargetID   int    `gorm:"column:target_id;not null"`

	Reason string `gorm:"column:reason;type:text;not null"`
	Status string `gorm:"column:status;type:varchar(20);not null;default:'pending'"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdateAt  time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (r *Report) TableName() string {
	return "reports"
}
