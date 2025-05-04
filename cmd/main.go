package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"

	"book-forum/internal/config"
	"book-forum/internal/db"
	"book-forum/internal/handler"
	"book-forum/internal/repository"
)

func main() {
	cfg := config.NewConfig()

	// Подключение к базе данных
	database, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing DB: %v", err)
		}
	}()

	// Проверка подключения к БД
	var version string
	err = database.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Failed to query DB version: %v", err)
	}
	log.Printf("Database version: %s", version)

	// Инициализация репозиториев и обработчиков
	db.CreateTable(database)
	bookRepo := repository.NewBookRepository(database)
	bookHandler := handler.NewBookHandler(bookRepo)

	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	// Инициализация роутера
	r := chi.NewRouter()
	r.Use(c.Handler)
	
	// Маршруты
	r.Get("/health", handler.HealthHandler)
	r.Route("/books", func(r chi.Router) {
	r.Get("/", bookHandler.GetAllBooks)
	r.Get("/{id}", bookHandler.GetBookByID)
		// Добавьте другие методы при необходимости
	})

	// Запуск сервера
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Server running on %s\n", addr)
	
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
