package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	DB *sql.DB
)

func InitDB() (*sql.DB, error) {
	connStr := "postgres://postgres:gulshan@localhost:5432/todoappdb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	fmt.Println("Connected to DB!")
	return db, nil
}
