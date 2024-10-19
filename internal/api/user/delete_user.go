package user

import (
	"context"
	"log"

	"github.com/bogdanove/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser - удаление пользователя из системы по его идентификатору
func (s *Server) DeleteUser(ctx context.Context, req *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	err := s.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	log.Printf("user was deleted with id: %d", req.Id)

	return &emptypb.Empty{}, nil
}
