package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	database "github.com/temuka-api-service/util/database"
)

type ReviewRepository interface {
	SetReview(ctx context.Context, review *model.Review) error
	DeleteReview(ctx context.Context, id int) error
	GetReviewsByUniversityID(ctx context.Context, universityID int) ([]model.Review, error)
}

type ReviewRepositoryImpl struct {
	db database.PostgresWrapper
}

func NewReviewRepository(db database.PostgresWrapper) ReviewRepository {
	return &ReviewRepositoryImpl{
		db: db,
	}
}

// SetReview = create/upsert
func (r *ReviewRepositoryImpl) SetReview(ctx context.Context, review *model.Review) error {
	if err := r.db.DB.WithContext(ctx).Create(review).Error; err != nil {
		return fmt.Errorf("failed to set review: %w", err)
	}
	return nil
}

func (r *ReviewRepositoryImpl) DeleteReview(ctx context.Context, id int) error {
	if err := r.db.DB.WithContext(ctx).Delete(&model.Review{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	return nil
}

func (r *ReviewRepositoryImpl) GetReviewsByUniversityID(ctx context.Context, universityID int) ([]model.Review, error) {
	var reviews []model.Review

	err := r.db.DB.WithContext(ctx).Where("university_id = ?", universityID).Find(&reviews).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	return reviews, nil
}
