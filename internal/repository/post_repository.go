package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	database "github.com/temuka-api-service/util/database"
	"gorm.io/gorm"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPostDetailByID(ctx context.Context, id int) (*model.Post, error)
	GetPostsByUserId(ctx context.Context, userId int) ([]model.Post, error)
	UpdatePost(ctx context.Context, id int, post *model.Post) error
	DeletePost(ctx context.Context, id int) error
	HasUserLikedPost(ctx context.Context, postId, userId int) (bool, error)
	GetLikedPostIdsForUser(ctx context.Context, userId int, postIds []int) (map[int]bool, error)
	CreatePostLike(ctx context.Context, postLike *model.PostLike) error
	DeletePostLike(ctx context.Context, userId, postId int) error
	IncrementPostLikeCount(ctx context.Context, postId int) error
	DecrementPostLikeCount(ctx context.Context, postId int) error
	CreatePostSave(ctx context.Context, postSave *model.PostSave) error
	DeletePostSave(ctx context.Context, postId, userId int) error
	HasUserSavedPost(ctx context.Context, postId, userId int) (bool, error)
	GetSavedPostIdsForUser(ctx context.Context, userId int, postIds []int) (map[int]bool, error)
	GetSavedPostsByUser(ctx context.Context, userId int) ([]model.Post, error)
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

func (r *PostRepositoryImpl) GetPostsByUserId(ctx context.Context, userId int) ([]model.Post, error) {
	var posts []model.Post

	q := r.db.DB.WithContext(ctx).Where("user_id = ?", userId)

	if err := q.Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to get posts by user id: %w", err)
	}

	return posts, nil
}

func (r *PostRepositoryImpl) HasUserLikedPost(ctx context.Context, postId, userId int) (bool, error) {
	var count int64
	err := r.db.DB.WithContext(ctx).
		Model(&model.PostLike{}).
		Where("post_id = ? AND user_id = ?", postId, userId).
		Count(&count).Error

	return count > 0, err
}

func (r *PostRepositoryImpl) GetLikedPostIdsForUser(ctx context.Context, userId int, postIds []int) (map[int]bool, error) {
	if len(postIds) == 0 {
		return map[int]bool{}, nil
	}

	var likedIDs []int
	err := r.db.DB.WithContext(ctx).
		Model(&model.PostLike{}).
		Where("user_id = ? AND post_id IN ?", userId, postIds).
		Pluck("post_id", &likedIDs).Error

	if err != nil {
		return nil, err
	}

	likedMap := make(map[int]bool, len(likedIDs))
	for _, id := range likedIDs {
		likedMap[id] = true
	}

	return likedMap, nil
}

func (r *PostRepositoryImpl) CreatePostLike(ctx context.Context, postLike *model.PostLike) error {
	return r.db.DB.WithContext(ctx).Create(postLike).Error
}

func (r *PostRepositoryImpl) DeletePostLike(ctx context.Context, postId, userId int) error {
	return r.db.DB.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postId, userId).
		Delete(&model.PostLike{}).Error
}

func (r *PostRepositoryImpl) IncrementPostLikeCount(ctx context.Context, postId int) error {
	return r.db.DB.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", postId).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

func (r *PostRepositoryImpl) DecrementPostLikeCount(ctx context.Context, postId int) error {
	return r.db.DB.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", postId).
		UpdateColumn("like_count", gorm.Expr("GREATEST(0, like_count - 1)")).Error
}

func (r *PostRepositoryImpl) CreatePostSave(ctx context.Context, postSave *model.PostSave) error {
	return r.db.DB.WithContext(ctx).Create(postSave).Error
}

func (r *PostRepositoryImpl) DeletePostSave(ctx context.Context, postId, userId int) error {
	return r.db.DB.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postId, userId).
		Delete(&model.PostSave{}).Error
}

func (r *PostRepositoryImpl) HasUserSavedPost(ctx context.Context, postId, userId int) (bool, error) {
	var count int64
	err := r.db.DB.WithContext(ctx).
		Model(&model.PostSave{}).
		Where("post_id = ? AND user_id = ?", postId, userId).
		Count(&count).Error

	return count > 0, err
}

func (r *PostRepositoryImpl) GetSavedPostsByUser(ctx context.Context, userId int) ([]model.Post, error) {
	var posts []model.Post

	err := r.db.DB.WithContext(ctx).
		Joins("JOIN post_saves ON post_saves.post_id = posts.id").
		Where("post_saves.user_id = ?", userId).
		Find(&posts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get saved posts by user id: %w", err)
	}

	return posts, nil
}

func (r *PostRepositoryImpl) GetSavedPostIdsForUser(ctx context.Context, userId int, postIds []int) (map[int]bool, error) {
	if len(postIds) == 0 {
		return map[int]bool{}, nil
	}

	var savedIDs []int
	err := r.db.DB.WithContext(ctx).
		Model(&model.PostSave{}).
		Where("user_id = ? AND post_id IN ?", userId, postIds).
		Pluck("post_id", &savedIDs).Error

	if err != nil {
		return nil, err
	}

	savedMap := make(map[int]bool, len(savedIDs))
	for _, id := range savedIDs {
		savedMap[id] = true
	}

	return savedMap, nil
}
