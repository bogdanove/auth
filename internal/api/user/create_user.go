package user

import (
	"context"
	"errors"
	"log"

	"github.com/bogdanove/auth/internal/converter"
	"github.com/bogdanove/auth/pkg/user_v1"
)

// CreateUser - создание нового пользователя в системе
func (s *Server) CreateUser(ctx context.Context, req *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	id, err := s.userService.CreateUser(ctx, converter.ToUserInfoFromPB(req.UserInfo))
	if err != nil {
		return nil, err
	}

	log.Printf("user was created with id: %d", id)

	return &user_v1.CreateResponse{
		Id: id,
	}, nil
}
