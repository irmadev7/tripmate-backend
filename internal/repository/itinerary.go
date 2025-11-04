package repository

import (
	"context"

	"github.com/irmadev7/tripmate-backend/internal/model"
	"gorm.io/gorm"
)

type itineraryRepository struct {
	db *gorm.DB
}

func NewItineraryRepository(db *gorm.DB) itineraryRepository {
	return itineraryRepository{db: db}
}

func (r *itineraryRepository) CreateItinerary(ctx context.Context, itinerary *model.Itinerary) error {
	return r.db.WithContext(ctx).Create(itinerary).Error
}

func (r *itineraryRepository) GetItineraryByUser(ctx context.Context, userId int) (*[]model.Itinerary, error) {
	var itineraries []model.Itinerary
	if err := r.db.Preload("Places").Where("user_id = ?", userId).Find(&itineraries).Error; err != nil {
		return nil, err
	}
	return &itineraries, nil
}

func (r *itineraryRepository) GetItineraryById(ctx context.Context, itineraryId int) (*model.Itinerary, error) {
	var itinerary model.Itinerary
	if err := r.db.Preload("Places").Where("id = ?", itineraryId).First(&itinerary).Error; err != nil {
		return nil, err
	}
	return &itinerary, nil
}
