package model

// UserInfo - структура запроса на добавление пользователя
type UserInfo struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            string
}

// UpdateUserInfo - структура для обновления данных пользователя
type UpdateUserInfo struct {
	ID   int64
	Name *string
	Role *string
}
