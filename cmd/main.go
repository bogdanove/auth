package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bogdanove/auth/internal/config"
	"github.com/bogdanove/auth/internal/config/env"
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

// CreateUser - создание нового пользователя в системе
func (s *server) CreateUser(ctx context.Context, req *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	log.Printf("start creating new user with name: %s", req.GetUserInfo().GetName())

	queryAddUsers, args, err := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role").
		Values(req.GetUserInfo().GetName(), req.GetUserInfo().GetEmail(),
			req.GetUserInfo().GetRole()).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	var userID int64
	err = s.pool.QueryRow(ctx, queryAddUsers, args...).Scan(&userID)
	if err != nil {
		log.Printf("failed to insert user: %v", err)
		return nil, err
	}

	log.Printf("new user was created with id: %d", userID)

	return &user_v1.CreateResponse{
		Id: userID,
	}, nil
}

// GetUser - получение информации о пользователе по его идентификатору
func (s *server) GetUser(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	log.Printf("receiving user with id: %d", req.GetId())

	query, args, err := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1).ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	var id int64
	var name, email, role string
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("failed to select user: %v", err)
		return nil, err
	}

	return &user_v1.GetResponse{
		User: &user_v1.User{
			Id:        id,
			Name:      name,
			Email:     email,
			Role:      user_v1.Role(user_v1.Role_value[role]),
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt.Time),
		},
	}, nil
}

// UpdateUser - обновление информации о пользователе по его идентификатору
func (s *server) UpdateUser(ctx context.Context, req *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("updating user with id: %d", req.GetId())

	query, args, err := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("name", req.UpdateUserInfo.GetName().Value).
		Set("role", req.UpdateUserInfo.GetRole()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.GetId()}).ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// DeleteUser - удаление пользователя из системы по его идентификатору
func (s *server) DeleteUser(ctx context.Context, req *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("deleting user with id: %d", req.GetId())

	queryDeleteUser, args, err := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	_, err = s.pool.Exec(ctx, queryDeleteUser, args...)
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
