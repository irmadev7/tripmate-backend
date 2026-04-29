package user

import (
	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/auth"
	"github.com/irmadev7/tripmate-backend/internal/user"
)

func RegisterRoutes(rg *gin.RouterGroup, userServiceV2 *user.ServiceV2, tokenSvc *auth.TokenService) {
	userHandler := NewHandler(userServiceV2)

	auth := rg.Group("/auth")
	{
		auth.POST("/login", userHandler.LoginV2Handler)
		auth.POST("/refresh", userHandler.RefreshTokenV2Handler)
		auth.POST("/logout", userHandler.LogoutV2Handler)
	}
}
