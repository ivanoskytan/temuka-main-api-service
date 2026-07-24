package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	database "github.com/temuka-api-service/util/database"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *model.Message) error
	DeleteMessage(ctx context.Context, id int) error
}

type MessageRepositoryImpl struct {
	db database.PostgresWrapper
}

func NewMessageRepositoryImpl(db database.PostgresWrapper) MessageRepository {
	return &MessageRepositoryImpl{
		db: db,
	}
}

func (r *MessageRepositoryImpl) CreateMessage(ctx context.Context, message *model.Message) error {
	if err := r.db.DB.WithContext(ctx).Create(message).Error; err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

func (r *MessageRepositoryImpl) DeleteMessage(ctx context.Context, id int) error {
	if err := r.db.DB.WithContext(ctx).Delete(&model.Message{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}
