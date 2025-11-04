package repository

import (
	"context"

	"github.com/irmadev7/tripmate-backend/internal/model"
	"gorm.io/gorm"
)

type placeRepository struct {
	db *gorm.DB
}

func NewPlaceRepository(db *gorm.DB) placeRepository {
	return placeRepository{db: db}
}

func (r *placeRepository) AddPlaceToItinerary(ctx context.Context, place *model.Destination) error {
	return r.db.WithContext(ctx).Create(place).Error
}
