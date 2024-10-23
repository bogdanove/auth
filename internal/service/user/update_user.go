package user

import (
	"context"

	"github.com/bogdanove/auth/internal/model"
	"github.com/bogdanove/auth/internal/repository/user/converter"
)

// UpdateUser - обновление информации о пользователе по его идентификатору
func (s *userService) UpdateUser(ctx context.Context, req *model.UpdateUserInfo) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.userRepository.UpdateUser(ctx, req)
		if errTx != nil {
			return errTx
		}

		errTx = s.userRepository.SaveLog(ctx, converter.FromServiceToLogRepo(&req.ID, updateAction))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
