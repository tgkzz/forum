package auth

import (
	"database/sql"
	authModel "forum/internal/model/auth"
)

type SQLAuthRepo struct {
	DB *sql.DB
}

type Authorization interface {
	InsertUser(user authModel.User) error
	// another methods will be implemented
}

func NewAuthRepo(db *sql.DB) *SQLAuthRepo {
	return &SQLAuthRepo{DB: db}
}

func (repo *SQLAuthRepo) InsertUser(user authModel.User) error {
	return nil
}
