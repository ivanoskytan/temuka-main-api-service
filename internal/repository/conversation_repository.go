package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	database "github.com/temuka-api-service/util/database"
)

type ConversationRepository interface {
	CreateConversation(ctx context.Context, conversation *model.Conversation) error
	GetConversationsByUserID(ctx context.Context, userID int) ([]model.Conversation, error)
	DeleteConversation(ctx context.Context, id int) error
	GetConversationDetailByID(ctx context.Context, id int) (*model.Conversation, error)
	AddParticipant(ctx context.Context, participant *model.Participant) error
	AddMessage(ctx context.Context, message *model.Message) error
	GetMessagesByConversationID(ctx context.Context, conversationID int) ([]model.Message, error)
}

type ConversationRepositoryImpl struct {
	db database.PostgresWrapper
}

func NewConversationRepository(db database.PostgresWrapper) ConversationRepository {
	return &ConversationRepositoryImpl{
		db: db,
	}
}

func (r *ConversationRepositoryImpl) CreateConversation(ctx context.Context, conversation *model.Conversation) error {
	if err := r.db.DB.WithContext(ctx).Create(conversation).Error; err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}
	return nil
}

func (r *ConversationRepositoryImpl) GetConversationsByUserID(ctx context.Context, userID int) ([]model.Conversation, error) {
	var conversations []model.Conversation

	err := r.db.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&conversations).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	return conversations, nil
}

func (r *ConversationRepositoryImpl) DeleteConversation(ctx context.Context, id int) error {
	if err := r.db.DB.WithContext(ctx).Delete(&model.Conversation{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}
	return nil
}

func (r *ConversationRepositoryImpl) AddMessage(ctx context.Context, message *model.Message) error {
	if err := r.db.DB.WithContext(ctx).Create(message).Error; err != nil {
		return fmt.Errorf("failed to add message: %w", err)
	}
	return nil
}

func (r *ConversationRepositoryImpl) GetConversationDetailByID(ctx context.Context, id int) (*model.Conversation, error) {
	var conversation model.Conversation

	err := r.db.DB.WithContext(ctx).First(&conversation, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation detail: %w", err)
	}

	return &conversation, nil
}

func (r *ConversationRepositoryImpl) AddParticipant(ctx context.Context, participant *model.Participant) error {
	if err := r.db.DB.WithContext(ctx).Create(participant).Error; err != nil {
		return fmt.Errorf("failed to add participant: %w", err)
	}
	return nil
}

func (r *ConversationRepositoryImpl) GetMessagesByConversationID(ctx context.Context, conversationID int) ([]model.Message, error) {
	var messages []model.Message

	err := r.db.DB.WithContext(ctx).Where("conversation_id = ?", conversationID).Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, nil
}
