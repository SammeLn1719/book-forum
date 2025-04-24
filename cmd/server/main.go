package main

import (
    "context"
    "database/sql"
    "log"
    "net/http"
    
    "book-forum/internal/config"
    "book-forum/internal/handler"
    "book-forum/internal/repository/postgres"
    "book-forum/internal/service"
    
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    _ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Error loading config:", err)
    }

    db, err := sql.Open("pgx", 
        "postgres://"+cfg.DBUser+":"+cfg.DBPassword+"@"+cfg.DBHost+":"+cfg.DBPort+"/"+cfg.DBName)
    if err != nil {
        log.Fatal("Error connecting to database:", err)
    }
    defer db.Close()

    // Инициализация репозиториев
    userRepo := postgres.NewUserRepository(db)
    
    // Инициализация сервисов
    authService := service.NewAuthService(userRepo, cfg.JWTSecret)
    
    // Инициализация обработчиков
    authHandler := handler.NewAuthHandler(authService)
    
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    
    r.Post("/api/register", authHandler.Register)
    r.Post("/api/login", authHandler.Login)
    
    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
