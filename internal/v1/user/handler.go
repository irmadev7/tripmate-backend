package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/model"
	"github.com/irmadev7/tripmate-backend/internal/pkg/response"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RegisterUser(c, req); err != nil {
		response.AppError(c, err, "failed to register user")
		return
	}

	c.JSON(http.StatusCreated, model.BaseResponse{Message: "User registered", Data: model.UserResponse{
		Email: req.Email,
		Name:  req.Name,
	}})
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusOK, model.BaseResponse{Message: "Login successful", Data: resp})
}

func (h *Handler) ProfileHandler(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email not found in context"})
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid email type"})
		return
	}

	profile, err := h.service.GetProfile(c, emailStr)
	if err != nil {
		response.AppError(c, err, "failed to get user")
		return
	}

	c.JSON(http.StatusOK, model.BaseResponse{Message: "Welcome to your profile!", Data: profile})
}

func (h *Handler) RefreshTokenHandler(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAccessToken, err := h.service.RefreshToken(c, req)
	if err != nil {
		response.AppError(c, err, "could not generate token")
		return
	}

	c.JSON(http.StatusOK, model.BaseResponse{Message: "Processed successfully", Data: newAccessToken})
}

func (h *Handler) LogoutHandler(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid email"})
		return
	}
	emailStr, ok := email.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid email"})
		return
	}
	err := h.service.Logout(c, emailStr)
	if err != nil {
		response.AppError(c, err, "failed to logout")
		return
	}

	c.JSON(http.StatusOK, model.BaseResponse{Message: "Logged out"})
}
