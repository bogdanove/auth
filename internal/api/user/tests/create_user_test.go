package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/bogdanove/auth/internal/api/user"
	"github.com/bogdanove/auth/internal/model"
	"github.com/bogdanove/auth/internal/service"
	serviceMocks "github.com/bogdanove/auth/internal/service/mocks"
	"github.com/bogdanove/auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *user_v1.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = int64(1)
		name     = "some name"
		email    = "some email"
		password = "some password"
		role     = user_v1.Role_USER

		serviceErr = fmt.Errorf("service error")

		userInfo = &user_v1.UserInfo{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            role,
		}

		req = &user_v1.CreateRequest{
			UserInfo: userInfo,
		}

		info = &model.UserInfo{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            role.String(),
		}

		res = &user_v1.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *user_v1.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, info).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, info).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			noteServiceMock := tt.userServiceMock(mc)
			api := user.NewServerImplementation(noteServiceMock)

			response, err := api.CreateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}
