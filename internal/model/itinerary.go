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
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type AddPlaceRequest struct {
	Name  string `json:"name"`
	Note  string `json:"note"`
	Day   int    `json:"day"`
	Order int    `json:"order"`
}
