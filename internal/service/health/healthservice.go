package health

import "forum/internal/repository/health"

type HealtServive struct {
	repo health.Health
}

type Healther interface {
	GetHealth() error
}

func NewHealthService(repository health.Health) *HealtServive {
	return &HealtServive{
		repo: repository,
	}
}

func (h *HealtServive) GetHealth() error {
	return h.repo.GetHealth()
}
