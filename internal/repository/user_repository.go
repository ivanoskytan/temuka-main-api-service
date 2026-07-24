package repository

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/model"
	database "github.com/temuka-api-service/util/database"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetFollowers(ctx context.Context, userId int) ([]model.UserFollow, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, userId int, user *model.User) error
	DeleteUser(ctx context.Context, id int) error
	CreateUserFollow(ctx context.Context, userFollow *model.UserFollow) error
}

type UserRepositoryImpl struct {
	db database.PostgresWrapper
}

func NewUserRepository(db database.PostgresWrapper) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.User) error {
	if err := r.db.DB.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) CreateUserFollow(ctx context.Context, userFollow *model.UserFollow) error {
	if err := r.db.DB.WithContext(ctx).Create(userFollow).Error; err != nil {
		return fmt.Errorf("failed to create user follow: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User

	if err := r.db.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	q := r.db.DB.WithContext(ctx).Where("email = ?", email)
	if err := q.First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User

	if err := r.db.DB.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	return users, nil
}

func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, id int) error {
	if err := r.db.DB.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, userId int, user *model.User) error {
	q := r.db.DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId)

	if err := q.Updates(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepositoryImpl) GetFollowers(ctx context.Context, userId int) ([]model.UserFollow, error) {
	var followers []model.UserFollow

	q := r.db.DB.WithContext(ctx).Where("follower_id = ?", userId)
	if err := q.Find(&followers).Error; err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	return followers, nil
}
