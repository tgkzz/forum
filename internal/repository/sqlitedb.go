package repository

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
	if err := db.Ping(); err != nil {
		return db, err
	}

	_, err = db.Exec("PRAGMA foreign_keys=ON;")
	if err != nil {
		return db, err
	}
	if err := createTable(db, config); err != nil {
		return db, err
	}
	if err := createCategory(db); err != nil {
		return db, err
	}
	return db, err
}

func createCategory(db *sql.DB) error {
	// add another categories
	request := "SELECT COUNT(*) FROM Category WHERE name IN ('comedy', 'horror', 'drama', 'other');"

	var count int
	err := db.QueryRow(request).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	insertQuery := "INSERT INTO Category (name) VALUES ('comedy'), ('horror'), ('drama'), ('other');"
	_, err = db.Exec(insertQuery)
	return err
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
