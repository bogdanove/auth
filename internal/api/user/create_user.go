package user

import (
	"context"
	"log"

	"github.com/bogdanove/auth/internal/converter"
	"github.com/bogdanove/auth/pkg/user_v1"
)

// CreateUser - создание нового пользователя в системе
func (s *Server) CreateUser(ctx context.Context, req *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	info, err := converter.ToUserInfoFromPB(req.UserInfo)
	if err != nil {
		return nil, err
	}

	id, err := s.userService.CreateUser(ctx, info)
	if err != nil {
		return nil, err
	}

	log.Printf("user was created with id: %d", id)

	return &user_v1.CreateResponse{
		Id: id,
	}, nil
}
