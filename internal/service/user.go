package service

import (
	"context"
	"github.com/DaniiLBez/avito_internship_task/internal/entities"
	"github.com/DaniiLBez/avito_internship_task/internal/repo"
)

type UserService struct {
	userRepo repo.User
}

func NewUserService(userRepo repo.User) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserById(ctx context.Context, id int) (entities.User, error) {
	return s.userRepo.GetUserById(ctx, id)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (entities.User, error) {
	return s.userRepo.GetUserByUsername(ctx, username)
}

func (s *UserService) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entities.User, error) {
	return s.userRepo.GetUserByUsernameAndPassword(ctx, username, password)
}
