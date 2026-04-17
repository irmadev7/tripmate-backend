package response

import (
	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/model"
)

func Success(c *gin.Context, status int, message string, data interface{}, meta interface{}) {
	c.JSON(status, model.BaseResponse{
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
