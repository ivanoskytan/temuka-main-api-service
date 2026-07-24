package service

import (
	"context"
	"errors"
	"strings"

	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
)

type CommunityService interface {
	CreateCommunity(ctx context.Context, data dto.CreateCommunityRequest) (*model.Community, error)
	GetCommunities(ctx context.Context) ([]model.Community, error)
	UpdateCommunity(ctx context.Context, id int, data dto.UpdateCommunityRequest) (*model.Community, error)
	DeleteCommunity(ctx context.Context, id int) error
	JoinCommunity(ctx context.Context, id int, data dto.JoinCommunityRequest) error
	GetCommunityPosts(ctx context.Context, id int, filters map[string]interface{}) ([]dto.CommunityPostData, error)
	GetCommunityDetail(ctx context.Context, slug string) (*model.Community, error)
	GetUserJoinedCommunities(ctx context.Context, data dto.GetUserJoinedCommunitiesRequest) ([]model.Community, error)
}

type CommunityServiceImpl struct {
	CommunityRepository repository.CommunityRepository
}

func NewCommunityService(repo repository.CommunityRepository) CommunityService {
	return &CommunityServiceImpl{
		CommunityRepository: repo,
	}
}

func (s *CommunityServiceImpl) CreateCommunity(ctx context.Context, data dto.CreateCommunityRequest) (*model.Community, error) {
	if !s.CommunityRepository.CheckCommunityNameAvailability(ctx, data.Name) {
		return nil, errors.New("community with the same name already exists")
	}

	newCommunity := model.Community{
		Name:         data.Name,
		Slug:         strings.ReplaceAll(strings.ToLower(data.Name), " ", "_"),
		Description:  data.Description,
		LogoPicture:  data.LogoPicture,
		CoverPicture: data.CoverPicture,
	}

	if err := s.CommunityRepository.CreateCommunity(ctx, &newCommunity); err != nil {
		return nil, errors.New("error creating community")
	}

	return &newCommunity, nil
}

func (s *CommunityServiceImpl) GetCommunities(ctx context.Context) ([]model.Community, error) {
	communities, err := s.CommunityRepository.GetCommunities(ctx)
	if err != nil {
		return nil, errors.New("error retrieving communities")
	}
	return communities, nil
}

func (s *CommunityServiceImpl) UpdateCommunity(ctx context.Context, id int, data dto.UpdateCommunityRequest) (*model.Community, error) {
	updated := model.Community{
		Name:         data.Name,
		Slug:         data.Slug,
		Description:  data.Description,
		LogoPicture:  data.LogoPicture,
		CoverPicture: data.CoverPicture,
	}

	if err := s.CommunityRepository.UpdateCommunity(ctx, id, &updated); err != nil {
		return nil, errors.New("error updating community")
	}
	return &updated, nil
}

func (s *CommunityServiceImpl) DeleteCommunity(ctx context.Context, id int) error {
	if err := s.CommunityRepository.DeleteCommunity(ctx, id); err != nil {
		return errors.New("error deleting community")
	}
	return nil
}

func (s *CommunityServiceImpl) JoinCommunity(ctx context.Context, id int, data dto.JoinCommunityRequest) error {
	community, err := s.CommunityRepository.GetCommunityDetailByID(ctx, id)
	if err != nil {
		return errors.New("error retrieving community")
	}
	if community == nil {
		return errors.New("community not found")
	}

	existingMember, err := s.CommunityRepository.CheckMembership(ctx, id, data.UserID)
	if err != nil {
		return errors.New("error checking membership")
	}
	if existingMember != nil {
		return errors.New("user already a member of the community")
	}

	newMember := model.CommunityMember{
		UserID:      data.UserID,
		CommunityID: id,
	}

	if err := s.CommunityRepository.AddCommunityMember(ctx, &newMember); err != nil {
		return errors.New("error adding community member")
	}

	community.MembersCount++
	if err := s.CommunityRepository.UpdateCommunity(ctx, id, community); err != nil {
		return errors.New("error updating community member count")
	}

	return nil
}

func (s *CommunityServiceImpl) GetCommunityPosts(ctx context.Context, id int, filters map[string]interface{}) ([]dto.CommunityPostData, error) {
	communityPosts, err := s.CommunityRepository.GetCommunityPosts(ctx, id, filters)
	if err != nil {
		return nil, errors.New("error retrieving community posts")
	}

	var posts []dto.CommunityPostData
	for _, cp := range communityPosts {
		posts = append(posts, dto.CommunityPostData{
			ID:           cp.ID,
			PostID:       cp.Post.ID,
			CommunityID:  cp.CommunityID,
			Title:        cp.Post.Title,
			Description:  cp.Post.Description,
			Image:        cp.Post.Image,
			UserID:       cp.Post.UserID,
			Topic:        cp.Topic,
			Mark:         cp.Mark,
			UpvoteCount:  len(cp.Post.Likes),
			CommentCount: len(cp.Post.Comments),
			CreatedAt:    cp.Post.CreatedAt,
		})
	}

	return posts, nil
}

func (s *CommunityServiceImpl) GetCommunityDetail(ctx context.Context, slug string) (*model.Community, error) {
	community, err := s.CommunityRepository.GetCommunityDetailBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("error retrieving community detail")
	}
	return community, nil
}

func (s *CommunityServiceImpl) GetUserJoinedCommunities(ctx context.Context, data dto.GetUserJoinedCommunitiesRequest) ([]model.Community, error) {
	communities, err := s.CommunityRepository.GetUserJoinedCommunities(ctx, data.UserID)
	if err != nil {
		return nil, errors.New("error retrieving user communities")
	}
	return communities, nil
}
