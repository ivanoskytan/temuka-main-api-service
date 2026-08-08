package model

import (
	"time"

	"gorm.io/gorm"
)

type University struct {
	gorm.Model
	ID             int     `gorm:"primary_key;university_id"`
	Name           string  `gorm:"column:name"`
	Slug           string  `gorm:"column:slug"`
	Logo           string  `gorm:"column:logo"`
	Summary        string  `gorm:"column:summary"`
	LocationID     int     `gorm:"column:location_id"`
	Website        string  `gorm:"column:website"`
	Address        string  `gorm:"column:address"`
	TotalReviews   *int    `gorm:"column:total_reviews"`
	TotalMajors    *int    `gorm:"column:total_majors"`
	Rating         *int    `gorm:"column:rating"`
	Type           string  `gorm:"column:type"`
	Accreditation  string  `gorm:"column:accreditation"`
	MinTuition     int     `gorm:"column:min_tuition"`
	MaxTuition     int     `gorm:"column:max_tuition"`
	AcceptanceRate float32 `gorm:"column:acceptance_rate"`

	Reviews []Review `gorm:"foreignKey:UniversityID"`
	Majors  []Major  `gorm:"foreignKey:UniversityID"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *University) TableName() string {
	return "universities"
}
