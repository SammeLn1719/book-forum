package repository

import (
    "database/sql"
    "book-forum/internal/models"
    "book-forum/internal/utils"
)

type UserRepository struct {
    DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
    hashedPassword, err := utils.HashPassword(user.PasswordHash)
    if err != nil {
        return err
    }

    query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
    return r.DB.QueryRow(query, user.Username, user.Email, hashedPassword).Scan(&user.ID)
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, username, email, password_hash, created_at FROM users WHERE email = $1`
    err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
    return user, err
}
func (r *UserRepository) IsEmailExists(email string) (bool, error) {
    var exists bool
    err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
    return exists, err
}
