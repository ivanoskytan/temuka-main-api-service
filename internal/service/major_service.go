package service

import (
	"context"
	"errors"

	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
)

type MajorService interface {
	AddMajor(ctx context.Context, req dto.AddMajorRequest) (*model.Major, error)
	GetMajors(ctx context.Context) ([]model.Major, error)
	GetMajorDetail(ctx context.Context, id int) (*model.Major, error)
	GetMajorsByUniversity(ctx context.Context, universityID int) ([]model.Major, error)
	AddMajorReview(ctx context.Context, req dto.AddMajorReviewRequest) (*model.MajorReview, error)
	GetMajorReviews(ctx context.Context, majorID int) ([]model.MajorReview, error)
}

type MajorServiceImpl struct {
	MajorRepository repository.MajorRepository
}

func NewMajorService(majorRepo repository.MajorRepository) MajorService {
	return &MajorServiceImpl{
		MajorRepository: majorRepo,
	}
}

func (s *MajorServiceImpl) AddMajor(ctx context.Context, req dto.AddMajorRequest) (*model.Major, error) {
	zeroValue := 0
	major := model.Major{
		Name:         req.Name,
		UniversityID: req.UniversityID,
		TotalReviews: &zeroValue,
		Rating:       &zeroValue,
	}

	if err := s.MajorRepository.CreateMajor(ctx, &major); err != nil {
		return nil, errors.New("failed to create major record")
	}
	return &major, nil
}

func (s *MajorServiceImpl) GetMajors(ctx context.Context) ([]model.Major, error) {
	return s.MajorRepository.GetMajorList(ctx)
}

func (s *MajorServiceImpl) GetMajorDetail(ctx context.Context, id int) (*model.Major, error) {
	return s.MajorRepository.GetMajorByID(ctx, id)
}

func (s *MajorServiceImpl) GetMajorsByUniversity(ctx context.Context, universityID int) ([]model.Major, error) {
	return s.MajorRepository.GetMajorsByUniversityID(ctx, universityID)
}

func (s *MajorServiceImpl) AddMajorReview(ctx context.Context, req dto.AddMajorReviewRequest) (*model.MajorReview, error) {
	review := model.MajorReview{
		UserID:  req.UserID,
		MajorID: req.MajorID,
		Text:    req.Text,
		Stars:   req.Stars,
	}

	if err := s.MajorRepository.SetMajorReview(ctx, &review); err != nil {
		return nil, errors.New("failed to save major review")
	}

	major, err := s.MajorRepository.GetMajorByID(ctx, req.MajorID)
	if err != nil {
		return nil, errors.New("associated major not found")
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
		return nil, errors.New("failed to update major moving averages")
	}

	return &review, nil
}

func (s *MajorServiceImpl) GetMajorReviews(ctx context.Context, majorID int) ([]model.MajorReview, error) {
	return s.MajorRepository.GetMajorReviewsByMajorID(ctx, majorID)
}
