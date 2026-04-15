package repository

import (
	"context"

	"github.com/irmadev7/tripmate-backend/internal/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) userRepository {
	return userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateRefreshToken(ctx context.Context, email, token string) error {
	query := `UPDATE users SET refresh_token=$1 WHERE email=$2`
	return r.db.WithContext(ctx).Exec(query, token, email).Error
}
