package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bogdanove/auth/pkg/user_v1"
	"github.com/bogdanove/platform_common/pkg/db"
	txMocks "github.com/bogdanove/platform_common/pkg/db/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/bogdanove/auth/internal/model"
	"github.com/bogdanove/auth/internal/repository"
	repoMocks "github.com/bogdanove/auth/internal/repository/mocks"
	repo "github.com/bogdanove/auth/internal/repository/user/model"
	"github.com/bogdanove/auth/internal/service/user"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userRepoMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.UserInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = int64(1)
		name     = "some name"
		email    = "some email"
		password = "some password"
		role     = user_v1.Role_USER

		repoErr = fmt.Errorf("repository error")

		req = &model.UserInfo{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            role.String(),
		}

		log = &repo.UserLog{
			UserID: &id,
			Action: "CREATE",
		}
	)

	tests := []struct {
		name          string
		args          args
		want          int64
		err           error
		userRepoMock  userRepoMockFunc
		txManagerMock txManagerMockFunc
	}{
		{
			name: "create user repository success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, req).Return(id, nil)
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
			name: "create user repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repoErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, req).Return(0, repoErr)
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

			response, err := service.CreateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}
