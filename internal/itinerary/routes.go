package itinerary

import (
	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/middleware"
	"github.com/irmadev7/tripmate-backend/internal/repository"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	itineraryRepo := repository.NewItineraryRepository(db)
	userRepo := repository.NewUserRepository(db)
	placeRepo := repository.NewPlaceRepository(db)

	itineraryHandler := ItineraryHandler{itineraryRepo: &itineraryRepo, userRepo: &userRepo, placeRepo: &placeRepo}
	protected := r.Group("/itineraries")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.POST("", itineraryHandler.CreateItinerary)
		protected.GET("", itineraryHandler.GetMyItineraries)
		protected.POST("/:id/places", itineraryHandler.AddPlaceToItinerary)
	}
}
