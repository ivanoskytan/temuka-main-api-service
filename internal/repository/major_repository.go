package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/util/database"
)

type MajorRepository interface {
	CreateMajor(ctx context.Context, major *model.Major) error
	UpdateMajor(ctx context.Context, id int, major *model.Major) error
	GetMajorByID(ctx context.Context, id int) (*model.Major, error)
	GetMajorList(ctx context.Context) ([]model.Major, error)
	GetMajorsByUniversityID(ctx context.Context, universityId int) ([]model.Major, error)
	GetMajorReviewsByMajorID(ctx context.Context, majorId int) ([]model.MajorReview, error)
	SetMajorReview(ctx context.Context, review *model.MajorReview) error
}

type MajorRepositoryImpl struct {
	db database.PostgresWrapper
}

func NewMajorRepository(db database.PostgresWrapper) MajorRepository {
	return &MajorRepositoryImpl{
		db: db,
	}
}

func (r *MajorRepositoryImpl) CreateMajor(ctx context.Context, major *model.Major) error {
	if err := r.db.Create(ctx, major); err != nil {
		return fmt.Errorf("failed to create major: %w", err)
	}
	return nil
}

func (r *MajorRepositoryImpl) UpdateMajor(ctx context.Context, id int, major *model.Major) error {
	q := r.db.Model(ctx, &model.Major{}).Where("id = ?", id)
	if err := q.Updates(major).Error; err != nil {
		return fmt.Errorf("failed to update major metrics: %w", err)
	}
	return nil
}

func (r *MajorRepositoryImpl) GetMajorsByUniversityID(ctx context.Context, universityId int) ([]model.Major, error) {
	var majors []model.Major

	if err := r.db.Where(ctx, &majors, "university_id", universityId); err != nil {
		return nil, fmt.Errorf("failed to get majors: %w", err)
	}

	return majors, nil
}

func (r *MajorRepositoryImpl) GetMajorList(ctx context.Context) ([]model.Major, error) {
	var majors []model.Major
	if err := r.db.Find(ctx, &majors); err != nil {
		return nil, fmt.Errorf("failed to retrieve complete major list: %w", err)
	}
	return majors, nil
}

func (r *MajorRepositoryImpl) GetMajorByID(ctx context.Context, id int) (*model.Major, error) {
	var major model.Major
	if err := r.db.First(ctx, &major, id); err != nil {
		return nil, fmt.Errorf("failed to get major by id: %w", err)
	}
	return &major, nil
}

func (r *MajorRepositoryImpl) GetMajorReviewsByMajorID(ctx context.Context, majorId int) ([]model.MajorReview, error) {
	var majorReviews []model.MajorReview

	if err := r.db.Where(ctx, &majorReviews, "major_id", majorId); err != nil {
		return nil, fmt.Errorf("failed to get major reviews: %w", err)
	}

	return majorReviews, nil
}

func (r *MajorRepositoryImpl) SetMajorReview(ctx context.Context, review *model.MajorReview) error {
	if err := r.db.Create(ctx, review); err != nil {
		return fmt.Errorf("failed to save major review entity: %w", err)
	}
	return nil
}
