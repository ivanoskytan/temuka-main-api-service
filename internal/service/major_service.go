package service

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/constant"
	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/publisher"
	"github.com/temuka-api-service/internal/repository"
)

type MajorService interface {
	AddMajor(ctx context.Context, req dto.AddMajorRequest) (*model.Major, error)
	GetMajorsByUniversity(ctx context.Context, universityID int) ([]model.Major, error)
	GetMajorDetail(ctx context.Context, id int) (*model.Major, error)
	GetMajors(ctx context.Context) ([]dto.MajorDetailResponse, error)
	AddMajorReview(ctx context.Context, req dto.AddMajorReviewRequest) (*model.MajorReview, error)
	GetMajorReviews(ctx context.Context, majorID int) ([]model.MajorReview, error)
}

type MajorServiceImpl struct {
	MajorRepository     repository.MajorRepository
	SuggestionPublisher publisher.SuggestionPublisher
}

func NewMajorService(majorRepo repository.MajorRepository, suggestionPublisher publisher.SuggestionPublisher) MajorService {
	return &MajorServiceImpl{
		MajorRepository:     majorRepo,
		SuggestionPublisher: suggestionPublisher,
	}
}

func (s *MajorServiceImpl) AddMajor(ctx context.Context, req dto.AddMajorRequest) (*model.Major, error) {
	zeroValue := 0
	major := model.Major{
		Name:         req.Name,
		UniversityID: req.UniversityID,
		Description:  req.Description,
		TotalReviews: &zeroValue,
		Rating:       &zeroValue,
	}

	if err := s.MajorRepository.CreateMajor(ctx, &major); err != nil {
		return nil, fmt.Errorf("failed to create major record")
	}

	go s.SuggestionPublisher.PublishSuggestionEvent(
		constant.EventOperationCreate,
		constant.EventEntityTypeMajor,
		fmt.Sprintf("%d", major.ID),
		map[string]interface{}{
			"title":   major.Name,
			"content": major.Description,
			"icon":    major.University.Logo,
			"slug":    fmt.Sprintf("major/%d", major.ID),
		},
	)
	return &major, nil
}

func (s *MajorServiceImpl) GetMajorsByUniversity(ctx context.Context, universityId int) ([]model.Major, error) {
	return s.MajorRepository.GetMajorsByUniversityID(ctx, universityId)
}

func (s *MajorServiceImpl) GetMajorDetail(ctx context.Context, id int) (*model.Major, error) {
	return s.MajorRepository.GetMajorByID(ctx, id)
}

func (s *MajorServiceImpl) GetMajors(ctx context.Context) ([]dto.MajorDetailResponse, error) {
	majors, err := s.MajorRepository.GetMajorList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve majors for the specified university")
	}

	var majorDetailResponse []dto.MajorDetailResponse
	for _, major := range majors {
		majorDetailResponse = append(majorDetailResponse, dto.MajorDetailResponse{
			ID:             major.ID,
			Name:           major.Name,
			Description:    major.Description,
			TotalReviews:   major.TotalReviews,
			Rating:         major.Rating,
			UniversityName: major.University.Name,
			UniversityLogo: major.University.Logo,
		})
	}

	return majorDetailResponse, nil
}

func (s *MajorServiceImpl) AddMajorReview(ctx context.Context, req dto.AddMajorReviewRequest) (*model.MajorReview, error) {
	review := model.MajorReview{
		UserID:  req.UserID,
		MajorID: req.MajorID,
		Text:    req.Text,
		Stars:   req.Stars,
	}

	if err := s.MajorRepository.SetMajorReview(ctx, &review); err != nil {
		return nil, fmt.Errorf("failed to save major review")
	}

	major, err := s.MajorRepository.GetMajorByID(ctx, req.MajorID)
	if err != nil {
		return nil, fmt.Errorf("associated major not found")
	}

	currentRating := 0
	if major.Rating != nil {
		currentRating = *major.Rating
	}

	currentTotalReviews := 0
	if major.TotalReviews != nil {
		currentTotalReviews = *major.TotalReviews
	}

	currentTotalReviews++
	newRating := (currentRating*(currentTotalReviews-1) + req.Stars) / currentTotalReviews

	major.Rating = &newRating
	major.TotalReviews = &currentTotalReviews

	if err := s.MajorRepository.UpdateMajor(ctx, req.MajorID, major); err != nil {
		return nil, fmt.Errorf("failed to update major moving averages")
	}

	go s.SuggestionPublisher.PublishSuggestionEvent(
		constant.EventOperationUpdate,
		constant.EventEntityTypeMajor,
		fmt.Sprintf("%d", major.ID),
		map[string]interface{}{},
	)

	return &review, nil
}

func (s *MajorServiceImpl) GetMajorReviews(ctx context.Context, majorID int) ([]model.MajorReview, error) {
	return s.MajorRepository.GetMajorReviewsByMajorID(ctx, majorID)
}
