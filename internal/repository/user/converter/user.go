package converter

import (
	"errors"

	service "github.com/bogdanove/auth/internal/model"
	repo "github.com/bogdanove/auth/internal/repository/user/model"
)

// ToUserFromRepo - конвертер структуры пользователя для сервисного слоя из репо
func ToUserFromRepo(user *repo.User) (*service.User, error) {
	if user == nil {
		return nil, errors.New("user for convert is nil")
	}
	return &service.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// FromServiceToLogRepo - конвертер для структуры сохранения действий пользователя
func FromServiceToLogRepo(id *int64, action string) *repo.UserLog {
	return &repo.UserLog{
		UserID: id,
		Action: action,
	}
}
