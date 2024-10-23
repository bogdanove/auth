package user

import (
	"context"

	srv "github.com/bogdanove/auth/internal/converter"
	"github.com/bogdanove/auth/internal/repository/user/converter"
	"github.com/bogdanove/auth/internal/repository/user/model"
	"github.com/bogdanove/auth/pkg/user_v1"
)

// GetUser - получение информации о пользователе по его идентификатору
func (s *userService) GetUser(ctx context.Context, id int64) (*user_v1.User, error) {
	var user *model.User
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		user, errTx = s.userRepository.GetUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.userRepository.SaveLog(ctx, converter.FromServiceToLogRepo(&id, getAction))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	result, err := srv.ToPBUserFromRepo(user)
	if err != nil {
		return nil, err
	}

	return result, nil
}
