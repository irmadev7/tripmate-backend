package user

import (
	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/auth"
	"github.com/irmadev7/tripmate-backend/internal/middleware"
	"github.com/irmadev7/tripmate-backend/internal/user"
)

func RegisterRoutes(rg *gin.RouterGroup, userService *user.Service, tokenSvc *auth.TokenService) {
	userHandler := NewHandler(userService)

	auth := rg.Group("/auth")
	{
		auth.POST("/register", userHandler.RegisterHandler)
		auth.POST("/login", userHandler.LoginHandler)
		auth.POST("/refresh", userHandler.RefreshTokenHandler)
	}

	user := rg.Group("/users")
	user.Use(middleware.JWTAuthMiddleware(tokenSvc))
	{
		user.GET("/profile", userHandler.ProfileHandler)
		user.POST("/logout", userHandler.LogoutHandler)
	}
}
