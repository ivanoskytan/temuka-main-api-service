package dto

type AddMajorRequest struct {
	Name         string `json:"name"`
	UniversityID int    `json:"university_id"`
	Description  string `json:"description"`
}

type AddMajorReviewRequest struct {
	UserID  int    `json:"user_id"`
	MajorID int    `json:"major_id"`
	Text    string `json:"text"`
	Stars   int    `json:"stars"`
}
