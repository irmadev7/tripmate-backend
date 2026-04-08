package repository

import (
	"context"

	"github.com/irmadev7/tripmate-backend/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type ItineraryRepository interface {
	CreateItinerary(ctx context.Context, itinerary *model.Itinerary) error
	GetItineraryByUser(ctx context.Context, userId int) (*[]model.Itinerary, error)
	GetItineraryById(ctx context.Context, itineraryId int) (*model.Itinerary, error)
}

type PlaceRepository interface {
	AddPlaceToItinerary(ctx context.Context, place *model.Destination) error
}
