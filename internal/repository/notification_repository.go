package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	database "github.com/temuka-api-service/util/database"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *model.Notification) error
	GetNotificationsByUserID(ctx context.Context, userId int) ([]model.Notification, error)
}

type NotificationRepositoryImpl struct {
	db database.PostgresWrapper
}

func NewNotificationRepository(db database.PostgresWrapper) NotificationRepository {
	return &NotificationRepositoryImpl{
		db: db,
	}
}

func (r *NotificationRepositoryImpl) CreateNotification(ctx context.Context, notification *model.Notification) error {
	if err := r.db.DB.WithContext(ctx).Create(notification).Error; err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

func (r *NotificationRepositoryImpl) GetNotificationsByUserID(ctx context.Context, userId int) ([]model.Notification, error) {
	var notifications []model.Notification

	if err := r.db.DB.WithContext(ctx).Where(&notifications, "user_id = ?", userId).Error; err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return notifications, nil
}
