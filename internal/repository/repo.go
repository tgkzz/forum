package repository

import (
	"database/sql"
	"forum/internal/repository/auth"
	"forum/internal/repository/filter"
	"forum/internal/repository/health"
	"forum/internal/repository/post"
)

type Repository struct {
	auth.Authorization
	health.Health
	post.Post
	filter.Filter
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Post:          post.NewPostRepo(db),
		Authorization: auth.NewAuthorizationRepo(db),
		Health:        health.NewHealthRepo(db),
		Filter:        filter.NewFilterRepo(db),
	}
}
