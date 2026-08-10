package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/temuka-api-service/internal/constant"
	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/publisher"
	"github.com/temuka-api-service/internal/repository"
	"github.com/temuka-api-service/util/database"
	"github.com/temuka-api-service/util/key_value_store"
	"gorm.io/gorm"
)

type PostService interface {
	CreatePost(ctx context.Context, req *dto.CreatePostRequest) (*model.Post, error)
	GetPostDetail(ctx context.Context, postID int) (*dto.PostDetailResponse, error)
	GetUserPosts(ctx context.Context, userID int) ([]model.Post, error)
	UpdatePost(ctx context.Context, postID int, req *dto.UpdatePostRequest) (*model.Post, error)
	DeletePost(ctx context.Context, postID int) error
	GetTimelinePosts(ctx context.Context, userID int) ([]model.Post, error)
	LikePost(ctx context.Context, postID, userID int) error
	UnlikePost(ctx context.Context, postID, userID int) error
}

type PostServiceImpl struct {
	database            database.PostgresWrapper
	postRepo            repository.PostRepository
	userRepo            repository.UserRepository
	commentRepo         repository.CommentRepository
	notificationRepo    repository.NotificationRepository
	communityRepo       repository.CommunityRepository
	redis               key_value_store.RedisWrapper
	suggestionPublisher publisher.SuggestionPublisher
}

func NewPostService(
	database database.PostgresWrapper,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	commentRepo repository.CommentRepository,
	notificationRepo repository.NotificationRepository,
	communityRepo repository.CommunityRepository,
	redis key_value_store.RedisWrapper,
	suggestionPublisher publisher.SuggestionPublisher,
) PostService {
	return &PostServiceImpl{
		database:            database,
		postRepo:            postRepo,
		userRepo:            userRepo,
		commentRepo:         commentRepo,
		notificationRepo:    notificationRepo,
		communityRepo:       communityRepo,
		redis:               redis,
		suggestionPublisher: suggestionPublisher,
	}
}

func (s *PostServiceImpl) CreatePost(ctx context.Context, req *dto.CreatePostRequest) (*model.Post, error) {
	newPost := model.Post{
		Title:       req.Title,
		Description: req.Description,
		UserID:      req.UserID,
		IsAnonymous: req.IsAnonymous,
	}

	if req.IsAnonymous {
		user, err := s.userRepo.GetUserByID(ctx, req.UserID)
		if err != nil {
			return nil, errors.New("error fetching user for anonymous post")
		}

		if user.UniversityID == nil || user.University == nil {
			return nil, errors.New("user must belong to a university to post anonymously")
		}

		newPost.UniversityOrigin = user.University.Name
	}

	if err := s.postRepo.CreatePost(ctx, &newPost); err != nil {
		return nil, errors.New("error creating post")
	}

	if req.CommunityID != nil && *req.CommunityID != 0 {
		var markVal, topicVal string
		if req.Mark != nil {
			markVal = *req.Mark
		}
		if req.Topic != nil {
			topicVal = *req.Topic
		}

		communityPost := model.CommunityPost{
			PostID:      newPost.ID,
			CommunityID: *req.CommunityID,
			Mark:        markVal,
			Topic:       topicVal,
		}

		if err := s.communityRepo.CreateCommunityPost(ctx, &communityPost); err != nil {
			return nil, errors.New("error linking post to community")
		}
		if err := s.communityRepo.UpdateCommunityPostsCount(ctx, *req.CommunityID); err != nil {
			return nil, errors.New("error updating community posts count")
		}
	}

	return &newPost, nil
}

func (s *PostServiceImpl) GetPostDetail(ctx context.Context, postID int) (*dto.PostDetailResponse, error) {
	post, err := s.postRepo.GetPostDetailByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByID(ctx, post.UserID)
	userSummary := dto.PostUserSummary{}
	if err == nil && user != nil {
		userSummary.Username = user.Username
		userSummary.ProfilePicture = user.ProfilePicture
	}

	comments, err := s.commentRepo.GetCommentsByPostID(ctx, postID)
	postComments := make([]dto.PostCommentSummary, 0)

	if err == nil {
		for _, c := range comments {
			postComments = append(postComments, dto.PostCommentSummary{
				ID:             c.ID,
				UserID:         c.UserID,
				Username:       c.User.Username,
				Content:        c.Content,
				ParentID:       c.ParentID,
				PostID:         c.PostID,
				ProfilePicture: c.User.Username,
				Votes:          len(c.Votes),
				CreatedAt:      c.CreatedAt,
			})
		}
	}

	return &dto.PostDetailResponse{
		Post:     post,
		User:     userSummary,
		Comments: postComments,
	}, nil
}

func (s *PostServiceImpl) GetUserPosts(ctx context.Context, userID int) ([]model.Post, error) {
	return s.postRepo.GetPostsByUserID(ctx, userID)
}

func (s *PostServiceImpl) UpdatePost(ctx context.Context, postID int, req *dto.UpdatePostRequest) (*model.Post, error) {
	updatedPost := model.Post{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
	}

	if err := s.postRepo.UpdatePost(ctx, postID, &updatedPost); err != nil {
		return nil, errors.New("error updating post")
	}

	go s.suggestionPublisher.PublishSuggestionEvent(
		constant.EventOperationUpdate,
		constant.EventEntityTypePost,
		fmt.Sprintf("%d", updatedPost.ID),
		map[string]interface{}{
			"title":   updatedPost.Title,
			"content": updatedPost.Description,
			"user_id": updatedPost.UserID,
		},
	)

	return &updatedPost, nil
}

func (s *PostServiceImpl) DeletePost(ctx context.Context, postID int) error {
	go s.suggestionPublisher.PublishSuggestionEvent(
		constant.EventOperationUpdate,
		constant.EventEntityTypePost,
		fmt.Sprintf("%d", postID),
		nil,
	)

	return s.postRepo.DeletePost(ctx, postID)
}

func (s *PostServiceImpl) GetTimelinePosts(ctx context.Context, userID int) ([]model.Post, error) {
	cacheKey := fmt.Sprintf("timeline_posts_user_%d", userID)

	var cached struct {
		Data []model.Post `json:"data"`
	}

	if err := s.redis.Get(cacheKey, cached); err == nil {
		return cached.Data, nil
	}

	userPosts, err := s.postRepo.GetPostsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	followers, err := s.userRepo.GetFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}

	var followerPosts []model.Post
	for _, f := range followers {
		if posts, err := s.postRepo.GetPostsByUserID(ctx, f.FollowingID); err == nil {
			followerPosts = append(followerPosts, posts...)
		}
	}

	allPosts := append(userPosts, followerPosts...)
	_ = s.redis.Set(cacheKey, struct {
		Data []model.Post `json:"data"`
	}{Data: allPosts}, 10*time.Minute)

	return allPosts, nil
}

func (s *PostServiceImpl) LikePost(ctx context.Context, postID, userID int) error {
	post, err := s.postRepo.GetPostDetailByID(ctx, postID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("post not found")
		}
		return err
	}

	isLiked, err := s.postRepo.HasUserLikedPost(ctx, postID, userID)
	if err != nil {
		return err
	}
	if isLiked {
		return nil
	}

	liker, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	err = s.database.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.postRepo.CreatePostLike(txCtx, &model.PostLike{PostID: postID, UserID: userID}); err != nil {
			return err
		}

		if err := s.postRepo.IncrementPostLikeCount(txCtx, postID); err != nil {
			return err
		}

		if post.UserID != userID {
			if err := s.userRepo.IncrementSocialPoint(txCtx, post.UserID); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	notification := model.Notification{
		UserID:  post.UserID,
		ActorID: userID,
		PostID:  post.ID,
		Type:    "like",
		Message: liker.Username + " liked your post: " + post.Title,
		Read:    false,
	}

	var icon, slug string

	if post.CommunityPost.Community != nil && post.CommunityPost.Community.ID != 0 {
		icon = post.CommunityPost.Community.LogoPicture
		slug = "comm/" + post.CommunityPost.Community.Name
	} else if post.User != nil {
		icon = post.User.ProfilePicture
		slug = "user/" + post.User.Username
	}

	go s.suggestionPublisher.PublishSuggestionEvent(
		constant.EventOperationUpdate,
		constant.EventEntityTypePost,
		fmt.Sprintf("%d", postID),
		map[string]interface{}{
			"title":   post.Title,
			"content": post.Description,
			"icon":    icon,
			"slug":    slug,
		},
	)
	return s.notificationRepo.CreateNotification(ctx, &notification)
}

func (s *PostServiceImpl) UnlikePost(ctx context.Context, postID, userID int) error {
	post, err := s.postRepo.GetPostDetailByID(ctx, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	isLiked, err := s.postRepo.HasUserLikedPost(ctx, postID, userID)
	if err != nil {
		return err
	}
	if !isLiked {
		return nil
	}

	return s.database.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.postRepo.DeletePostLike(txCtx, postID, userID); err != nil {
			return err
		}

		if err := s.postRepo.DecrementPostLikeCount(txCtx, postID); err != nil {
			return err
		}

		if post.UserID != userID {
			if err := s.userRepo.DecrementSocialPoint(txCtx, post.UserID); err != nil {
				return err
			}
		}

		return nil
	})
}
