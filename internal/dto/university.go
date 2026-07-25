package dto

import "time"

type AddUniversityRequest struct {
	Name          string `json:"name"`
	Summary       string `json:"summary"`
	LocationID    int    `json:"location_id"`
	Website       string `json:"website"`
	Address       string `json:"address"`
	MinTuition    int    `json:"min_tuition"`
	MaxTuition    int    `json:"max_tuition"`
	TotalMajors   int    `json:"total_majors"`
	Logo          string `json:"logo"`
	Type          string `json:"type"`
	Accreditation string `json:"accreditation"`
}

type UpdateUniversityRequest struct {
	Name          string `json:"name"`
	Summary       string `json:"summary"`
	LocationID    int    `json:"location_id"`
	Website       string `json:"website"`
	Address       string `json:"address"`
	MinTuition    int    `json:"min_tuition"`
	MaxTuition    int    `json:"max_tuition"`
	TotalMajors   int    `json:"total_majors"`
	Logo          string `json:"logo"`
	Type          string `json:"type"`
	Accreditation string `json:"accreditation"`
}

type AddReviewRequest struct {
	UserID       int    `json:"user_id"`
	UniversityID int    `json:"university_id"`
	Text         string `json:"text"`
	Stars        int    `json:"stars"`
}

type UniversityReviewResponse struct {
	ID           int `json:"ID"`
	UniversityID int `json:"UniversityID"`
	UserID       int `json:"UserID"`

	Username       string    `json:"Username"`
	ProfilePicture string    `json:"ProfilePicture"`
	Stars          int       `json:"Stars"`
	Text           string    `json:"Text"`
	CreatedAt      time.Time `json:"CreatedAt"`
	UpdatedAt      time.Time `json:"UpdatedAt"`
}
