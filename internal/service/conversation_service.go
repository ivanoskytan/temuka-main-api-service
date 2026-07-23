package service

import (
	"context"
	"errors"

	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
)

type ConversationService interface {
	AddConversation(ctx context.Context, req dto.AddConversationRequest) (*model.Conversation, error)
	AddMessage(ctx context.Context, req dto.AddMessageRequest) (*model.Message, error)
	AddParticipant(ctx context.Context, req dto.AddParticipantRequest) error
	GetConversationsByUserID(ctx context.Context, userID int) ([]model.Conversation, error)
	GetConversationDetail(ctx context.Context, id int) (*model.Conversation, error)
	DeleteConversation(ctx context.Context, id int) error
	RetrieveMessages(ctx context.Context, conversationID int) ([]model.Message, error)
}

type ConversationServiceImpl struct {
	ConversationRepository repository.ConversationRepository
	UserRepository         repository.UserRepository
}

func NewConversationService(conversationRepo repository.ConversationRepository, userRepo repository.UserRepository) ConversationService {
	return &ConversationServiceImpl{
		ConversationRepository: conversationRepo,
		UserRepository:         userRepo,
	}
}

func (s *ConversationServiceImpl) AddConversation(ctx context.Context, req dto.AddConversationRequest) (*model.Conversation, error) {
	conversation := model.Conversation{
		UserID: req.UserID,
		Title:  req.Title,
	}
	if err := s.ConversationRepository.CreateConversation(ctx, &conversation); err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (s *ConversationServiceImpl) AddMessage(ctx context.Context, req dto.AddMessageRequest) (*model.Message, error) {
	conv, err := s.ConversationRepository.GetConversationDetailByID(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}

	var participantID int
	for _, p := range conv.Participants {
		if p.UserID == req.UserID {
			participantID = p.ID
			break
		}
	}

	if participantID == 0 {
		return nil, errors.New("user is not a participant in the conversation")
	}

	msg := &model.Message{
		UserID: participantID,
		Text:   req.Text,
	}

	if err := s.ConversationRepository.AddMessage(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *ConversationServiceImpl) AddParticipant(ctx context.Context, req dto.AddParticipantRequest) error {
	participant := model.Participant{
		UserID:         req.UserID,
		ConversationID: req.ConversationID,
	}
	return s.ConversationRepository.AddParticipant(ctx, &participant)
}

func (s *ConversationServiceImpl) GetConversationsByUserID(ctx context.Context, userID int) ([]model.Conversation, error) {
	return s.ConversationRepository.GetConversationsByUserID(ctx, userID)
}

func (s *ConversationServiceImpl) GetConversationDetail(ctx context.Context, id int) (*model.Conversation, error) {
	return s.ConversationRepository.GetConversationDetailByID(ctx, id)
}

func (s *ConversationServiceImpl) DeleteConversation(ctx context.Context, id int) error {
	return s.ConversationRepository.DeleteConversation(ctx, id)
}

func (s *ConversationServiceImpl) RetrieveMessages(ctx context.Context, conversationID int) ([]model.Message, error) {
	return s.ConversationRepository.GetMessagesByConversationID(ctx, conversationID)
}
