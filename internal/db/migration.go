package db

import (
	"database/sql"
	"log"
)

func CreateTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		description TEXT,
		price DECIMAL(10,2) NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create books table:", err)
	}
}
