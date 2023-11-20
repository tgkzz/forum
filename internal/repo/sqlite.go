package repo

import (
	"database/sql"
	"forum/config"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqlite(config config.Config) (*sql.DB, error) {
	db, err := sql.Open(config.DB.DriverName, config.DB.DataSourceName)
	if err != nil {
		return db, err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return db, err
	}

	if err := createTable(db, config); err != nil {
		return db, err
	}

	return db, err
}

func createTable(db *sql.DB, config config.Config) error {
	fileSql, err := os.ReadFile(config.DB.Sql)
	if err != nil {
		return err
	}

	requests := strings.Split(string(fileSql), ";")
	for _, request := range requests {
		_, err = db.Exec(request)
		if err != nil {
			return err
		}
	}

	return err
}
