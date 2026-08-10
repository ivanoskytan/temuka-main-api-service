package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
)

type ModeratorService interface {
	GetModeratorsByCommunityID(ctx context.Context, communityID int) (*dto.ModeratorListResponse, error)
	SendModeratorRequest(ctx context.Context, data dto.SendModeratorRequest) error
	RemoveModerator(ctx context.Context, moderatorID int) error
}

type ModeratorServiceImpl struct {
	ModeratorRepository    repository.ModeratorRepository
	NotificationRepository repository.NotificationRepository
}

func NewModeratorService(moderatorRepo repository.ModeratorRepository, notificationRepo repository.NotificationRepository) ModeratorService {
	return &ModeratorServiceImpl{
		ModeratorRepository:    moderatorRepo,
		NotificationRepository: notificationRepo,
	}
}

func (s *ModeratorServiceImpl) GetModeratorsByCommunityID(ctx context.Context, communityID int) (*dto.ModeratorListResponse, error) {
	moderators, err := s.ModeratorRepository.GetModeratorsByCommunityID(ctx, communityID)
	moderatorList := make([]dto.ModeratorDetail, 0)

	if err != nil {
		for _, m := range moderators {
			moderatorList = append(moderatorList, dto.ModeratorDetail{
				ID:             m.ID,
				UserID:         m.CommunityMember.UserID,
				Username:       m.CommunityMember.User.Username,
				ProfilePicture: m.CommunityMember.User.ProfilePicture,
				SocialPoint:    m.CommunityMember.User.SocialPoint,
				CreatedAt:      m.CreatedAt,
			})
		}
	}
	return &dto.ModeratorListResponse{
		ModeratorList: moderatorList,
	}, nil
}

func (s *ModeratorServiceImpl) SendModeratorRequest(ctx context.Context, data dto.SendModeratorRequest) error {
	notification := model.Notification{
		UserID:  data.CommunityMemberID,
		Type:    "request",
		Message: "You have been requested to be a moderator in community with ID " + strconv.Itoa(data.CommunityID),
	}

	if err := s.NotificationRepository.CreateNotification(ctx, &notification); err != nil {
		return errors.New("error creating moderator notification")
	}

	return nil
}

func (s *ModeratorServiceImpl) RemoveModerator(ctx context.Context, moderatorID int) error {
	if err := s.ModeratorRepository.DeleteModerator(ctx, moderatorID); err != nil {
		return errors.New("error removing moderator")
	}
	return nil
}
