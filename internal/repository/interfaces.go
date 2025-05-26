package repository

import (
	"context"

	"github.com/irmadev7/tripmate-backend/internal/model"
)

type Sql interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
}
