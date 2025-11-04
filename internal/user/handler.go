package user

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/irmadev7/tripmate-backend/internal/model"
	"github.com/irmadev7/tripmate-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(r repository.UserRepository) *UserHandler {
	return &UserHandler{repo: r}
}

func (h *UserHandler) RegisterHandler(c *gin.Context) {
	var req model.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := h.repo.GetUserByUsername(c, req.Username)
	if err == nil || existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hashed),
	}
	if err := h.repo.CreateUser(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

func (h *UserHandler) LoginHandler(c *gin.Context) {
	var req model.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetUserByUsername(c, req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "defaultsecret"
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (h *UserHandler) ProfileHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(500, gin.H{"error": "username not found in context"})
		return
	}

	c.JSON(200, gin.H{
		"message":  "Welcome to your profile!",
		"username": username,
	})
}
