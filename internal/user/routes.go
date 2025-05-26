package user

import (
	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/middleware"
	"github.com/irmadev7/tripmate-backend/internal/repository"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)

	userHandler := UserHandler{repo: &userRepo}
	r.POST("/register", userHandler.RegisterHandler)
	r.POST("/login", userHandler.LoginHandler)

	protected := r.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/profile", userHandler.ProfileHandler)
	}

}
