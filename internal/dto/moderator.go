package dto

import "time"

type SendModeratorRequest struct {
	CommunityID       int `json:"community_id"`
	CommunityMemberID int `json:"communitymember_id"`
}

type ModeratorDetail struct {
	ID             int       `json:"ID"`
	UserID         int       `json:"user_id"`
	Username       string    `json:"username"`
	ProfilePicture string    `json:"profile_picture"`
	SocialPoint    int       `json:"social_point"`
	CreatedAt      time.Time `json:"created_at"`
}

type ModeratorListResponse struct {
	ModeratorList []ModeratorDetail `json:"moderator_list"`
}
