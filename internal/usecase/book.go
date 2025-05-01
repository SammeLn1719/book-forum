package usecase

import (
	"database/sql"
	"fmt"

	"book-forum/internal/models"
)

func InsertBook(db *sql.DB, book models.Book) (int, error) {
	query := `INSERT INTO books (title, author, description, price) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id`

	var id int
	err := db.QueryRow(query,
		book.Title,
		book.Author,
		book.Description,
		book.Price,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert book: %v", err)
	}

	return id, nil
}
