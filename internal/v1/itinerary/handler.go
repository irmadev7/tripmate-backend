package itinerary

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/itinerary"
	"github.com/irmadev7/tripmate-backend/internal/model"
	"github.com/irmadev7/tripmate-backend/internal/pkg/apperror"
	"github.com/irmadev7/tripmate-backend/internal/pkg/response"
	"github.com/irmadev7/tripmate-backend/internal/pkg/utils"
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
		response.AppError(c, apperror.New(apperror.InvalidInput, "invalid request body", err), "invalid request")
		return
	}

	email, err := utils.GetEmail(c)
	if err != nil {
		response.AppError(c, apperror.New(apperror.Unauthorized, "unauthorized", err), "unauthorized email")
		return
	}

	itinerary := model.CreateItineraryRequest{
		Title:       input.Title,
		Description: input.Description,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Email:       email,
	}
	if err := h.service.CreateItinerary(c, itinerary); err != nil {
		response.AppError(c, err, "failed to create itinerary")
		return
	}

	response.Success(c, http.StatusCreated, "itinerary created", itinerary)
}

func (h *Handler) GetMyItineraries(c *gin.Context) {
	email, err := utils.GetEmail(c)
	if err != nil {
		response.AppError(c, apperror.New(apperror.Unauthorized, "unauthorized", err), "unauthorized email")
		return
	}

	itineraries, err := h.service.GetMyItineraries(c, email)
	if err != nil {
		response.AppError(c, err, "failed to fetch itineraries")
		return
	}

	response.Success(c, http.StatusOK, "processed successfully", itineraries)
}

func (h *Handler) AddPlaceToItinerary(c *gin.Context) {
	email, err := utils.GetEmail(c)
	if err != nil {
		response.AppError(c, apperror.New(apperror.Unauthorized, "unauthorized", err), "unauthorized email")
		return
	}

	var input model.AddPlaceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.AppError(c, apperror.New(apperror.InvalidInput, "invalid request body", err), "invalid request")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.AppError(c, apperror.New(apperror.InvalidInput, "invalid request param", err), "invalid request")
		return
	}

	input.ItineraryID = id
	input.Email = email
	if err := h.service.AddPlaceToItinerary(c, input); err != nil {
		response.AppError(c, err, "failed to add place")
		return
	}

	response.Success(c, http.StatusCreated, "place added", input)
}
