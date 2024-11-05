package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bogdanove/auth/internal/api/user"
	"github.com/bogdanove/auth/internal/service"
	serviceMocks "github.com/bogdanove/auth/internal/service/mocks"
	"github.com/bogdanove/auth/pkg/user_v1"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *user_v1.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = int64(1)
		name      = "some name"
		email     = "some email"
		role      = user_v1.Role_USER
		createdAt = timestamppb.New(time.Time{})
		updatedAt = timestamppb.New(time.Time{})

		serviceErr = fmt.Errorf("service error")

		userInfo = &user_v1.User{
			Id:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		req = &user_v1.GetRequest{
			Id: id,
		}

		res = &user_v1.GetResponse{
			User: userInfo,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *user_v1.GetResponse
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
				mock.GetUserMock.Expect(ctx, req.GetId()).Return(userInfo, nil)
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
				mock.GetUserMock.Expect(ctx, req.GetId()).Return(nil, serviceErr)
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

			response, err := api.GetUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}
