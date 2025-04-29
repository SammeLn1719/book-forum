
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

  "book-forum/internal/config"
	"book-forum/internal/db"
	"book-forum/internal/handler"
)

func main() {
	cfg := config.NewConfig()

	// Подключение к базе
	database, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	// Инициализация роутера
	r := chi.NewRouter()
	r.Get("/health", handler.HealthHandler)

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
