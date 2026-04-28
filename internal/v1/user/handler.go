package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/model"
	"github.com/irmadev7/tripmate-backend/internal/pkg/apperror"
	"github.com/irmadev7/tripmate-backend/internal/pkg/response"
	"github.com/irmadev7/tripmate-backend/internal/pkg/utils"
	"github.com/irmadev7/tripmate-backend/internal/user"
)

type Handler struct {
	service *user.Service
}

func NewHandler(service *user.Service) *Handler {
	return &Handler{service: service}
}

// RegisterHandler godoc
// @Summary Register new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.UserRequest true "Register request"
// @Success 201 {object} model.BaseResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *Handler) RegisterHandler(c *gin.Context) {
	var req model.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, apperror.New(apperror.InvalidInput, "invalid request body", err), "invalid request")
		return
	}

	resp, err := h.service.RegisterUser(c, req)
	if err != nil {
		response.AppError(c, err, "failed to register user")
		return
	}

	response.Success(c, http.StatusCreated, "user registered", resp, nil)
}

// LoginHandler godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login request"
// @Success 200 {object} model.BaseResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *Handler) LoginHandler(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, apperror.New(apperror.InvalidInput, "invalid request body", err), "invalid request")
		return
	}

	resp, err := h.service.Login(c, req)
	if err != nil {
		response.AppError(c, err, "failed to login")
		return
	}

	response.Success(c, http.StatusOK, "login successful", resp, nil)
}

// ProfileHandler godoc
// @Summary Get user profile
// @Tags user
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.BaseResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/v1/users/profile [get]
func (h *Handler) ProfileHandler(c *gin.Context) {
	email, err := utils.GetEmail(c)
	if err != nil {
		response.AppError(c, apperror.New(apperror.Unauthorized, "unauthorized", err), "unauthorized email")
		return
	}

	profile, err := h.service.GetProfile(c, email)
	if err != nil {
		response.AppError(c, err, "failed to get user")
		return
	}

	response.Success(c, http.StatusOK, "profile fetched", profile, nil)
}

// RefreshTokenHandler godoc
// @Summary refresh token user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} model.BaseResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshTokenHandler(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, apperror.New(apperror.InvalidInput, "invalid request body", err), "invalid request")
		return
	}

	newAccessToken, err := h.service.RefreshToken(c, req)
	if err != nil {
		response.AppError(c, err, "could not generate token")
		return
	}

	response.Success(c, http.StatusOK, "token refreshed", newAccessToken, nil)
}

// LogoutHandler godoc
// @Summary logout user
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.BaseResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/v1/users/logout [post]
func (h *Handler) LogoutHandler(c *gin.Context) {
	email, err := utils.GetEmail(c)
	if err != nil {
		response.AppError(c, apperror.New(apperror.Unauthorized, "unauthorized", err), "unauthorized email")
		return
	}
	if err := h.service.Logout(c, email); err != nil {
		response.AppError(c, err, "failed to logout")
		return
	}

	response.Success(c, http.StatusOK, "logged out", nil, nil)
}
