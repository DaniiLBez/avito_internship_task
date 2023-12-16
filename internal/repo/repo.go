package repo

import (
	"context"
	"github.com/DaniiLBez/avito_internship_task/internal/entities"
	"github.com/DaniiLBez/avito_internship_task/internal/repo/pgbd"
	"github.com/DaniiLBez/avito_internship_task/pkg/postgres"
)

type User interface {
	CreateUser(ctx context.Context, user entities.User) (int, error)
	GetUserById(ctx context.Context, id int) (entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (entities.User, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entities.User, error)
}

type Slug interface {
	CreateSlug(ctx context.Context, slug entities.Slug) (int, error)
	DeleteSlug(ctx context.Context, slugName string) error
	GetSlugByName(ctx context.Context, slugName string) (entities.Slug, error)
}

type Membership interface {
	AddUserToSlug(ctx context.Context, slugsToAdd []int, slugsToDelete []int, userId int) error
	GetActiveSlugs(ctx context.Context, userId int) ([]entities.Slug, error)
	AddUserSlugs(ctx context.Context, slugs []int, userId int) error
	DeleteUserSlugs(ctx context.Context, slugs []int, userId int) error
}

type Repositories struct {
	User       User
	Slug       Slug
	Membership Membership
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User:       pgbd.NewUserRepo(pg),
		Slug:       pgbd.NewSlugRepo(pg),
		Membership: pgbd.NewMembershipRepo(pg),
	}
}
