package itinerary

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/itinerary"
	"github.com/irmadev7/tripmate-backend/internal/model"
	"github.com/irmadev7/tripmate-backend/internal/pkg/response"
)

type Handler struct {
	service *itinerary.Service
}

func NewHandler(service *itinerary.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateItinerary(c *gin.Context) {
	var input model.CreateItineraryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	emailStr, ok := email.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	itinerary := model.CreateItineraryRequest{
		Title:       input.Title,
		Description: input.Description,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Email:       emailStr,
	}
	err := h.service.CreateItinerary(c, itinerary)
	if err != nil {
		response.AppError(c, err, "failed to create itinerary")
		return
	}

	c.JSON(http.StatusCreated, model.BaseResponse{Message: "Processed successfully", Data: itinerary})
}

func (h *Handler) GetMyItineraries(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	itineraries, err := h.service.GetMyItineraries(c, emailStr)
	if err != nil {
		response.AppError(c, err, "failed to fetch itineraries")
		return
	}

	c.JSON(http.StatusOK, model.BaseResponse{Message: "Processed successfully", Data: itineraries})
}

func (h *Handler) AddPlaceToItinerary(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	emailStr, ok := email.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input model.AddPlaceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid itinerary id"})
		return
	}

	input.ItineraryID = id
	input.Email = emailStr
	if err := h.service.AddPlaceToItinerary(c, input); err != nil {
		response.AppError(c, err, "failed to add place")
		return
	}

	c.JSON(http.StatusCreated, model.BaseResponse{Message: "Processed successfully", Data: input})
}
