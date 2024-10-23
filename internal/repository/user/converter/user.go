package converter

import (
	repo "github.com/bogdanove/auth/internal/repository/user/model"
)

// FromServiceToLogRepo - конвертер для структуры сохранения действий пользователя
func FromServiceToLogRepo(id *int64, action string) *repo.UserLog {
	return &repo.UserLog{
		UserID: id,
		Action: action,
	}
}
