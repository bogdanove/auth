package app

import (
	"context"
	"log"

	server "github.com/bogdanove/auth/internal/api/user"
	"github.com/bogdanove/auth/internal/client/db"
	"github.com/bogdanove/auth/internal/client/db/pg"
	"github.com/bogdanove/auth/internal/client/db/transaction"
	"github.com/bogdanove/auth/internal/closer"
	"github.com/bogdanove/auth/internal/config"
	"github.com/bogdanove/auth/internal/config/env"
	"github.com/bogdanove/auth/internal/repository"
	userRepo "github.com/bogdanove/auth/internal/repository/user"
	"github.com/bogdanove/auth/internal/service"

	userService "github.com/bogdanove/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository

	userService service.UserService

	userServer *server.Server
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PGConfig - конфигурация подключения к бд
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// GRPCConfig - конфигурация сервера GRPC
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// DBClient - клиент для базы данных
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

// TxManager - менеджер транзакций
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// UserRepository - инициализация репозитория пользователя
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewUserRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

// UserService - инициализация сервиса пользователя
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewUserService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

// UserImpl - инициализация сервера
func (s *serviceProvider) UserImpl(ctx context.Context) *server.Server {
	if s.userServer == nil {
		s.userServer = server.NewServerImplementation(s.UserService(ctx))
	}

	return s.userServer
}
