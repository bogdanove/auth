package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bogdanove/auth/internal/api/user"
	"github.com/bogdanove/auth/internal/model"
	"github.com/bogdanove/auth/internal/service"
	serviceMocks "github.com/bogdanove/auth/internal/service/mocks"
	"github.com/bogdanove/auth/pkg/user_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *user_v1.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id   = int64(1)
		name = wrapperspb.StringValue{Value: "some name"}
		role = user_v1.Role_USER

		nameUpd = name.Value
		roleUpd = role.String()

		serviceErr = fmt.Errorf("service error")

		userInfoAll = &user_v1.UpdateUserInfo{
			Name: &name,
			Role: &role,
		}

		userInfoName = &user_v1.UpdateUserInfo{
			Name: &name,
		}

		userInfoRole = &user_v1.UpdateUserInfo{
			Role: &role,
		}

		reqAll = &user_v1.UpdateRequest{
			Id:             id,
			UpdateUserInfo: userInfoAll,
		}

		reqName = &user_v1.UpdateRequest{
			Id:             id,
			UpdateUserInfo: userInfoName,
		}

		reqRole = &user_v1.UpdateRequest{
			Id:             id,
			UpdateUserInfo: userInfoRole,
		}

		res = &emptypb.Empty{}

		infoAll = &model.UpdateUserInfo{
			ID:   id,
			Name: &nameUpd,
			Role: &roleUpd,
		}

		infoName = &model.UpdateUserInfo{
			ID:   id,
			Name: &nameUpd,
		}

		infoRole = &model.UpdateUserInfo{
			ID:   id,
			Role: &roleUpd,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case with all parameter",
			args: args{
				ctx: ctx,
				req: reqAll,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, infoAll).Return(nil)
				return mock
			},
		},
		{
			name: "success case with name parameter",
			args: args{
				ctx: ctx,
				req: reqName,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, infoName).Return(nil)
				return mock
			},
		},
		{
			name: "success case with role parameter",
			args: args{
				ctx: ctx,
				req: reqRole,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, infoRole).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: reqAll,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, infoAll).Return(serviceErr)
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

			response, err := api.UpdateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}
