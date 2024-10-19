package user

import (
	"context"
	"log"

	"github.com/bogdanove/auth/pkg/user_v1"
)

// GetUser - получение информации о пользователе по его идентификатору
func (s *Server) GetUser(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	user, err := s.userService.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	log.Printf("user was received with id: %d", req.Id)

	return &user_v1.GetResponse{
		User: user,
	}, err
}
