package models

import "time"

type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username" validate:"required,min=3,max=50"`
    Email     string    `json:"email" validate:"required,email"`
    PasswordHash   string    `json:"password_hash" validate:"required,min=8"`
    CreatedAt time.Time `json:"created_at"`
}

