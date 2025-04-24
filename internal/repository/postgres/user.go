package postgres

import (
    "context"
    "database/sql"
    "fmt"
    "book-forum/internal/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (username, email, password_hash)
        VALUES ($1, $2, $3)
        RETURNING id, created_at`
    
    return r.db.QueryRowContext(ctx, query,
        user.Username,
        user.Email,
        user.PasswordHash,
    ).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    query := `SELECT id, username, email, password_hash, created_at FROM users WHERE email = $1`
    
    user := &models.User{}
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
    )
    
    if err != nil {
        return nil, fmt.Errorf("error getting user: %w", err)
    }
    
    return user, nil
}
