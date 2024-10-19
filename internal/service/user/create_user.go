package user

import (
	"context"

	"github.com/bogdanove/auth/internal/model"
	"github.com/bogdanove/auth/internal/repository/user/converter"
)

const createAction = "CREATE"

// CreateUser - создание нового пользователя в системе
func (s *userService) CreateUser(ctx context.Context, req *model.UserInfo) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.CreateUser(ctx, converter.FromServiceToUserInfoRepo(req))
		if errTx != nil {
			return errTx
		}

		errTx = s.userRepository.SaveLog(ctx, converter.FromServiceToLogRepo(id, createAction))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
