package itinerary

import (
	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/auth"
	"github.com/irmadev7/tripmate-backend/internal/itinerary"
	"github.com/irmadev7/tripmate-backend/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, itineraryService *itinerary.Service, tokenSvc *auth.TokenService) {
	itineraryHandler := NewHandler(itineraryService)

	itinerary := rg.Group("/itineraries")
	itinerary.Use(middleware.JWTAuthMiddleware(tokenSvc))
	{
		itinerary.POST("", itineraryHandler.CreateItinerary)
		itinerary.GET("", itineraryHandler.GetMyItineraries)
		itinerary.POST("/:id/places", itineraryHandler.AddPlaceToItinerary)
	}
}
