package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"book-forum/internal/config"
	"book-forum/internal/db"
	"book-forum/internal/handler"
	"book-forum/internal/models"
	"book-forum/internal/repository"
)

func main() {
	cfg := config.NewConfig()

	// Подключение к базе
	// Создаем подключение
	database, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing DB: %v", err)
		}
	}()

	// Теперь работаем напрямую с *sql.DB
	var version string
	err = database.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Failed to query DB version: %v", err)
	}
	log.Printf("Database version: %s", version)

	// Создаем таблицу
	db.CreateTable(database)

	// Добавляем книгу
	newBook := models.Book{
		Title:       "1984",
		Author:      "George Orwell",
		Description: "Dystopian novel",
		Price:       12.99,
	}

	id, err := repository.InsertBook(database, newBook)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Inserted book with ID: %d\n", id)
	bookRepo := repository.NewBookRepository(database)
	bookHandler := handler.NewBookHandler(bookRepo)
	// Инициализация роутера
	r := chi.NewRouter()
	r.Get("/health", handler.HealthHandler)
	r.Get("/books", bookHandler.GetAllBooks)
	// Старт сервера
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Server running on %s\n", addr)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
