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
	"book-forum/internal/middleware"
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
	db.CreateTableBooks(database)
	if err := db.CreateTableUsers(database); err != nil {
    log.Fatalf("Failed to create tables: %v", err) // Добавьте эту проверку
  }
	bookRepo := repository.NewBookRepository(database)
	bookHandler := handler.NewBookHandler(bookRepo)
	userRepo := repository.NewUserRepository(database)
	sessionRepo := repository.NewSessionRepository(database)
	authHandler := handler.NewAuthHandler(userRepo, sessionRepo)

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
	profileHandler := handler.ProfileHandler
	// Отдача статических файлов (картинок)
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Маршруты
	r.Get("/health", handler.HealthHandler)
	r.Route("/books", func(r chi.Router) {
		r.Get("/", bookHandler.GetAllBooks)
		r.Get("/{id}", bookHandler.GetBookByID)
		// Добавьте другие методы при необходимости
	})
	r.Post("/register", authHandler.Register)
  r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)

// Защищенные роуты
r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(sessionRepo))
		r.Get("/profile", profileHandler)
		r.Post("/logout", authHandler.Logout)
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
