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

type UniversityService interface {
	AddUniversity(ctx context.Context, req dto.AddUniversityRequest) (*model.University, error)
	UpdateUniversity(ctx context.Context, id int, req dto.UpdateUniversityRequest) (*model.University, error)
	GetUniversityDetail(ctx context.Context, slug string) (*model.University, error)
	GetUniversities(ctx context.Context) ([]model.University, error)
	AddReview(ctx context.Context, req dto.AddReviewRequest) (*model.Review, error)
	GetUniversityReviews(ctx context.Context, universityID int) ([]dto.UniversityReviewResponse, error)
}

type UniversityServiceImpl struct {
	UniversityRepository repository.UniversityRepository
	ReviewRepository     repository.ReviewRepository
	SuggestionPublisher  publisher.SuggestionPublisher
}

func NewUniversityService(universityRepo repository.UniversityRepository, reviewRepo repository.ReviewRepository, suggestionPublisher publisher.SuggestionPublisher) UniversityService {
	return &UniversityServiceImpl{
		UniversityRepository: universityRepo,
		ReviewRepository:     reviewRepo,
		SuggestionPublisher:  suggestionPublisher,
	}
}

func (s *UniversityServiceImpl) AddUniversity(ctx context.Context, req dto.AddUniversityRequest) (*model.University, error) {
	university := model.University{
		Name:          req.Name,
		Slug:          strings.ReplaceAll(strings.ToLower(req.Name), " ", "_"),
		Summary:       req.Summary,
		LocationID:    req.LocationID,
		Website:       req.Website,
		Address:       req.Address,
		TotalMajors:   &req.TotalMajors,
		MinTuition:    req.MinTuition,
		MaxTuition:    req.MaxTuition,
		Logo:          req.Logo,
		Type:          req.Type,
		Accreditation: req.Accreditation,
	}

	if err := s.UniversityRepository.CreateUniversity(ctx, &university); err != nil {
		return nil, errors.New("failed to create university")
	}

	go s.SuggestionPublisher.PublishSuggestionEvent(
		constant.EventOperationCreate,
		constant.EventEntityTypeUniversity,
		fmt.Sprintf("%d", university.ID),
		map[string]interface{}{
			"title":   university.Name,
			"content": university.Summary,
			"icon":    university.Logo,
		},
	)

	return &university, nil
}

func (s *UniversityServiceImpl) UpdateUniversity(ctx context.Context, id int, req dto.UpdateUniversityRequest) (*model.University, error) {
	existing, err := s.UniversityRepository.GetUniversityByID(ctx, id)
	if err != nil {
		return nil, errors.New("university not found")
	}

	existing.Name = req.Name
	existing.Slug = strings.ReplaceAll(strings.ToLower(req.Name), " ", "_")
	existing.Summary = req.Summary
	existing.LocationID = req.LocationID
	existing.Website = req.Website
	existing.Address = req.Address
	existing.TotalMajors = &req.TotalMajors
	existing.MinTuition = req.MinTuition
	existing.MaxTuition = req.MaxTuition
	existing.Type = req.Type
	existing.Accreditation = req.Accreditation
	existing.Logo = req.Logo

	if err := s.UniversityRepository.UpdateUniversity(ctx, id, existing); err != nil {
		return nil, errors.New("failed to update university")
	}

	go s.SuggestionPublisher.PublishSuggestionEvent(
		constant.EventOperationUpdate,
		constant.EventEntityTypeUniversity,
		fmt.Sprintf("%d", id),
		map[string]interface{}{
			"title":   existing.Name,
			"content": existing.Summary,
			"icon":    existing.Logo,
		},
	)

	return existing, nil
}

func (s *UniversityServiceImpl) GetUniversityDetail(ctx context.Context, slug string) (*model.University, error) {
	university, err := s.UniversityRepository.GetUniversityBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("university not found")
	}
	return university, nil
}

func (s *UniversityServiceImpl) GetUniversities(ctx context.Context) ([]model.University, error) {
	return s.UniversityRepository.GetUniversityList(ctx)
}

func (s *UniversityServiceImpl) AddReview(ctx context.Context, req dto.AddReviewRequest) (*model.Review, error) {
	review := model.Review{
		UserID:       req.UserID,
		UniversityID: req.UniversityID,
		Text:         req.Text,
		Stars:        req.Stars,
	}

	if err := s.ReviewRepository.SetReview(ctx, &review); err != nil {
		return nil, errors.New("failed to create review")
	}

	university, err := s.UniversityRepository.GetUniversityByID(ctx, req.UniversityID)
	if err != nil {
		return nil, errors.New("university not found")
	}

	universityRating := 0
	if university.Rating != nil {
		universityRating = *university.Rating
	}

	universityTotalReviews := 0
	if university.TotalReviews != nil {
		universityTotalReviews = *university.TotalReviews
	}

	universityTotalReviews++
	newRating := (universityRating*universityTotalReviews + req.Stars) / universityTotalReviews

	university.Rating = &newRating
	university.TotalReviews = &universityTotalReviews

	if err := s.UniversityRepository.UpdateUniversity(ctx, req.UniversityID, university); err != nil {
		return nil, errors.New("failed to update university rating")
	}

	go s.SuggestionPublisher.PublishSuggestionEvent(
		constant.EventOperationUpdate,
		constant.EventEntityTypeUniversity,
		fmt.Sprintf("%d", university.ID),
		map[string]interface{}{},
	)

	return &review, nil
}

func (s *UniversityServiceImpl) GetUniversityReviews(ctx context.Context, universityID int) ([]dto.UniversityReviewResponse, error) {
	reviews, err := s.ReviewRepository.GetReviewsByUniversityID(ctx, universityID)
	if err != nil {
		return nil, errors.New("failed to get reviews")
	}

	var universityReviewResponse []dto.UniversityReviewResponse
	for _, review := range reviews {
		universityReviewResponse = append(universityReviewResponse, dto.UniversityReviewResponse{
			ID:             review.ID,
			UserID:         review.UserID,
			UniversityID:   review.UniversityID,
			Text:           review.Text,
			Stars:          review.Stars,
			Username:       review.User.Username,
			ProfilePicture: review.User.ProfilePicture,
			CreatedAt:      review.CreatedAt,
		})
	}

	return universityReviewResponse, nil
}
