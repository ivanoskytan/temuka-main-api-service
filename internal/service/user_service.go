package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/temuka-api-service/internal/constant"
	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/publisher"
	"github.com/temuka-api-service/internal/repository"
)

type UserService interface {
	SearchUsers(ctx context.Context, data dto.SearchUsersDTO) ([]model.User, error)
	GetUserDetail(ctx context.Context, data dto.GetUserDetailDTO) (*model.User, error)
	CreateUser(ctx context.Context, data dto.CreateUserDTO) (*model.User, error)
	UpdateUser(ctx context.Context, data dto.UpdateUserDTO) error
	FollowUser(ctx context.Context, data dto.FollowUserDTO) error
	GetFollowers(ctx context.Context, data dto.GetFollowersDTO) ([]model.UserFollow, error)
}

type UserServiceImpl struct {
	UserRepository      repository.UserRepository
	SuggestionPublisher publisher.SuggestionPublisher
}

func NewUserService(userRepository repository.UserRepository, suggestionPublisher publisher.SuggestionPublisher) UserService {
	return &UserServiceImpl{
		UserRepository:      userRepository,
		SuggestionPublisher: suggestionPublisher,
	}
}

func (s *UserServiceImpl) SearchUsers(ctx context.Context, data dto.SearchUsersDTO) ([]model.User, error) {
	users, err := s.UserRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []model.User
	for _, user := range users {
		if data.Name == "" || strings.Contains(strings.ToLower(user.Username), strings.ToLower(data.Name)) {
			filtered = append(filtered, user)
		}
	}

	if len(filtered) == 0 {
		return nil, errors.New("no users found")
	}

	return filtered, nil
}

func (s *UserServiceImpl) GetUserDetail(ctx context.Context, data dto.GetUserDetailDTO) (*model.User, error) {
	user, err := s.UserRepository.GetUserByID(ctx, data.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, data dto.CreateUserDTO) (*model.User, error) {
	newUser := model.User{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
	}

	if err := s.UserRepository.CreateUser(ctx, &newUser); err != nil {
		return nil, errors.New("error creating user")
	}

	return &newUser, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, data dto.UpdateUserDTO) error {
	updatedUser := model.User{
		Username:       data.Username,
		Desc:           data.Desc,
		Displayname:    data.Displayname,
		ProfilePicture: data.ProfilePicture,
	}

	if err := s.UserRepository.UpdateUser(ctx, data.UserID, &updatedUser); err != nil {
		return errors.New("error updating user")
	}

	return nil
}

func (s *UserServiceImpl) FollowUser(ctx context.Context, data dto.FollowUserDTO) error {
	if _, err := s.UserRepository.GetUserByID(ctx, data.TargetID); err != nil {
		return errors.New("target user not found")
	}

	newFollow := model.UserFollow{
		FollowerID:  data.CurrentUserID,
		FollowingID: data.TargetID,
	}

	if err := s.UserRepository.CreateUserFollow(ctx, &newFollow); err != nil {
		return errors.New("error following user")
	}

	go s.SuggestionPublisher.PublishSuggestionEvent(
		constant.EventOperationUpdate,
		constant.EventEntityTypeUser,
		fmt.Sprintf("%d", data.TargetID),
		map[string]interface{}{},
	)

	return nil
}

func (s *UserServiceImpl) GetFollowers(ctx context.Context, data dto.GetFollowersDTO) ([]model.UserFollow, error) {
	if _, err := s.UserRepository.GetUserByID(ctx, data.UserID); err != nil {
		return nil, errors.New("user not found")
	}

	followers, err := s.UserRepository.GetFollowers(ctx, data.UserID)
	if err != nil {
		return nil, errors.New("error retrieving followers")
	}

	return followers, nil
}
