package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	database "github.com/temuka-api-service/util/database"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPostDetailByID(ctx context.Context, id int) (*model.Post, error)
	GetPostsByUserID(ctx context.Context, userId int) ([]model.Post, error)
	UpdatePost(ctx context.Context, id int, post *model.Post) error
	DeletePost(ctx context.Context, id int) error
}

type PostRepositoryImpl struct {
	db database.PostgresWrapper
}

func NewPostRepository(db database.PostgresWrapper) PostRepository {
	return &PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) CreatePost(ctx context.Context, post *model.Post) error {
	if err := r.db.DB.WithContext(ctx).Create(post).Error; err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

func (r *PostRepositoryImpl) GetPostDetailByID(ctx context.Context, id int) (*model.Post, error) {
	var post model.Post

	if err := r.db.DB.WithContext(ctx).
		Preload("User").
		Preload("Likes").
		Preload("CommunityPost").
		Preload("CommunityPost.Community").
		First(&post, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get post detail: %w", err)
	}

	return &post, nil
}

func (r *PostRepositoryImpl) DeletePost(ctx context.Context, id int) error {
	if err := r.db.Delete(ctx, &model.Post{}, id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

func (r *PostRepositoryImpl) UpdatePost(ctx context.Context, id int, post *model.Post) error {
	q := r.db.DB.WithContext(ctx).Model(&model.Post{}).Where("id = ?", id)

	if err := q.Updates(post).Error; err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

func (r *PostRepositoryImpl) GetPostsByUserID(ctx context.Context, userId int) ([]model.Post, error) {
	var posts []model.Post

	q := r.db.DB.WithContext(ctx).Where("user_id = ?", userId)

	if err := q.Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to get posts by user id: %w", err)
	}

	return posts, nil
}
