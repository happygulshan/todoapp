package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() *sql.DB {
	var err error
	connStr := "postgres://postgres:gulshan@localhost:5432/todoappdb?sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("error :)")
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		fmt.Println("error :) 2")
		panic(err)
	}

	fmt.Println("Connected to DB!")

	return DB

}
