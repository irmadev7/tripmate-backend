package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *TokenService) GenerateRefreshTokenV2(email string) (string, string, error) {
	jti := uuid.NewString()

	claims := jwt.MapClaims{
		"email": email,
		"jti":   jti,
		"type":  "refresh",
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", "", err
	}

	return tokenStr, jti, nil
}
