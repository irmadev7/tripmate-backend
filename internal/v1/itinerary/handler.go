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

// CreateItinerary godoc
// @Summary Create itinerary by user
// @Description Create itinerary by user
// @Tags itinerary
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateItineraryRequest true "Create itinerary"
// @Success 201 {object} model.BaseResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/v1/itineraries [post]
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

	input.Email = email
	if err := h.service.CreateItinerary(c, input); err != nil {
		response.AppError(c, err, "failed to create itinerary")
		return
	}

	response.Success(c, http.StatusCreated, "itinerary created", input, nil)
}

// GetMyItineraries godoc
// @Summary Get user's itineraries
// @Tags itinerary
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Page Limit"
// @Param search query string false "Search by title"
// @Success 200 {object} model.BaseResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/v1/itineraries [get]
func (h *Handler) GetMyItineraries(c *gin.Context) {
	email, err := utils.GetEmail(c)
	if err != nil {
		response.AppError(c, apperror.New(apperror.Unauthorized, "unauthorized", err), "unauthorized email")
		return
	}

	var req model.GetMyItineraryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.AppError(c,
			apperror.New(apperror.InvalidInput, "invalid query params", err), "invalid request",
		)
		return
	}

	// default value
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	req.Email = email
	itineraries, err := h.service.GetMyItineraries(c, req)
	if err != nil {
		response.AppError(c, err, "failed to fetch itineraries")
		return
	}

	response.Success(c, http.StatusOK, "processed successfully", itineraries.Data, itineraries.Meta)
}

// AddPlaceToItinerary godoc
// @Summary Add place to itinerary
// @Description Add place itinerary
// @Tags itinerary
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Itinerary ID"
// @Param request body model.AddPlaceRequest true "Add place to itinerary"
// @Success 201 {object} model.BaseResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/v1/itineraries/{id}/places [post]
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

	response.Success(c, http.StatusCreated, "place added", input, nil)
}
