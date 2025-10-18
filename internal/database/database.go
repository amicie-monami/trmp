package database

import (
	"database/sql"

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

	return db, nil
}
