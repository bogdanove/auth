package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bogdanove/platform_common/pkg/db"
	txMocks "github.com/bogdanove/platform_common/pkg/db/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bogdanove/auth/internal/repository"
	repoMocks "github.com/bogdanove/auth/internal/repository/mocks"
	repo "github.com/bogdanove/auth/internal/repository/user/model"
	"github.com/bogdanove/auth/internal/service/user"
	"github.com/bogdanove/auth/pkg/user_v1"
)

func TestGet(t *testing.T) {
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

		id        = int64(1)
		name      = "some name"
		email     = "some email"
		role      = user_v1.Role_USER
		createdAt = timestamppb.New(time.Time{})

		repoErr = fmt.Errorf("repository error")

		res = &user_v1.User{
			Id:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: createdAt,
		}

		info = &repo.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      role.String(),
			CreatedAt: createdAt.AsTime(),
		}

		log = &repo.UserLog{
			UserID: &id,
			Action: "GET",
		}
	)

	tests := []struct {
		name          string
		args          args
		want          *user_v1.User
		err           error
		userRepoMock  userRepoMockFunc
		txManagerMock txManagerMockFunc
	}{
		{
			name: "get user repository success case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: res,
			err:  nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(info, nil)
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
			name: "get user repository error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: nil,
			err:  repoErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(nil, repoErr)
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

			response, err := service.GetUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}
