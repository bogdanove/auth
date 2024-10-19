package converter

import (
	service "github.com/bogdanove/auth/internal/model"
	repo "github.com/bogdanove/auth/internal/repository/user/model"
)

// ToUserFromRepo - конвертер структуры пользователя для сервисного слоя из репо
func ToUserFromRepo(user *repo.User) *service.User {
	return &service.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// FromServiceToUserInfoRepo - конвертер информации о пользователе из сервисного в репо слой
func FromServiceToUserInfoRepo(user *service.UserInfo) *repo.UserInfo {
	return &repo.UserInfo{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            user.Role,
	}
}

// FromServiceToUpdateUserInfoRepo - конвертер для структуры обновления данных пользователя из сервисного в репо слой
func FromServiceToUpdateUserInfoRepo(info *service.UpdateUserInfo) *repo.UpdateUserInfo {
	return &repo.UpdateUserInfo{
		ID:   info.ID,
		Name: info.Name,
		Role: info.Role,
	}
}

// FromServiceToLogRepo - конвертер для структуры сохранения действий пользователя
func FromServiceToLogRepo(id int64, action string) *repo.UserLog {
	return &repo.UserLog{
		UserID: id,
		Action: action,
	}
}
