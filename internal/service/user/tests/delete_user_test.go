package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bogdanove/platform_common/pkg/db"
	txMocks "github.com/bogdanove/platform_common/pkg/db/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/bogdanove/auth/internal/repository"
	repoMocks "github.com/bogdanove/auth/internal/repository/mocks"
	repo "github.com/bogdanove/auth/internal/repository/user/model"
	"github.com/bogdanove/auth/internal/service/user"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type userRepoMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = int64(1)

		repoErr = fmt.Errorf("repository error")

		log = &repo.UserLog{
			UserID: &id,
			Action: "DELETE",
		}
	)

	tests := []struct {
		name          string
		args          args
		err           error
		userRepoMock  userRepoMockFunc
		txManagerMock txManagerMockFunc
	}{
		{
			name: "delete user repository success case",
			args: args{
				ctx: ctx,
				req: id,
			},
			err: nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(nil)
				mock.SaveLogMock.Expect(ctx, log).Return(nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "delete user repository error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			err: repoErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(repoErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepoMock(mc)
			txManagerMock := tt.txManagerMock(mc)
			service := user.NewUserService(userRepoMock, txManagerMock)

			err := service.DeleteUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
