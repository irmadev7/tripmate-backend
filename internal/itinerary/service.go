package itinerary

import (
	"context"
	"errors"
	"math"

	"github.com/irmadev7/tripmate-backend/internal/model"
	"github.com/irmadev7/tripmate-backend/internal/pkg/apperror"
	"github.com/irmadev7/tripmate-backend/internal/repository"
	"gorm.io/gorm"
)

type Service struct {
	itineraryRepo repository.ItineraryRepository
	userRepo      repository.UserRepository
	placeRepo     repository.PlaceRepository
}

func NewService(itineraryRepo repository.ItineraryRepository, userRepo repository.UserRepository, placeRepo repository.PlaceRepository) *Service {
	return &Service{
		itineraryRepo: itineraryRepo,
		userRepo:      userRepo,
		placeRepo:     placeRepo,
	}
}

func (s *Service) CreateItinerary(ctx context.Context, req model.CreateItineraryRequest) error {
	loginUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil || loginUser == nil {
		return apperror.New(apperror.Unauthorized, "user doesn't have access", err)
	}

	itinerary := model.Itinerary{
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		UserID:      loginUser.ID,
	}
	if err := s.itineraryRepo.CreateItinerary(ctx, &itinerary); err != nil {
		return apperror.New(apperror.Internal, "failed to create itinerary", err)
	}

	return nil
}

func (s *Service) GetMyItineraries(ctx context.Context, req model.GetMyItineraryRequest) (*model.GetMyItineraryResponse, error) {
	loginUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil || loginUser == nil {
		return nil, apperror.New(apperror.Unauthorized, "user doesn't have access", err)
	}

	pagination := model.PaginationQuery{
		Page:   req.Page,
		Limit:  req.Limit,
		Search: req.Search,
	}
	itineraries, total, err := s.itineraryRepo.GetItineraryByUser(ctx, pagination, int(loginUser.ID))
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to get itineraries", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &model.GetMyItineraryResponse{Data: itineraries, Meta: model.Meta{
		Page:       req.Page,
		Limit:      req.Limit,
		TotalData:  total,
		TotalPages: totalPages,
	}}, nil
}

func (s *Service) AddPlaceToItinerary(ctx context.Context, req model.AddPlaceRequest) error {
	itinerary, err := s.itineraryRepo.GetItineraryById(ctx, req.ItineraryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.New(apperror.NotFound, "itinerary not found", err)
		}
		return apperror.New(apperror.Internal, "failed to get itinerary", err)
	}

	loginUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return apperror.New(apperror.Unauthorized, "user doesn't have access", err)
	}
	if itinerary.UserID != loginUser.ID {
		return apperror.New(apperror.Unauthorized, "you don't own this itinerary", err)
	}

	place := model.Destination{
		ItineraryID: itinerary.ID,
		Name:        req.Name,
		Note:        req.Note,
		Day:         req.Day,
		Order:       req.Order,
	}
	if err := s.placeRepo.AddPlaceToItinerary(ctx, &place); err != nil {
		return apperror.New(apperror.Internal, "failed to add place", err)
	}

	return nil

}
