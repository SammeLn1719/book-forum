package db

import (
	"database/sql"
	"fmt"
	_"github.com/lib/pq"
	"book-forum/internal/config"
)

func NewPostgres(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
