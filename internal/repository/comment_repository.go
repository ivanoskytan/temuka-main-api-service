package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	database "github.com/temuka-api-service/util/database"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetCommentsByPostID(ctx context.Context, postID int) ([]model.Comment, error)
	DeleteComment(ctx context.Context, commentID int) error
	GetRepliesByParentID(ctx context.Context, parentID int) ([]model.Comment, error)
	GetCommentDetailByID(ctx context.Context, id int) (*model.Comment, error)
	GetCommentsByUserID(ctx context.Context, userID int) ([]model.Comment, error)
}

type CommentRepositoryImpl struct {
	db database.PostgresWrapper
}

func NewCommentRepository(db database.PostgresWrapper) CommentRepository {
	return &CommentRepositoryImpl{
		db: db,
	}
}

func (r *CommentRepositoryImpl) CreateComment(ctx context.Context, comment *model.Comment) error {
	err := r.db.DB.WithContext(ctx).Create(comment).Error
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

func (r *CommentRepositoryImpl) GetCommentsByPostID(ctx context.Context, postID int) ([]model.Comment, error) {
	var comments []model.Comment

	db := r.db.DB.WithContext(ctx).Preload("User").Where("post_id = ?", postID)
	if err := db.Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, nil
}

func (r *CommentRepositoryImpl) GetCommentsByUserID(ctx context.Context, userID int) ([]model.Comment, error) {
	var comments []model.Comment

	db := r.db.DB.WithContext(ctx).Where("user_id = ?", userID)
	if err := db.Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to get comments by user: %w", err)
	}

	return comments, nil
}

func (r *CommentRepositoryImpl) DeleteComment(ctx context.Context, commentID int) error {
	if err := r.db.DB.WithContext(ctx).Delete(&model.Comment{}, commentID).Error; err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

func (r *CommentRepositoryImpl) GetRepliesByParentID(ctx context.Context, parentID int) ([]model.Comment, error) {
	var replies []model.Comment

	db := r.db.DB.WithContext(ctx).Where("parent_id = ?", parentID)
	if err := db.Find(&replies).Error; err != nil {
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}

	return replies, nil
}

func (r *CommentRepositoryImpl) GetCommentDetailByID(ctx context.Context, id int) (*model.Comment, error) {
	var comment model.Comment

	if err := r.db.DB.WithContext(ctx).First(&comment, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get comment detail: %w", err)
	}

	return &comment, nil
}
