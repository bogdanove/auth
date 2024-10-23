package user

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/bogdanove/auth/internal/client/db"
	service "github.com/bogdanove/auth/internal/model"
	"github.com/bogdanove/auth/internal/repository"
	"github.com/bogdanove/auth/internal/repository/user/model"
)

const (
	usersTableName    = "users"
	usersLogTableName = "users_log"

	idColumn             = "id"
	usersNameColumn      = "name"
	usersEmailColumn     = "email"
	usersRoleColumn      = "role"
	usersCreatedAtColumn = "created_at"
	usersUpdatedAtColumn = "updated_at"

	usersLogUserIDColumn = "user_id"
	usersLogActionColumn = "action"
)

type userRepo struct {
	db db.Client
}

// NewUserRepository - конструктор создание репозитория
func NewUserRepository(db db.Client) repository.UserRepository {
	return &userRepo{db: db}
}

// CreateUser - создание нового пользователя в системе
func (r *userRepo) CreateUser(ctx context.Context, req *service.UserInfo) (int64, error) {
	log.Printf("start creating new user with name: %s", req.Name)

	queryAddUsers, args, err := sq.Insert(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(usersNameColumn, usersEmailColumn, usersRoleColumn).
		Values(req.Name, req.Email, req.Role).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return 0, err
	}

	qc := db.Query{
		Name:     "user_repository.CreateUser",
		QueryRaw: queryAddUsers,
	}

	var userID int64
	err = r.db.DB().QueryRowContext(ctx, qc, args...).Scan(&userID)
	if err != nil {
		log.Printf("failed to insert user: %v", err)
		return 0, err
	}

	return userID, nil
}

// GetUser - получение информации о пользователе по его идентификатору
func (r *userRepo) GetUser(ctx context.Context, req int64) (*model.User, error) {
	log.Printf("receiving user with id: %d", req)

	query, args, err := sq.Select(idColumn, usersNameColumn, usersEmailColumn,
		usersRoleColumn, usersCreatedAtColumn, usersUpdatedAtColumn).
		From(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: req}).
		Limit(1).ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetUser",
		QueryRaw: query,
	}

	var user model.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		log.Printf("failed to select user: %v", err)
		return nil, err
	}

	return &user, nil
}

// UpdateUser - обновление информации о пользователе по его идентификатору
func (r *userRepo) UpdateUser(ctx context.Context, req *service.UpdateUserInfo) error {
	log.Printf("updating user with id: %d", req.ID)

	builder := sq.Update(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Set(usersUpdatedAtColumn, time.Now())
	if req.Name != nil {
		builder = builder.Set(usersNameColumn, req.Name)
	}
	if req.Role != nil {
		builder = builder.Set(usersRoleColumn, req.Role)
	}

	query, args, err := builder.Where(sq.Eq{idColumn: req.ID}).ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	q := db.Query{
		Name:     "user_repository.UpdateUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return err
	}

	return nil
}

// DeleteUser - удаление пользователя из системы по его идентификатору
func (r *userRepo) DeleteUser(ctx context.Context, req int64) error {
	log.Printf("deleting user with id: %d", req)

	queryDeleteUser, args, err := sq.Delete(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: req}).ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	q := db.Query{
		Name:     "user_repository.DeleteUser",
		QueryRaw: queryDeleteUser,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return err
	}

	return nil
}

// SaveLog - сохранение записи о действиях пользователя
func (r *userRepo) SaveLog(ctx context.Context, req *model.UserLog) error {
	log.Printf("create new users_log with user_id: %d", req.UserID)

	queryUserLog, args, err := sq.Insert(usersLogTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(usersLogUserIDColumn, usersLogActionColumn).
		Values(req.UserID, req.Action).ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	ql := db.Query{
		Name:     "user_repository.SaveLog",
		QueryRaw: queryUserLog,
	}

	_, err = r.db.DB().ExecContext(ctx, ql, args...)
	if err != nil {
		log.Printf("failed to insert users_log: %v", err)
		return err
	}
	return nil
}
