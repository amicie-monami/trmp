package database

import (
	"database/sql"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func InitSQLiteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./writers.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	sqlBytes, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	queries := strings.Split(string(sqlBytes), ";")
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	sqlBytes, err = os.ReadFile("testdata.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(sqlBytes)); err != nil {
		return err
	}

	return nil
}
