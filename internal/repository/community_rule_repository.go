package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	database "github.com/temuka-api-service/util/database"
)

type CommunityRuleRepository interface {
	CreateCommunityRule(ctx context.Context, communityRule *model.CommunityRule) error
	CreateCommunityRules(ctx context.Context, communityRules []model.CommunityRule) error
	UpdateCommunityRule(ctx context.Context, communityRule *model.CommunityRule) error
	DeleteCommunityRule(ctx context.Context, communityRuleID int) error
}

type CommunityRuleRepositoryImpl struct {
	db database.PostgresWrapper
}

func NewCommunityRuleRepository(db database.PostgresWrapper) CommunityRuleRepository {
	return &CommunityRuleRepositoryImpl{
		db: db,
	}
}

func (r *CommunityRuleRepositoryImpl) CreateCommunityRule(ctx context.Context, communityRule *model.CommunityRule) error {
	err := r.db.DB.WithContext(ctx).Create(communityRule).Error
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

func (r *CommunityRuleRepositoryImpl) CreateCommunityRules(ctx context.Context, communityRules []model.CommunityRule) error {
	err := r.db.DB.WithContext(ctx).Create(&communityRules).Error
	if err != nil {
		return fmt.Errorf("failed to create comments: %w", err)
	}
	return nil
}

func (r *CommunityRuleRepositoryImpl) UpdateCommunityRule(ctx context.Context, communityRule *model.CommunityRule) error {
	err := r.db.DB.WithContext(ctx).Save(communityRule).Error
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	return nil
}

func (r *CommunityRuleRepositoryImpl) DeleteCommunityRule(ctx context.Context, communityRuleID int) error {
	if err := r.db.DB.WithContext(ctx).Delete(&model.CommunityRule{}, communityRuleID).Error; err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}
