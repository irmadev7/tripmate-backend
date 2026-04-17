package user

import (
	"log"
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

func (h *Handler) RegisterHandler(c *gin.Context) {
	var req model.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		response.AppError(c, apperror.New(apperror.InvalidInput, "invalid request body", err), "invalid request")
		return
	}

	resp, err := h.service.RegisterUser(c, req)
	if err != nil {
		response.AppError(c, err, "failed to register user")
		return
	}

	response.Success(c, http.StatusCreated, "user registered", resp)
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, apperror.New(apperror.InvalidInput, "invalid request body", err), "invalid request")
		return
	}

	resp, err := h.service.Login(c, model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.AppError(c, err, "failed to login")
		return
	}

	response.Success(c, http.StatusOK, "login successful", resp)
}

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

	response.Success(c, http.StatusOK, "welcome to your profile!", profile)
}

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

	response.Success(c, http.StatusOK, "processed successfully", newAccessToken)
}

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

	response.Success(c, http.StatusOK, "logged out", nil)
}
