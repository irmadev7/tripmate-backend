package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/irmadev7/tripmate-backend/internal/pkg/apperror"
)

func AppError(c *gin.Context, err error, fallbackMsg string) {
	var appErr *apperror.Error
	if !errors.As(err, &appErr) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fallbackMsg})
		return
	}

	c.JSON(statusCode(appErr.Code), gin.H{"error": appErr.Message})
}

func statusCode(code apperror.Code) int {
	switch code {
	case apperror.InvalidInput:
		return http.StatusBadRequest
	case apperror.Unauthorized:
		return http.StatusUnauthorized
	case apperror.Conflict:
		return http.StatusConflict
	case apperror.NotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
