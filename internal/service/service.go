package service

import (
	"context"
	"github.com/DaniiLBez/avito_internship_task/internal/entities"
	"github.com/DaniiLBez/avito_internship_task/internal/repo"
	"github.com/DaniiLBez/avito_internship_task/pkg/hasher"
	"time"
)

type AuthCreateUserInput struct {
	Username string
	Password string
}

type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type Auth interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput) (int, error)
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(token string) (int, error)
}

type Slug interface {
	CreateSlug(ctx context.Context, slugName string) (int, error)
	DeleteSlug(ctx context.Context, slugName string) error
	AddUserToSlug(ctx context.Context, slugsToAdd []string, slugsToDelete []string, userId int) error
	GetActiveSlugs(ctx context.Context, userId int) ([]entities.Slug, error)
}

type User interface {
	GetUserById(ctx context.Context, id int) (entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (entities.User, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entities.User, error)
}

type Services struct {
	User User
	Slug Slug
	Auth Auth
}

type ServicesDependencies struct {
	Repos    *repo.Repositories
	Hasher   hasher.PasswordHasher
	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps *ServicesDependencies) *Services {
	return &Services{
		Auth: NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
		User: NewUserService(deps.Repos.User),
		Slug: NewSlugService(deps.Repos.User, deps.Repos.Slug, deps.Repos.Membership),
	}
}
