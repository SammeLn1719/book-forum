package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"book-forum/internal/models"
)

var ErrBookNotFound = errors.New("book not found")

type BookRepository struct {
	DB *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (r *BookRepository) GetBookByID(id int) (*models.Book, error) {
	query := `
		SELECT id, title, author, description, price 
		FROM books 
		WHERE id = $1
	`

	var book models.Book
	err := r.DB.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Description,
		&book.Price,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return &book, nil
}

// Остальные методы остаются без изменений
func (r *BookRepository) GetAllBooks() ([]models.Book, error) {
	query := `SELECT id, title, author, description, price FROM books`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get books: %v", err)
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Description,
			&book.Price,
		); err != nil {
			return nil, fmt.Errorf("failed to scan book: %v", err)
		}
		books = append(books, book)
	}
	return books, nil
}

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
