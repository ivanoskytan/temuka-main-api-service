package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	database "github.com/temuka-api-service/util/database"
	"gorm.io/gorm"
)

type CommunityRepository interface {
	CreateCommunity(ctx context.Context, community *model.Community) error
	CheckCommunityNameAvailability(ctx context.Context, name string) bool
	UpdateCommunity(ctx context.Context, id int, community *model.Community) error
	GetCommunities(ctx context.Context) ([]model.Community, error)
	GetUserJoinedCommunities(ctx context.Context, userID int) ([]model.Community, error)
	GetCommunityDetailByID(ctx context.Context, id int) (*model.Community, error)
	CheckMembership(ctx context.Context, communityID, userID int) (*model.CommunityMember, error)
	AddCommunityMember(ctx context.Context, member *model.CommunityMember) error
	GetCommunityPosts(ctx context.Context, id int, filters map[string]interface{}) ([]model.CommunityPost, error)
	UpdateCommunityPostsCount(ctx context.Context, id int) error
	UpdateCommunityMembersCount(ctx context.Context, id int) error
	DeleteCommunity(ctx context.Context, id int) error
	GetCommunityDetailBySlug(ctx context.Context, slug string) (*model.Community, error)
	CreateCommunityPost(ctx context.Context, communityPost *model.CommunityPost) error
}

type CommunityRepositoryImpl struct {
	db database.PostgresWrapper
}

func NewCommunityRepository(db database.PostgresWrapper) CommunityRepository {
	return &CommunityRepositoryImpl{
		db: db,
	}
}

func (r *CommunityRepositoryImpl) CreateCommunity(ctx context.Context, community *model.Community) error {
	if err := r.db.Create(ctx, community); err != nil {
		return fmt.Errorf("failed to create community: %w", err)
	}
	return nil
}

func (r *CommunityRepositoryImpl) CreateCommunityPost(ctx context.Context, communityPost *model.CommunityPost) error {
	if err := r.db.Create(ctx, communityPost); err != nil {
		return fmt.Errorf("failed to create community post: %w", err)
	}
	return nil
}

func (r *CommunityRepositoryImpl) CheckCommunityNameAvailability(ctx context.Context, name string) bool {
	var count int64

	err := r.db.Model(ctx, &model.Community{}).
		Where("name = ?", name).
		Count(&count).Error

	if err != nil {
		return false
	}

	return count == 0
}

func (r *CommunityRepositoryImpl) UpdateCommunity(ctx context.Context, id int, community *model.Community) error {
	if err := r.db.Model(ctx, &model.Community{}).
		Where("id = ?", id).
		Updates(community).Error; err != nil {
		return fmt.Errorf("failed to update community: %w", err)
	}
	return nil
}

func (r *CommunityRepositoryImpl) GetCommunityDetailByID(ctx context.Context, id int) (*model.Community, error) {
	var community model.Community
	if err := r.db.First(ctx, &community, id); err != nil {
		return nil, fmt.Errorf("failed to get community detail: %w", err)
	}
	return &community, nil
}

func (r *CommunityRepositoryImpl) GetCommunities(ctx context.Context) ([]model.Community, error) {
	var communities []model.Community
	if err := r.db.Find(ctx, &communities); err != nil {
		return nil, fmt.Errorf("failed to get communities: %w", err)
	}
	return communities, nil
}

func (r *CommunityRepositoryImpl) DeleteCommunity(ctx context.Context, id int) error {
	if err := r.db.Delete(ctx, &model.Community{}, id); err != nil {
		return fmt.Errorf("failed to delete community: %w", err)
	}
	return nil
}

func (r *CommunityRepositoryImpl) AddCommunityMember(ctx context.Context, member *model.CommunityMember) error {
	if err := r.db.Create(ctx, member); err != nil {
		return fmt.Errorf("failed to add community member: %w", err)
	}
	return nil
}

func (r *CommunityRepositoryImpl) UpdateCommunityPostsCount(ctx context.Context, id int) error {
	if err := r.db.Model(ctx, &model.Community{}).
		Where("id = ?", id).
		Update("posts_count", gorm.Expr("posts_count + 1")).Error; err != nil {
		return fmt.Errorf("failed to update community posts count: %w", err)
	}
	return nil
}

func (r *CommunityRepositoryImpl) UpdateCommunityMembersCount(ctx context.Context, id int) error {
	if err := r.db.Model(ctx, &model.Community{}).
		Where("id = ?", id).
		Update("members_count", gorm.Expr("members_count + 1")).Error; err != nil {
		return fmt.Errorf("failed to update community members count: %w", err)
	}
	return nil
}

func (r *CommunityRepositoryImpl) CheckMembership(ctx context.Context, communityID, userID int) (*model.CommunityMember, error) {
	var member model.CommunityMember

	err := r.db.Where(ctx, "community_id = ? AND user_id = ?", communityID, userID).
		First(&member).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to check membership: %w", err)
	}

	return &member, nil
}

func (r *CommunityRepositoryImpl) GetCommunityPosts(ctx context.Context, communityID int, filters map[string]interface{}) ([]model.CommunityPost, error) {
	var posts []model.CommunityPost

	query := r.db.DB.WithContext(ctx).
		Preload("Post").
		Preload("Post.Comments").
		Preload("Post.Likes").
		Where("community_id = ?", communityID)

	for key, value := range filters {
		if key == "sort" || key == "sort_by" {
			continue
		}
		query = query.Where(key+" = ?", value)
	}

	sortBy, sortExists := filters["sort_by"].(string)
	sortOrder, orderExists := filters["sort"].(string)

	if sortExists && orderExists {
		query = query.Order(sortBy + " " + sortOrder)
	} else if sortExists {
		query = query.Order(sortBy + " asc")
	} else {
		query = query.Order("created_at desc")
	}

	if err := query.Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to get community posts: %w", err)
	}

	return posts, nil
}

func (r *CommunityRepositoryImpl) GetUserJoinedCommunities(ctx context.Context, userID int) ([]model.Community, error) {
	var communities []model.Community

	rawQuery := `
		SELECT c.*
		FROM community_members cm
		INNER JOIN communities c ON cm.community_id = c.id
		WHERE cm.user_id = ? AND cm.banned = false
	`

	if err := r.db.DB.WithContext(ctx).Raw(rawQuery, userID).Scan(&communities).Error; err != nil {
		return nil, fmt.Errorf("failed to get user joined communities: %w", err)
	}

	return communities, nil
}

func (r *CommunityRepositoryImpl) GetCommunityDetailBySlug(ctx context.Context, slug string) (*model.Community, error) {
	var community model.Community
	err := r.db.Where(ctx, "slug = ?", slug).First(&community).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get community detail by slug: %w", err)
	}
	return &community, nil
}
