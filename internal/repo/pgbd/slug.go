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

type SlugRepo struct {
	*postgres.Postgres
}

func NewSlugRepo(pg *postgres.Postgres) *SlugRepo {
	return &SlugRepo{pg}
}

func (r *SlugRepo) CreateSlug(ctx context.Context, slug entities.Slug) (int, error) {
	sql, args, _ := r.Builder.
		Insert("slugs").
		Columns("name").
		Values(slug.Name).
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
		return 0, fmt.Errorf("SlugRepo.CreateSlug - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *SlugRepo) DeleteSlug(ctx context.Context, slugName string) error {
	sql, args, _ := r.Builder.
		Delete("slugs").
		Where("name = ?", slugName).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return repoerrs.ErrCannotDelete
			}
		}
		return fmt.Errorf("SlugRepo.DeleteSlug - r.Pool.Exec: %v", err)
	}

	return nil
}

func (r *SlugRepo) GetSlugByName(ctx context.Context, slugName string) (entities.Slug, error) {

	sql, args, _ := r.Builder.
		Select("id, name, created_at").
		From("slugs").
		Where("name = ?", slugName).
		ToSql()

	var slug entities.Slug
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&slug.Id,
		&slug.Name,
		&slug.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Slug{}, repoerrs.ErrNotFound
		}
		return entities.Slug{}, fmt.Errorf("SlugRepo.GetSlagByName - r.Pool.QueryRow: %v", err)
	}

	return slug, nil
}
