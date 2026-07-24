package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	database "github.com/temuka-api-service/util/database"
)

type UniversityRepository interface {
	CreateUniversity(ctx context.Context, university *model.University) error
	UpdateUniversity(ctx context.Context, id int, university *model.University) error
	GetUniversityList(ctx context.Context) ([]model.University, error)
	DeleteUniversity(ctx context.Context, id int) error
	GetUniversityByID(ctx context.Context, id int) (*model.University, error)
	GetUniversityBySlug(ctx context.Context, slug string) (*model.University, error)
}

type UniversityRepositoryImpl struct {
	db database.PostgresWrapper
}

func NewUniversityRepository(db database.PostgresWrapper) UniversityRepository {
	return &UniversityRepositoryImpl{
		db: db,
	}
}

func (r *UniversityRepositoryImpl) CreateUniversity(ctx context.Context, university *model.University) error {
	if err := r.db.DB.WithContext(ctx).Create(university).Error; err != nil {
		return fmt.Errorf("failed to create university: %w", err)
	}
	return nil
}

func (r *UniversityRepositoryImpl) DeleteUniversity(ctx context.Context, id int) error {
	if err := r.db.DB.WithContext(ctx).Delete(&model.University{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete university: %w", err)
	}
	return nil
}

func (r *UniversityRepositoryImpl) UpdateUniversity(ctx context.Context, id int, university *model.University) error {
	q := r.db.Model(ctx, &model.University{}).Where("id = ?", id)

	if err := q.Updates(university).Error; err != nil {
		return fmt.Errorf("failed to update university: %w", err)
	}

	return nil
}

func (r *UniversityRepositoryImpl) GetUniversityList(ctx context.Context) ([]model.University, error) {
	var universities []model.University

	if err := r.db.DB.WithContext(ctx).Find(&universities).Error; err != nil {
		return nil, fmt.Errorf("failed to get university list: %w", err)
	}

	return universities, nil
}

func (r *UniversityRepositoryImpl) GetUniversityByID(ctx context.Context, id int) (*model.University, error) {
	var university model.University

	if err := r.db.DB.WithContext(ctx).First(&university, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get university by id: %w", err)
	}

	return &university, nil
}

func (r *UniversityRepositoryImpl) GetUniversityBySlug(ctx context.Context, slug string) (*model.University, error) {
	var university model.University

	err := r.db.DB.WithContext(ctx).Where("slug = ?", slug).First(&university).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get university by slug: %w", err)
	}

	return &university, nil
}
