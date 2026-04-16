package user

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/irmadev7/tripmate-backend/internal/auth"
	"github.com/irmadev7/tripmate-backend/internal/model"
	"github.com/irmadev7/tripmate-backend/internal/pkg/apperror"
	"github.com/irmadev7/tripmate-backend/internal/pkg/utils"
	"github.com/irmadev7/tripmate-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo     repository.UserRepository
	tokenSvc *auth.TokenService
}

func NewService(repo repository.UserRepository, tokenSvc *auth.TokenService) *Service {
	return &Service{repo: repo, tokenSvc: tokenSvc}
}

func (s *Service) RegisterUser(ctx context.Context, req model.UserRequest) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.New(apperror.Internal, "failed to hash password", err)
	}

	user := model.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashed),
	}
	if err := s.repo.CreateUser(ctx, &user); err != nil {
		if repository.IsDuplicateKeyError(err) {
			return apperror.New(apperror.Conflict, "email already exists", err)
		}
		return apperror.New(apperror.Internal, "failed to create user", err)
	}

	go utils.SendEmail(user.Email, user.Name)

	return nil
}

func (s *Service) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.Unauthorized, "invalid credentials", err)
		}
		return nil, apperror.New(apperror.Internal, "failed to get user", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperror.New(apperror.Unauthorized, "invalid credentials", err)
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(user.Email)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "could not generate access token", err)

	}

	refreshToken, err := s.tokenSvc.GenerateRefreshToken(user.Email)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "could not generate refresh token", err)

	}

	if err = s.repo.UpdateRefreshToken(ctx, user.Email, refreshToken); err != nil {
		return nil, apperror.New(apperror.Internal, "failed to save refresh token", err)

	}
	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) GetProfile(ctx context.Context, email string) (*model.UserResponse, error) {
	profile, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.NotFound, "user not found", err)
		}
		return nil, apperror.New(apperror.Internal, "failed to get user", err)
	}

	return &model.UserResponse{Email: profile.Email, Name: profile.Name}, nil
}

func (s *Service) RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (*model.RefreshTokenResponse, error) {
	token, err := s.tokenSvc.Parse(req.RefreshToken)
	if err != nil || !token.Valid {
		return nil, apperror.New(apperror.Unauthorized, "invalid refresh token", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apperror.New(apperror.Unauthorized, "invalid claims", nil)
	}

	// validate token type
	tokenType, ok := claims["type"]
	if !ok {
		return nil, apperror.New(apperror.Unauthorized, "invalid token type", nil)
	}
	tokenTypeStr, ok := tokenType.(string)
	if !ok {
		return nil, apperror.New(apperror.Unauthorized, "invalid token type", nil)
	}
	if tokenTypeStr != "refresh" {
		return nil, apperror.New(apperror.Unauthorized, "invalid token type", nil)

	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, apperror.New(apperror.Unauthorized, "invalid token email", nil)

	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil || user.RefreshToken != req.RefreshToken {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.Unauthorized, "invalid refresh token", err)
		}
		if err != nil {
			return nil, apperror.New(apperror.Internal, "failed to get user", err)
		}
		return nil, apperror.New(apperror.Unauthorized, "invalid refresh token", nil)
	}

	newAccessToken, err := s.tokenSvc.GenerateAccessToken(email)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "could not generate access token", err)
	}

	return &model.RefreshTokenResponse{AccessToken: newAccessToken}, nil
}

func (s *Service) Logout(ctx context.Context, email string) error {
	err := s.repo.UpdateRefreshToken(ctx, email, "")
	if err != nil {
		return apperror.New(apperror.Internal, "failed to logout", err)
	}

	return nil
}
