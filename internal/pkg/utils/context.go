package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetEmail(c *gin.Context) (string, error) {
	email, ok := c.Get("email")
	if !ok {
		return "", errors.New("email not found")
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", errors.New("invalid email type")
	}

	return emailStr, nil
}
