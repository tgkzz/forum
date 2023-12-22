package service

// worked in database, problem in post
import (
	"forum/internal/repository"
	authService "forum/internal/service/auth"
	filterService "forum/internal/service/filter"
	healthService "forum/internal/service/health"
	postService "forum/internal/service/post"
)

type Service struct {
	authService.Auth
	healthService.Healther
	postService.Poster
	filterService.Filterer
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth:     authService.NewAuthorizationService(repo.Authorization),
		Healther: healthService.NewHealthService(repo.Health),
		Poster:   postService.NewPostService(repo.Post),
		Filterer: filterService.NewFilterService(repo.Filter),
	}
}
