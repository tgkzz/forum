package service

import (
	"forum/internal/repo"
	authService "forum/internal/service/auth"
)

type Service struct {
	authService.AuthService
}

func NewService(repo *repo.Repo) *Service {
	return &Service{
		AuthService: *authService.NewAuthService(repo),
	}
}
