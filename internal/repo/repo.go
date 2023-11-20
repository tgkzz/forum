package repo

import (
	"database/sql"
	authRepo "forum/internal/repo/auth"
)

type Repo struct {
	authRepo.Authorization
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		Authorization: authRepo.NewAuthRepo(db),
	}
}
