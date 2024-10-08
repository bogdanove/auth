package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bogdanove/auth/pkg/user_v1"
)

const grpcPort = 50051

type server struct {
	user_v1.UnimplementedUserV1Server
}

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	reflection.Register(s)

	user_v1.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at: %v", listen.Addr())

	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Create - создание нового пользователя в системе
func (s *server) Create(_ context.Context, req *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	log.Printf("start creating new user with name: %s", req.GetUserInfo().GetName())

	id := gofakeit.Int64()

	log.Printf("new user was created with id: %d", id)

	return &user_v1.CreateResponse{
		Id: id,
	}, nil
}

// Get - получение информации о пользователе по его идентификатору
func (s *server) Get(_ context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	log.Printf("receiving user with id: %d", req.GetId())

	return &user_v1.GetResponse{
		User: &user_v1.User{
			Id:        req.GetId(),
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			Role:      user_v1.Role_USER,
			CreatedAt: timestamppb.New(time.Now()),
		},
	}, nil
}

// Update - обновление информации о пользователе по его идентификатору
func (s *server) Update(_ context.Context, req *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("updating user with id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

// Delete - удаление пользователя из системы по его идентификатору
func (s *server) Delete(_ context.Context, req *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("deleting user with id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}
