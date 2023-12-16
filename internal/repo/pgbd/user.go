package pgbd

import (
	"context"
	"errors"
	"fmt"
	"github.com/DaniiLBez/avito_internship_task/internal/entities"
	"github.com/DaniiLBez/avito_internship_task/internal/repo/repoerrs"
	"github.com/DaniiLBez/avito_internship_task/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) CreateUser(ctx context.Context, user entities.User) (int, error) {
	sql, args, _ := r.Builder.
		Insert("users").
		Columns("username", "password").
		Values(user.Username, user.Password).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrs.ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("UserRepo.CreateUser - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *UserRepo) GetUserById(ctx context.Context, id int) (entities.User, error) {
	sql, args, _ := r.Builder.
		Select("id, username, password, created_at").
		From("users").
		Where("id = ?", id).
		ToSql()

	var user entities.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.User{}, repoerrs.ErrNotFound
		}
		return entities.User{}, fmt.Errorf("UserRepo.GetUserById - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (entities.User, error) {
	sql, args, _ := r.Builder.
		Select("id, username, password, created_at").
		From("users").
		Where("username = ?", username).
		ToSql()

	var user entities.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.User{}, repoerrs.ErrNotFound
		}
		return entities.User{}, fmt.Errorf("UserRepo.GetUserByUsername - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entities.User, error) {
	sql, args, _ := r.Builder.
		Select("id, username, password, created_at").
		From("users").
		Where("username = ? AND password = ?", username, password).
		ToSql()

	var user entities.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.User{}, repoerrs.ErrNotFound
		}
		return entities.User{}, fmt.Errorf("UserRepo.GetUserByUsernameAndPassword - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}
