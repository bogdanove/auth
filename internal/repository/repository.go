package repository

import (
	"context"

	service "github.com/bogdanove/auth/internal/model"
	"github.com/bogdanove/auth/internal/repository/user/model"
)

// UserRepository - репозиторий пользователя
type UserRepository interface {
	CreateUser(ctx context.Context, req *model.UserInfo) (int64, error)
	GetUser(ctx context.Context, id int64) (*service.User, error)
	UpdateUser(ctx context.Context, req *model.UpdateUserInfo) error
	DeleteUser(ctx context.Context, id int64) error
	SaveLog(ctx context.Context, req *model.UserLog) error
}
