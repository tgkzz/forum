package health

import "database/sql"

type HealthRepo struct {
	DB *sql.DB
}

type Health interface {
	GetHealth() error
}

func NewHealthRepo(db *sql.DB) *HealthRepo {
	return &HealthRepo{
		DB: db,
	}
}

func (h *HealthRepo) GetHealth() error {
	return h.DB.Ping()
}
