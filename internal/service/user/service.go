package user

import (
	"github.com/bogdanove/auth/internal/client/db"
	"github.com/bogdanove/auth/internal/repository"
	"github.com/bogdanove/auth/internal/service"
)

const createAction = "CREATE"
const updateAction = "UPDATE"
const deleteAction = "DELETE"
const getAction = "GET"

type userService struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

// NewUserService - конструктор сервиса чат
func NewUserService(
	userRepository repository.UserRepository,
	txManager db.TxManager,
) service.UserService {
	return &userService{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
