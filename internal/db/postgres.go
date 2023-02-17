package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var database *sql.DB

func ConnectDatabase(url string) *sql.DB {
	var err error
	database, err = sql.Open("pgx", url)
	if err != nil {
		log.Fatalf("Failed to open database! => %v\n", err)
	}

	if err = database.Ping(); err != nil {
		log.Fatalf("Failed to ping database! => %v\n", err)
	}

	return database
}

func CloseConnection() {
	if err := database.Close(); err != nil {
		log.Fatalf("Failed to close database! => %v\n", err)
	}
}
