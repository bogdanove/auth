package model

import (
	"database/sql"
	"time"
)

// UserInfo - структура запроса на добавление пользователя
type UserInfo struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            string
}

// User - структура пользователя
type User struct {
	ID        int64
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// UpdateUserInfo - структура для обновления данных пользователя
type UpdateUserInfo struct {
	ID   int64
	Name *string
	Role *string
}
