package data

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(dsn string) *sql.DB {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("failed to open SQLite DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to SQLite DB: %v", err)
	}

	return db
}
