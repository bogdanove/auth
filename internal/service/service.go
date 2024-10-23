package service

import (
	"context"

	"github.com/bogdanove/auth/internal/model"
	"github.com/bogdanove/auth/pkg/user_v1"
)

// UserService - сервис пользователя
type UserService interface {
	CreateUser(ctx context.Context, req *model.UserInfo) (int64, error)
	GetUser(ctx context.Context, id int64) (*user_v1.User, error)
	UpdateUser(ctx context.Context, req *model.UpdateUserInfo) error
	DeleteUser(ctx context.Context, id int64) error
}
