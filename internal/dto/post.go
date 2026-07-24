package dto

import (
	"time"

	"github.com/temuka-api-service/internal/model"
)

type CreatePostRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	UserID      int     `json:"user_id"`
	CommunityID *int    `json:"community_id"`
	Mark        *string `json:"mark"`
	Topic       *string `json:"topic"`
}

type UpdatePostRequest struct {
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type LikePostRequest struct {
	UserID int `json:"user_id"`
}

type PostUserSummary struct {
	Username       string `json:"Username"`
	ProfilePicture string `json:"ProfilePicture"`
}

type PostCommentSummary struct {
	ID             int       `json:"ID"`
	UserID         int       `json:"UserID"`
	Username       string    `json:"Username"`
	ProfilePicture string    `json:"ProfilePicture"`
	PostID         int       `json:"PostID"`
	ParentID       *int      `json:"ParentID"`
	Votes          int       `json:"Votes"`
	Content        string    `json:"Content"`
	CreatedAt      time.Time `json:"CreatedAt"`
}

type PostDetailData struct {
	Post     *model.Post          `json:"post"`
	User     PostUserSummary      `json:"user"`
	Comments []PostCommentSummary `json:"comments"`
}

type PostCreatedEventData struct {
	PostID      int    `json:"post_id"`
	UserID      int    `json:"user_id"`
	CommunityID int    `json:"community_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PostLikedEventData struct {
	PostID        int `json:"post_id"`
	PostOwnerID   int `json:"post_owner_id"`
	LikedByUserID int `json:"liked_by_user_id"`
	CommunityID   int `json:"community_id"`
}
