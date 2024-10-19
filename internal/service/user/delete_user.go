package user

import (
	"context"

	"github.com/bogdanove/auth/internal/repository/user/converter"
)

const deleteAction = "DELETE"

// DeleteUser - удаление пользователя из системы по его идентификатору
func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.userRepository.DeleteUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.userRepository.SaveLog(ctx, converter.FromServiceToLogRepo(id, deleteAction))
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
