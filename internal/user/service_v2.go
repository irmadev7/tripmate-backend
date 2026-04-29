package user

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/irmadev7/tripmate-backend/internal/auth"
	"github.com/irmadev7/tripmate-backend/internal/model"
	"github.com/irmadev7/tripmate-backend/internal/pkg/apperror"
	"github.com/irmadev7/tripmate-backend/internal/repository"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type ServiceV2 struct {
	repo     repository.UserRepository
	tokenSvc *auth.TokenService
	cache    *redis.Client
}

func NewServiceV2(repo repository.UserRepository, tokenSvc *auth.TokenService, redis *redis.Client) *ServiceV2 {
	return &ServiceV2{repo: repo, tokenSvc: tokenSvc, cache: redis}
}

func (s *ServiceV2) LoginV2(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperror.New(apperror.Unauthorized, "invalid credentials", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperror.New(apperror.Unauthorized, "invalid credentials", err)
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(user.Email)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to generate access token", err)
	}

	refreshToken, jti, err := s.tokenSvc.GenerateRefreshTokenV2(user.Email)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to generate refresh token", err)
	}

	key := fmt.Sprintf("refresh:%s", jti)
	err = s.cache.Set(ctx, key, user.Email, 7*24*time.Hour).Err()
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to store refresh token", err)
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *ServiceV2) RefreshTokenV2(ctx context.Context, req model.RefreshTokenRequest) (*model.RefreshTokenResponse, error) {
	token, err := s.tokenSvc.Parse(req.RefreshToken)
	if err != nil {
		return nil, apperror.New(apperror.Unauthorized, "invalid token", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apperror.New(apperror.Unauthorized, "invalid token claims", nil)
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, apperror.New(apperror.Unauthorized, "invalid jti", nil)
	}

	email := claims["email"].(string)

	key := fmt.Sprintf("refresh:%s", jti)
	val, err := s.cache.Get(ctx, key).Result()
	if err != nil || val != email {
		return nil, apperror.New(apperror.Unauthorized, "invalid refresh token", err)
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(email)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to generate token", err)
	}

	return &model.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *ServiceV2) LogoutV2(ctx context.Context, refreshToken string) error {
	token, err := s.tokenSvc.Parse(refreshToken)
	if err != nil {
		return apperror.New(apperror.Unauthorized, "invalid token", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return apperror.New(apperror.Unauthorized, "invalid token claims", nil)
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return apperror.New(apperror.Unauthorized, "invalid jti", nil)
	}

	key := fmt.Sprintf("refresh:%s", jti)
	if err := s.cache.Del(ctx, key).Err(); err != nil {
		return apperror.New(apperror.Internal, "failed to logout", err)
	}

	return nil
}
