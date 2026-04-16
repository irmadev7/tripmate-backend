package model

import (
	"time"
)

type Itinerary struct {
	BaseModel
	Title       string        `json:"title"`
	Description string        `json:"description"`
	StartDate   time.Time     `json:"start_date"`
	EndDate     time.Time     `json:"end_date"`
	UserID      uint          `json:"user_id"`
	Places      []Destination `json:"places" gorm:"foreignKey:ItineraryID"`
}

type Destination struct {
	BaseModel
	ItineraryID uint   `json:"itinerary_id"`
	Name        string `json:"name"`
	Note        string `json:"note"`
	Day         int    `json:"day"`
	Order       int    `json:"order"` // order by day
}

type CreateItineraryRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Email       string    `json:"-"`
}

type AddPlaceRequest struct {
	Name        string `json:"name" binding:"required"`
	Note        string `json:"note"`
	Day         int    `json:"day" binding:"required"`
	Order       int    `json:"order" binding:"required"`
	ItineraryID int    `json:"-"`
	Email       string `json:"-"`
}
