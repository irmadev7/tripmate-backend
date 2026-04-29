package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/model"
	"github.com/irmadev7/tripmate-backend/internal/pkg/apperror"
	"github.com/irmadev7/tripmate-backend/internal/pkg/response"
	"github.com/irmadev7/tripmate-backend/internal/user"
)

type Handler struct {
	service *user.ServiceV2
}

func NewHandler(service *user.ServiceV2) *Handler {
	return &Handler{service: service}
}
func (h *Handler) LoginV2Handler(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, apperror.New(apperror.InvalidInput, "invalid request", err), "invalid request")
		return
	}

	res, err := h.service.LoginV2(c, req)
	if err != nil {
		response.AppError(c, err, "login failed")
		return
	}

	response.Success(c, http.StatusOK, "login success", res, nil)
}

func (h *Handler) RefreshTokenV2Handler(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, apperror.New(apperror.InvalidInput, "invalid request body", err), "invalid request")
		return
	}

	newAccessToken, err := h.service.RefreshTokenV2(c, req)
	if err != nil {
		response.AppError(c, err, "could not generate token")
		return
	}

	response.Success(c, http.StatusOK, "token refreshed", newAccessToken, nil)
}

func (h *Handler) LogoutV2Handler(c *gin.Context) {
	var req model.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, apperror.New(apperror.InvalidInput, "invalid request body", err), "invalid request")
		return
	}

	if err := h.service.LogoutV2(c, req.RefreshToken); err != nil {
		response.AppError(c, err, "failed to logout")
		return
	}

	response.Success(c, http.StatusOK, "logged out", nil, nil)
}
