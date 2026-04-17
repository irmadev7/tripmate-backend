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

func (r *itineraryRepository) GetItineraryByUser(ctx context.Context, query model.PaginationQuery, userId int) (*[]model.Itinerary, int64, error) {
	var itineraries []model.Itinerary
	var total int64

	db := r.db.WithContext(ctx)

	baseQuery := r.db.WithContext(ctx).Model(&model.Itinerary{}).Where("user_id = ?", userId)
	if query.Search != "" {
		baseQuery = baseQuery.Where("title ILIKE ?", "%"+query.Search+"%")
	}
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (query.Page - 1) * query.Limit

	if err := db.Model(&model.Itinerary{}).Where("user_id = ?", userId).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := baseQuery.Preload("Places").Where("user_id = ?", userId).Limit(query.Limit).Offset(offset).Order("updated_at DESC").Find(&itineraries).Error; err != nil {
		return nil, 0, err
	}

	return &itineraries, total, nil
}

func (r *itineraryRepository) GetItineraryById(ctx context.Context, itineraryId int) (*model.Itinerary, error) {
	var itinerary model.Itinerary
	if err := r.db.Preload("Places").Where("id = ?", itineraryId).First(&itinerary).Error; err != nil {
		return nil, err
	}
	return &itinerary, nil
}
