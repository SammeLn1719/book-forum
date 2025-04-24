package models

import "time"

type User struct {
    ID           int       `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}

type Book struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Author      string    `json:"author"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
}

type Topic struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    BookID    int       `json:"book_id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    TopicID   int       `json:"topic_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}
