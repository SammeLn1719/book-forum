package repository

import (
	"database/sql"
	"time"  // Теперь используется в методе CreateSession
	"book-forum/internal/models"
)

type SessionRepository struct {
	DB *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (r *SessionRepository) CreateSession(session *models.Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(
		query,
		session.ID,
		session.UserID,
		session.ExpiresAt.Format(time.RFC3339),  // Явное использование time
	)
	return err
}

func (r *SessionRepository) GetSession(id string) (*models.Session, error) {
	session := &models.Session{}
	var expiresAtStr string
	
	query := `SELECT id, user_id, expires_at FROM sessions WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(
		&session.ID,
		&session.UserID,
		&expiresAtStr,
	)
	
	if err != nil {
		return nil, err
	}
	
	// Парсим строку времени
	session.ExpiresAt, err = time.Parse(time.RFC3339, expiresAtStr)
	return session, err
}

func (r *SessionRepository) DeleteSession(id string) error {
	_, err := r.DB.Exec("DELETE FROM sessions WHERE id = $1", id)
	return err
}
