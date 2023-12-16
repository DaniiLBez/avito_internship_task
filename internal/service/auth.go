package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/DaniiLBez/avito_internship_task/internal/entities"
	"github.com/DaniiLBez/avito_internship_task/internal/repo"
	"github.com/DaniiLBez/avito_internship_task/internal/repo/repoerrs"
	"github.com/DaniiLBez/avito_internship_task/pkg/hasher"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/exp/slog"
	"time"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int
}

type AuthService struct {
	userRepo       repo.User
	passwordHasher hasher.PasswordHasher
	signKey        string
	tokenTTL       time.Duration
}

func NewAuthService(
	userRepo repo.User,
	passwordHasher hasher.PasswordHasher,
	signKey string,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{userRepo: userRepo, passwordHasher: passwordHasher, signKey: signKey, tokenTTL: tokenTTL}
}

func (s *AuthService) CreateUser(ctx context.Context, input AuthCreateUserInput) (int, error) {

	user := entities.User{
		Username: input.Username,
		Password: s.passwordHasher.Hash(input.Password),
	}

	userId, err := s.userRepo.CreateUser(ctx, user)

	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return 0, ErrUserAlreadyExists
		}
		slog.Error("AuthService.CreateUser", err)
		return 0, ErrCannotCreateUser
	}
	return userId, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error) {
	// get user from DB
	user, err := s.userRepo.GetUserByUsernameAndPassword(ctx, input.Username, s.passwordHasher.Hash(input.Password))
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return "", ErrUserNotFound
		}
		slog.Error("AuthService.GenerateToken: cannot get user", err)
		return "", ErrCannotGetUser
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		UserId: user.Id,
	})

	// sign token
	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		slog.Error("AuthService.GenerateToken: cannot sign token", err)
		return "", ErrCannotSignToken
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return 0, ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, ErrCannotParseToken
	}

	return claims.UserId, nil
}
