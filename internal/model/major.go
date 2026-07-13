package model

import (
	"time"

	"gorm.io/gorm"
)

type Major struct {
	gorm.Model
	ID           int           `gorm:"primary_key;university_id"`
	UniversityID int           `gorm:"column:university_id"`
	Name         string        `gorm:"column:name"`
	Description  string        `gorm:"column:description"`
	TotalReviews *int          `gorm:"column:total_reviews"`
	Rating       *int          `gorm:"column:rating"`
	Reviews      []MajorReview `gorm:"foreignKey:MajorID"`
	CreatedAt    time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time     `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (m *Major) TableName() string {
	return "majors"
}
