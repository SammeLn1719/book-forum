package db

import (
	"database/sql"
	"log"
	"fmt"
)

func CreateTableBooks(db *sql.DB) {
	booksTable := `CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		description TEXT,
		price DECIMAL(10,2) NOT NULL
	)`

	_, err := db.Exec(booksTable)
	if err != nil {
		log.Fatal("Failed to create books table:", err)
	}
}
func CreateTableUsers(db *sql.DB) error {
    // Создание таблицы пользователей
    usersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        password_hash VARCHAR(100) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

        // Добавьте другие таблицы по аналогии

    _, err := db.Exec(usersTable)
    if err != nil {
        return fmt.Errorf("failed to create users table: %w", err)
    }
  
    return nil
}
