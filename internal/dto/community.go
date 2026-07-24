package dto

import "time"

type CreateCommunityRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	LogoPicture  string `json:"logo_picture"`
	CoverPicture string `json:"cover_picture"`
}

type UpdateCommunityRequest struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	LogoPicture  string `json:"logo_picture"`
	CoverPicture string `json:"cover_picture"`
}

type JoinCommunityRequest struct {
	UserID int `json:"user_id"`
}

type GetUserJoinedCommunitiesRequest struct {
	UserID int `json:"user_id"`
}

type CommunityPostData struct {
	ID          int `json:"id"`
	PostID      int `json:"post_id"`
	CommunityID int `json:"community_id"`
	UserID      int `json:"user_id"`

	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Image        string    `json:"image"`
	Topic        string    `json:"topic"`
	Mark         string    `json:"mark"`
	UpvoteCount  int       `json:"upvote_count"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
}
