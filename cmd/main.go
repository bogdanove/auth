package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"github.com/bogdanove/auth/internal/config"
	"github.com/bogdanove/auth/internal/config/env"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bogdanove/auth/pkg/user_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	user_v1.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	listen, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	user_v1.RegisterUserV1Server(s, &server{pool: pool})

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
