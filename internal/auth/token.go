package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret []byte
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{
		secret: []byte(secret),
	}
}

func (s *TokenService) GenerateAccessToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"type":  "access",
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	})

	return token.SignedString(s.secret)
}

func (s *TokenService) GenerateRefreshToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"type":  "refresh",
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	return token.SignedString(s.secret)
}

func (s *TokenService) Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}

		return s.secret, nil
	})
}
