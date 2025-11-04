package itinerary

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/model"
	"github.com/irmadev7/tripmate-backend/internal/repository"
)

type ItineraryHandler struct {
	itineraryRepo repository.ItineraryRepository
	userRepo      repository.UserRepository
	placeRepo     repository.PlaceRepository
}

func NewItineraryHandler(itineraryRepo repository.ItineraryRepository, userRepo repository.UserRepository, placeRepo repository.PlaceRepository) *ItineraryHandler {
	return &ItineraryHandler{
		itineraryRepo: itineraryRepo,
		userRepo:      userRepo,
		placeRepo:     placeRepo,
	}
}

func (h *ItineraryHandler) CreateItinerary(c *gin.Context) {
	var input model.CreateItineraryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	loginUser, err := h.userRepo.GetUserByUsername(c, username.(string))
	if err != nil || loginUser == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user lookup failed"})
		return
	}

	itinerary := model.Itinerary{
		Title:       input.Title,
		Description: input.Description,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		UserID:      loginUser.ID,
	}
	if err := h.itineraryRepo.CreateItinerary(c, &itinerary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch itineraries"})
		return
	}

	c.JSON(http.StatusCreated, itinerary)
}

func (h *ItineraryHandler) GetMyItineraries(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	loginUser, err := h.userRepo.GetUserByUsername(c, username.(string))
	if err != nil || loginUser == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user lookup failed"})
		return
	}

	itineraries, err := h.itineraryRepo.GetItineraryByUser(c, int(loginUser.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch itineraries"})
		return
	}

	c.JSON(http.StatusOK, itineraries)
}

func (h *ItineraryHandler) AddPlaceToItinerary(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	itineraryID := c.Param("id")
	itineraryIDInt, _ := strconv.Atoi(itineraryID)
	itinerary, err := h.itineraryRepo.GetItineraryById(c, itineraryIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "itinerary not found"})
		return
	}

	loginUser, err := h.userRepo.GetUserByUsername(c, username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user lookup failed"})
		return
	}
	if itinerary.UserID != loginUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't own this itinerary"})
		return
	}

	var input model.AddPlaceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	place := model.Destination{
		ItineraryID: itinerary.ID,
		Name:        input.Name,
		Note:        input.Note,
		Day:         input.Day,
		Order:       input.Order,
	}

	if err := h.placeRepo.AddPlaceToItinerary(c, &place).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to add place"})
		return
	}

	c.JSON(201, place)
}
