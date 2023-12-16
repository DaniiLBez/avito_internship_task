package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/DaniiLBez/avito_internship_task/internal/entities"
	"github.com/DaniiLBez/avito_internship_task/internal/repo"
	"github.com/DaniiLBez/avito_internship_task/internal/repo/repoerrs"
	"golang.org/x/exp/slog"
)

type SlugService struct {
	userRepo       repo.User
	slugRepo       repo.Slug
	membershipRepo repo.Membership
}

func NewSlugService(userRepo repo.User, slugRepo repo.Slug, membershipRepo repo.Membership) *SlugService {
	return &SlugService{
		userRepo:       userRepo,
		slugRepo:       slugRepo,
		membershipRepo: membershipRepo,
	}
}

func (s *SlugService) CreateSlug(ctx context.Context, slugName string) (int, error) {

	slug := entities.Slug{Name: slugName}

	slugId, err := s.slugRepo.CreateSlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return 0, ErrSlugAlreadyExists
		}
		slog.Error("SlugService.CreateSlug", err)
		return 0, ErrCannotCreateSlug
	}
	return slugId, nil
}

func (s *SlugService) DeleteSlug(ctx context.Context, slugName string) error {
	err := s.slugRepo.DeleteSlug(ctx, slugName)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return ErrSlugNotFound
		}
		slog.Error("SlugService.DeleteSlug", err)
		return ErrCannotDeleteSlug
	}
	return nil
}

func (s *SlugService) AddUserToSlug(ctx context.Context, slugsToAdd []string, slugsToDelete []string, userId int) error {
	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		switch {
		case errors.Is(err, repoerrs.ErrNotFound):
			return ErrUserNotFound
		case errors.Is(err, repoerrs.ErrCannotGet):
			return ErrCannotGetUser
		}
	}

	slugsIdToAdd := make([]int, 0)
	slugsIdToDelete := make([]int, 0)

	for _, slugName := range slugsToAdd {

		var slug entities.Slug
		slug, err = s.slugRepo.GetSlugByName(ctx, slugName)
		slog.Info(fmt.Sprintf("Founded slug: %v", slug))
		if err != nil {
			switch {
			case errors.Is(err, repoerrs.ErrNotFound):
				return ErrSlugNotFound
			case errors.Is(err, repoerrs.ErrCannotGet):
				return ErrCannotGetSlug
			}
		}
		slugsIdToAdd = append(slugsIdToAdd, slug.Id)
	}

	for _, slugName := range slugsToDelete {

		var slag entities.Slug
		slag, err = s.slugRepo.GetSlugByName(ctx, slugName)

		if err != nil {
			switch {
			case errors.Is(err, repoerrs.ErrNotFound):
				return ErrSlugNotFound
			case errors.Is(err, repoerrs.ErrCannotGet):
				return ErrCannotGetSlug
			}
		}
		slugsIdToDelete = append(slugsIdToDelete, slag.Id)
	}

	slog.Info(fmt.Sprintf("To add: %v, To delete: %v", slugsIdToAdd, slugsIdToDelete))

	return s.membershipRepo.AddUserToSlug(ctx, slugsIdToAdd, slugsIdToDelete, user.Id)
}

func (s *SlugService) GetActiveSlugs(ctx context.Context, userId int) ([]entities.Slug, error) {
	return s.membershipRepo.GetActiveSlugs(ctx, userId)
}
