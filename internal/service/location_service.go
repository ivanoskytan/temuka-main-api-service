package service

import (
	"context"
	"fmt"

	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
)

type LocationService interface {
	AddLocation(ctx context.Context, req *dto.AddLocationRequest) (*model.Location, error)
	UpdateLocation(ctx context.Context, id int, req *dto.UpdateLocationRequest) (*model.Location, error)
	GetLocations(ctx context.Context) ([]model.Location, error)
}

type LocationServiceImpl struct {
	locationRepo repository.LocationRepository
}

func NewLocationService(repo repository.LocationRepository) LocationService {
	return &LocationServiceImpl{
		locationRepo: repo,
	}
}

func (s *LocationServiceImpl) AddLocation(ctx context.Context, req *dto.AddLocationRequest) (*model.Location, error) {
	newLocation := model.Location{
		Name: req.Name,
	}

	if err := s.locationRepo.AddLocation(ctx, &newLocation); err != nil {
		return nil, err
	}

	return &newLocation, nil
}

func (s *LocationServiceImpl) UpdateLocation(ctx context.Context, id int, req *dto.UpdateLocationRequest) (*model.Location, error) {
	location, err := s.locationRepo.GetLocationById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("location not found")
	}

	location.Name = req.Name
	return location, nil
}

func (s *LocationServiceImpl) GetLocations(ctx context.Context) ([]model.Location, error) {
	locations, err := s.locationRepo.GetLocations(ctx)
	if err != nil {
		return nil, err
	}
	return locations, nil
}
