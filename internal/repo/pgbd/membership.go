package pgbd

import (
	"context"
	"errors"
	"fmt"
	"github.com/DaniiLBez/avito_internship_task/internal/entities"
	"github.com/DaniiLBez/avito_internship_task/internal/repo/repoerrs"
	"github.com/DaniiLBez/avito_internship_task/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

type MembershipRepo struct {
	*postgres.Postgres
}

func NewMembershipRepo(pb *postgres.Postgres) *MembershipRepo {
	return &MembershipRepo{pb}
}

func (r *MembershipRepo) AddUserToSlug(ctx context.Context, slugsToAdd []int, slugsToDelete []int, userId int) error {
	err := r.AddUserSlugs(ctx, slugsToAdd, userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return repoerrs.ErrAlreadyExists
			}
		}
		return fmt.Errorf("MembershipRepo.AddUserToSlug - r.Pool.Exec: %v", err)
	}

	err = r.DeleteUserSlugs(ctx, slugsToDelete, userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return repoerrs.ErrCannotDelete
			}
		}
		return fmt.Errorf("MembershipRepo.AddUserToSlug - r.Pool.Exec: %v", err)
	}

	return nil
}

func (r *MembershipRepo) GetActiveSlugs(ctx context.Context, userId int) ([]entities.Slug, error) {
	sql, args, _ := r.Builder.
		Select("slugs.id", "slugs.name", "slugs.created_at").
		From("membership").
		Join("slugs ON membership.slug_id = slugs.id").
		Where("membership.user_id = ?", userId).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("MembershipRepo.GetActiveSlugs - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	var activeSlugs []entities.Slug
	for rows.Next() {
		var slug entities.Slug
		if err = rows.Scan(&slug.Id, &slug.Name, &slug.CreatedAt); err != nil {
			return nil, fmt.Errorf("MembershipRepo.GetActiveSlugs - rows.Scan: %v", err)
		}
		activeSlugs = append(activeSlugs, slug)
	}

	return activeSlugs, nil
}

func (r *MembershipRepo) AddUserSlugs(ctx context.Context, slugs []int, userId int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("MembershipRepo.AddUserSlugs - r.Pool.Begin: %v\n", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	builder := r.Builder.Insert("membership").Columns("user_id", "slug_id")

	for _, id := range slugs {
		builder = builder.Values(userId, id)
	}

	sql, args, _ := builder.ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("MembershipRepo.AddUserSlugs - tx.Exec: %v\n", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("MembershipRepo.AddUserSlugs - tx.Commit: %v\n", err)
	}

	return nil
}

func (r *MembershipRepo) DeleteUserSlugs(ctx context.Context, slugs []int, userId int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("MembershipRepo.DeleteUserSlugs - r.Pool.Begin: %v\n", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := r.Builder.Delete("membership").
		Where(
			squirrel.And{
				squirrel.Eq{"slug_id": slugs}, squirrel.Eq{"user_id": userId},
			}).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("MembershipRepo.DeleteUserSlugs - tx.Exec: %v\n", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("MembershipRepo.DeleteUserSlugs - tx.Commit: %v\n", err)
	}

	return nil
}
