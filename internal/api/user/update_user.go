package user

import (
	"context"
	"log"

	"github.com/bogdanove/auth/internal/converter"
	"github.com/bogdanove/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser - обновление информации о пользователе по его идентификатору
func (s *Server) UpdateUser(ctx context.Context, req *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	info, err := converter.ToUpdateUserInfoFromPB(req)
	if err != nil {
		return nil, err
	}

	err = s.userService.UpdateUser(ctx, info)
	if err != nil {
		return nil, err
	}

	log.Printf("user was updated with id: %d", req.Id)

	return &emptypb.Empty{}, nil
}
