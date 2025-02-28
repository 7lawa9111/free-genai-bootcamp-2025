package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return DB
} 