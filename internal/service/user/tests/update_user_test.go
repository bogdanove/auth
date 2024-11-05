package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bogdanove/platform_common/pkg/db"
	txMocks "github.com/bogdanove/platform_common/pkg/db/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/bogdanove/auth/internal/model"
	"github.com/bogdanove/auth/internal/repository"
	repoMocks "github.com/bogdanove/auth/internal/repository/mocks"
	repo "github.com/bogdanove/auth/internal/repository/user/model"
	"github.com/bogdanove/auth/internal/service/user"
	"github.com/bogdanove/auth/pkg/user_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userRepoMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.UpdateUserInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id   = int64(1)
		name = "some name"
		role = user_v1.Role_USER.String()

		serviceErr = fmt.Errorf("repository error")

		reqAll = &model.UpdateUserInfo{
			ID:   id,
			Name: &name,
			Role: &role,
		}

		reqName = &model.UpdateUserInfo{
			ID:   id,
			Name: &name,
		}

		reqRole = &model.UpdateUserInfo{
			ID:   id,
			Role: &role,
		}

		log = &repo.UserLog{
			UserID: &id,
			Action: "UPDATE",
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
			name: "update user success case with all parameter",
			args: args{
				ctx: ctx,
				req: reqAll,
			},
			err: nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, reqAll).Return(nil)
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
			name: "update user success case with name parameter",
			args: args{
				ctx: ctx,
				req: reqName,
			},
			err: nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, reqName).Return(nil)
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
			name: "update user success case with role parameter",
			args: args{
				ctx: ctx,
				req: reqRole,
			},
			err: nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, reqRole).Return(nil)
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
			name: "update user repository error case",
			args: args{
				ctx: ctx,
				req: reqAll,
			},
			err: serviceErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, reqAll).Return(serviceErr)
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

			err := service.UpdateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
