package model

import (
	"database/sql"
	"time"
)

// UserInfo - структура запроса на добавление пользователя репо
type UserInfo struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            string
}

// User - структура пользователя для бд
type User struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Role      string       `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// UpdateUserInfo - структура для обновления данных пользователя в бд
type UpdateUserInfo struct {
	ID   int64
	Name string
	Role string
}

// UserLog - структура для хранения действий пользователя
type UserLog struct {
	UserID int64
	Action string
}
