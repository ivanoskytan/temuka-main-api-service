package service

import (
	"context"
	"errors"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/publisher"
	"github.com/temuka-api-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, data dto.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, data dto.LoginRequest) (map[string]interface{}, error)
	ResetPassword(ctx context.Context, data dto.ResetPasswordRequest) error
}

type AuthServiceImpl struct {
	UserRepository       repository.UserRepository
	searchIndexPublisher publisher.SearchIndexPublisher
}

func NewAuthService(userRepository repository.UserRepository, searchIndexPublisher publisher.SearchIndexPublisher) AuthService {
	return &AuthServiceImpl{
		UserRepository: userRepository,
	}
}

func (c *AuthServiceImpl) Register(ctx context.Context, data dto.RegisterRequest) (*model.User, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error hashing password")
	}

	newUser := model.User{
		Username:       data.Username,
		Email:          data.Email,
		Password:       string(hashedPwd),
		ProfilePicture: "",
		CoverPicture:   "",
	}

	if err := c.UserRepository.CreateUser(ctx, &newUser); err != nil {
		return nil, errors.New("error creating user")
	}

	return &newUser, nil
}

func (c *AuthServiceImpl) Login(ctx context.Context, data dto.LoginRequest) (map[string]interface{}, error) {
	user, err := c.UserRepository.GetUserByEmail(ctx, data.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})

	tokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return nil, errors.New("error generating token")
	}

	response := map[string]interface{}{
		"token": tokenString,
	}

	return response, nil
}

func (c *AuthServiceImpl) ResetPassword(ctx context.Context, data dto.ResetPasswordRequest) error {
	token, err := jwt.Parse(data.ResetToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["email"] != data.Email {
			return errors.New("invalid token email")
		}

		if data.NewPassword != data.NewPasswordConfirmation {
			return errors.New("passwords do not match")
		}

		hashedNewPwd, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("error hashing new password")
		}

		userID, err := strconv.Atoi(data.UserID)
		if err != nil {
			return errors.New("invalid user ID")
		}

		user, err := c.UserRepository.GetUserByID(ctx, userID)
		if err != nil {
			return errors.New("user not found")
		}

		user.Password = string(hashedNewPwd)
		if err := c.UserRepository.UpdateUser(ctx, userID, user); err != nil {
			return errors.New("error updating password")
		}

		return nil
	}

	return errors.New("invalid token")
}
