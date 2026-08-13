package service

import (
	"context"
	"errors"

	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
)

type CommentService interface {
	AddComment(ctx context.Context, data dto.AddCommentRequest) (*model.Comment, error)
	ShowCommentsByPost(ctx context.Context, data dto.ShowCommentsRequest) ([]model.Comment, error)
	DeleteComment(ctx context.Context, commentID int) error
	ShowReplies(ctx context.Context, data dto.ShowRepliesRequest) ([]model.Comment, error)
	GetCommentsByUser(ctx context.Context, userID int) ([]model.Comment, error)
}

type CommentServiceImpl struct {
	CommentRepository      repository.CommentRepository
	PostRepository         repository.PostRepository
	NotificationRepository repository.NotificationRepository
	ReportRepository       repository.ReportRepository
}

func NewCommentService(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	notificationRepo repository.NotificationRepository,
	reportRepo repository.ReportRepository,
) CommentService {
	return &CommentServiceImpl{
		CommentRepository:      commentRepo,
		PostRepository:         postRepo,
		NotificationRepository: notificationRepo,
		ReportRepository:       reportRepo,
	}
}

func (s *CommentServiceImpl) AddComment(ctx context.Context, data dto.AddCommentRequest) (*model.Comment, error) {
	post, err := s.PostRepository.GetPostDetailByID(ctx, data.PostID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	newComment := model.Comment{
		UserID:   data.UserID,
		PostID:   data.PostID,
		ParentID: data.ParentID,
		Content:  data.Content,
	}

	if err := s.CommentRepository.CreateComment(ctx, &newComment); err != nil {
		return nil, errors.New("error creating comment")
	}

	// Create notification if commenter isn’t post owner
	if post.UserID != data.UserID {
		newNotification := model.Notification{
			UserID:    post.UserID,
			ActorID:   data.UserID,
			PostID:    &data.PostID,
			CommentID: &newComment.ID,
			Type:      "comment",
			Message:   "New comment on your post",
			Read:      false,
		}

		if err := s.NotificationRepository.CreateNotification(ctx, &newNotification); err != nil {
			return nil, errors.New("error creating notification")
		}
	}

	return &newComment, nil
}

func (s *CommentServiceImpl) ShowCommentsByPost(ctx context.Context, data dto.ShowCommentsRequest) ([]model.Comment, error) {
	comments, err := s.CommentRepository.GetCommentsByPostID(ctx, data.PostID)
	if err != nil {
		return nil, errors.New("error retrieving comments")
	}
	return comments, nil
}

func (s *CommentServiceImpl) GetCommentsByUser(ctx context.Context, userID int) ([]model.Comment, error) {
	comments, err := s.CommentRepository.GetCommentsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("error retrieving comments by user")
	}
	return comments, nil
}

func (s *CommentServiceImpl) DeleteComment(ctx context.Context, commentID int) error {
	if err := s.CommentRepository.DeleteComment(ctx, commentID); err != nil {
		return errors.New("error deleting comment")
	}
	return nil
}

func (s *CommentServiceImpl) ShowReplies(ctx context.Context, data dto.ShowRepliesRequest) ([]model.Comment, error) {
	var fetchReplies func(parentID int) ([]model.Comment, error)
	fetchReplies = func(parentID int) ([]model.Comment, error) {
		comments, err := s.CommentRepository.GetRepliesByParentID(ctx, parentID)
		if err != nil {
			return nil, err
		}

		for i := range comments {
			replies, err := fetchReplies(comments[i].ID)
			if err != nil {
				return nil, err
			}
			comments[i].Replies = replies
		}
		return comments, nil
	}

	return fetchReplies(data.ParentID)
}
