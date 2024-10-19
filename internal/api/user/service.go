package user

import (
	"github.com/bogdanove/auth/internal/service"
	"github.com/bogdanove/auth/pkg/user_v1"
)

// Server - сервер GRPC
type Server struct {
	user_v1.UnimplementedUserV1Server
	userService service.UserService
}

// NewServerImplementation - имплементация сервера GRPC
func NewServerImplementation(userService service.UserService) *Server {
	return &Server{
		userService: userService,
	}
}
